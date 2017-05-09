package beeaccount

import (
	"23333/utils/web/beenh/beepermission"
	"23333/utils/web/verify"
	"errors"

	context "github.com/astaxie/beego/context"
)

// mgrsStrore managers storage
var mgrsStore map[string]*AccountMgr

// AccountMgr account manager
// member:domain domain of this manager
// member:pc permission checker of account actions
type AccountMgr struct {
	domain string
	model  AccountModel
	pc     beepermission.PermissionChecker
}

// GetAccountMgr get account manager with domain name
func GetAccountMgr(domain string) (*AccountMgr, error) {
	if mgrsStore != nil {
		if mgr, ok := mgrsStore[domain]; ok {
			return mgr, nil
		}
	}
	return nil, errors.New("not exist")
}

// NewAccountMgr constructor of AccountMgr
func NewAccountMgr(domain string, model AccountModel, pc beepermission.PermissionChecker) *AccountMgr {
	if mgrsStore == nil {
		mgrsStore = make(map[string]*AccountMgr)
	} else {
		if mgr, ok := mgrsStore[domain]; ok {
			mgr.model = model
			return mgr
		}
	}
	if model == nil {
		return nil
	}
	mgr := &AccountMgr{domain, model, pc}
	mgrsStore[domain] = mgr
	return mgr
}

// GetLoginUserId get login user id
func (a *AccountMgr) GetLoginUserId(c *context.Context) *UserId {
	ss := c.Input.CruSession
	if userId, ok := ss.Get("LoginUser").(UserId); ok {
		return &userId
	}
	return nil
}

// CurrentAccount get current(login) account invoker
func (a *AccountMgr) CurrentAccount(c *context.Context) *accountInvoker {
	userId := a.GetLoginUserId(c)
	return &accountInvoker{a.model, userId, userId, c, a}
}

// OtherAccount get other account invoker
func (a *AccountMgr) OtherAccount(c *context.Context, userId *UserId) *accountInvoker {
	return &accountInvoker{a.model, a.GetLoginUserId(c), userId, c, a}
}

// Register register an account with account info
func (a *AccountMgr) Register(c *context.Context, info *AccountInfo) error {
	ss := c.Input.CruSession
	ssUser := ss.Get("LoginUser")
	if ssUser != nil {
		return errors.New("account is logged in")
	}
	valErr := info.Validate()
	if valErr != nil {
		return valErr
	}
	return a.model.Add(info)
}

// VerifyId verify account login id
func (a *AccountMgr) VerifyId(c *context.Context, v *verify.Verifier, loginId string) error {
	if v == nil {
		return errors.New("invalid verifier")
	}
	vErr := v.Verify()
	if vErr == nil {
		userId, uidErr := a.GetUserId(c, loginId)
		if uidErr != nil {
			return uidErr
		}
		accountSchema, schemaErr := GetAccountSchema(a.domain)
		if schemaErr != nil {
			return schemaErr
		}
		name, matchErr := accountSchema.MatchLoginId(loginId)
		if matchErr != nil {
			return matchErr
		}
		basicInfo, getAccountErr := a.model.GetAccountBasicInfo(userId.Uid)
		if getAccountErr != nil {
			return getAccountErr
		}
		basicInfo.LoginIds[name] = NewLoginId(basicInfo.LoginIds[name].Id, true)
		return a.model.UpdateAccountBasicInfo(userId.Uid, basicInfo)
	}
	return errors.New("invalid token")
}

// ResetPwd reset password
func (a *AccountMgr) ResetPwd(c *context.Context, v *verify.Verifier, loginId string, newPwd *LoginPwd) error {
	if v == nil {
		return errors.New("invalid verifier")
	}
	vErr := v.Verify()
	if vErr == nil {
		userId, uidErr := a.GetUserId(c, loginId)
		if uidErr != nil {
			return errors.New("unknown account")
		}

		AccountBasicInfo, accountErr := a.model.GetAccountBasicInfo(userId.Uid)
		if accountErr != nil {
			return accountErr
		}

		newPwdRaw, newPwdErr := newPwd.GetPwd()
		if newPwdErr == nil {
			AccountBasicInfo.Password.SetEncryptedPwd(newPwdRaw)
			return a.model.UpdateAccountBasicInfo(userId.Uid, AccountBasicInfo)
		}
	}
	return errors.New("invalid token")
}

// GetUserId get user id by login string
// Universal Interface, can be called without login
func (a *AccountMgr) GetUserId(c *context.Context, loginId string) (*UserId, error) {
	accountSchema, schemaErr := GetAccountSchema(a.domain)
	if schemaErr != nil {
		return nil, schemaErr
	}
	name, matchErr := accountSchema.MatchLoginId(loginId)
	if matchErr != nil {
		return nil, matchErr
	}
	return a.model.GetUserId(name, loginId)
}

// Login login an account with password
func (a *AccountMgr) Login(c *context.Context, loginId string, pwd *LoginPwd) error {
	userId, findErr := a.GetUserId(c, loginId)
	if findErr != nil {
		return errors.New("invalid user id")
	}
	basicInfo, accountErr := a.model.GetAccountBasicInfo(userId.Uid)
	if accountErr != nil {
		return accountErr
	}
	accountSchema, schemaErr := GetAccountSchema(a.domain)
	if schemaErr != nil {
		return schemaErr
	}
	name, matchErr := accountSchema.MatchLoginId(loginId)
	if matchErr != nil {
		return matchErr
	}
	if !basicInfo.LoginIds[name].Verified {
		return errors.New("account id not verified")
	}
	LoginPwd, pwdErr := basicInfo.Password.GetPwd()
	if pwdErr == nil {
		userPwd, _ := pwd.GetPwd()
		if userPwd != LoginPwd {
			return errors.New("invalid password")
		}
	} else {
		return pwdErr
	}
	ss := c.Input.CruSession
	ss.Set("LoginUser", basicInfo.UserId)
	return nil
}
