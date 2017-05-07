package beeaccount

import (
	"23333/utils/web/beenh/beepermission"
	"23333/utils/web/verify"
	"errors"

	context "github.com/astaxie/beego/context"
)

var mgrsStore map[string]*AccountMgr

type sessionKeys struct {
	loginUser string
}

type AccountMgr struct {
	domain string
	model  AccountModel
	pc     beepermission.PermissionChecker
}

func GetAccountMgr(domain string) (*AccountMgr, error) {
	if mgrsStore != nil {
		if mgr, ok := mgrsStore[domain]; ok {
			return mgr, nil
		}
	}
	return nil, errors.New("not exist")
}

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

func (a *AccountMgr) GetLoginUserId(c *context.Context) *UserId {
	ss := c.Input.CruSession
	if userId, ok := ss.Get("LoginUser").(UserId); ok {
		return &userId
	}
	return nil
}

func (a *AccountMgr) CurrentAccount(c *context.Context) *accountInvoker {
	userId := a.GetLoginUserId(c)
	return &accountInvoker{a.model, userId, userId, c, a}
}

func (a *AccountMgr) OtherAccount(c *context.Context, userId *UserId) *accountInvoker {
	return &accountInvoker{a.model, a.GetLoginUserId(c), userId, c, a}
}

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
		desc, _ := MatchIdDescriptor(loginId)
		baseInfo, getAccountErr := a.model.GetAccountBaseInfo(userId.Uid)
		if getAccountErr != nil {
			return getAccountErr
		}
		baseInfo.LoginIds[desc.Name] = NewLoginId(baseInfo.LoginIds[desc.Name].Id, true)
		return a.model.UpdateAccountBaseInfo(userId.Uid, baseInfo)
	}
	return errors.New("invalid token")
}

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

		accountBaseInfo, accountErr := a.model.GetAccountBaseInfo(userId.Uid)
		if accountErr != nil {
			return accountErr
		}

		newPwdRaw, newPwdErr := newPwd.GetPwd()
		if newPwdErr == nil {
			accountBaseInfo.Password.SetEncryptedPwd(newPwdRaw)
			return a.model.UpdateAccountBaseInfo(userId.Uid, accountBaseInfo)
		}
	}
	return errors.New("invalid token")
}

// Universal Interface, can be called without login
func (a *AccountMgr) GetUserId(c *context.Context, loginId string) (*UserId, error) {
	desc, matchErr := MatchIdDescriptor(loginId)
	if matchErr != nil {
		return nil, matchErr
	}
	return a.model.GetUserId(desc.Name, loginId)
}

func (a *AccountMgr) Login(c *context.Context, loginId string, pwd *LoginPwd) error {
	userId, findErr := a.GetUserId(c, loginId)
	if findErr != nil {
		return errors.New("invalid user id")
	}
	baseInfo, accountErr := a.model.GetAccountBaseInfo(userId.Uid)
	if accountErr != nil {
		return accountErr
	}
	desc, _ := MatchIdDescriptor(loginId)
	if !baseInfo.LoginIds[desc.Name].Verified {
		return errors.New("account id not verified")
	}
	LoginPwd, pwdErr := baseInfo.Password.GetPwd()
	if pwdErr == nil {
		userPwd, _ := pwd.GetPwd()
		if userPwd != LoginPwd {
			return errors.New("invalid password")
		}
	} else {
		return pwdErr
	}
	ss := c.Input.CruSession
	ss.Set("LoginUser", baseInfo.UserId)
	return nil
}
