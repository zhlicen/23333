package beaccount

import (
	"23333/utils/web/verify"
	"errors"

	"fmt"

	"github.com/astaxie/beego"
	context "github.com/astaxie/beego/context"
	"github.com/astaxie/beego/session"
)

var mgrsStore map[string]*AccountMgr

type sessionKeys struct {
	loginUser string
}

type accountInvoker struct {
	model        AccountModel
	loginUserId  *UserId
	invokeUserId *UserId
	context      *context.Context
}

type AccountMgr struct {
	domain string
	model  AccountModel
}

func GetAccountMgr(domain string) (*AccountMgr, error) {
	if mgrsStore != nil {
		if mgr, ok := mgrsStore[domain]; ok {
			return mgr, nil
		}
	}
	return nil, errors.New("not exist")
}

func NewAccountMgr(domain string, model AccountModel) *AccountMgr {
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
	mgr := &AccountMgr{domain, model}
	mgrsStore[domain] = mgr
	return mgr
}

func (a *AccountMgr) getOperateUid(c *context.Context) (string, error) {
	ss := c.Input.CruSession
	ssUser := ss.Get("LoginUser")
	otherUser := ss.Get("OtherUser")

	if ssUser != nil {
		return ssUser.(string), nil
	}
	if otherUser != nil {
		return otherUser.(string), nil
	}
	return "", errors.New("account is logged in")
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
	return &accountInvoker{a.model, userId, userId, c}
}

func (a *AccountMgr) IOtherAccount(c *context.Context, userId *UserId) *accountInvoker {
	return &accountInvoker{a.model, a.GetLoginUserId(c), userId, c}
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

// Require user verify
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

func (a *AccountMgr) LogoutSession(c *context.Context, sid string) error {
	fmt.Println("loging out " + sid)
	curSs := c.Input.CruSession
	var ss session.Store
	if sid == curSs.SessionID() {
		ss = curSs
	} else {
		var ssErr error
		ss, ssErr = beego.GlobalSessions.GetSessionStore(sid)
		if ssErr != nil {
			return ssErr
		}
	}

	ssUser := ss.Get("LoginUser")
	if ssUser == nil {
		return errors.New("no login user with this session")
	}
	curUser, _ := a.getOperateUid(c)
	fmt.Println(curUser + " : " + ssUser.(string))
	if curUser != ssUser.(string) {
		return errors.New("login user not match")
	}
	fmt.Println("Logged out")
	// ****
	ss.Flush()
	ss.SessionRelease(c.ResponseWriter)
	fmt.Println(ss.Get("LoginUser"))

	return nil
}

func (a *AccountMgr) ChangePwd(c *context.Context, oldPwd *LoginPwd, newPwd *LoginPwd) error {
	uid, err := a.getOperateUid(c)
	if err != nil {
		return err
	}
	accountBaseInfo, accountErr := a.model.GetAccountBaseInfo(uid)
	if accountErr != nil {
		return accountErr
	}
	oldPwdSaved, oldPwdSavedErr := accountBaseInfo.Password.GetPwd()
	oldPwdUser, oldPwdUserErr := oldPwd.GetPwd()

	if oldPwdSavedErr == nil && oldPwdUserErr == nil && oldPwdSaved == oldPwdUser {
		newPwdRaw, newPwdErr := newPwd.GetPwd()
		if newPwdErr == nil {
			accountBaseInfo.Password.SetEncryptedPwd(newPwdRaw)
			return a.model.UpdateAccountBaseInfo(uid, accountBaseInfo)
		}
	}
	return nil
}

func (a *AccountMgr) GetAccountBaseInfo(c *context.Context) (*AccountBaseInfo, error) {
	uid, err := a.getOperateUid(c)
	if err != nil {
		return nil, err
	}
	return a.model.GetAccountBaseInfo(uid)
}

func (a *AccountMgr) UpdateAccountBaseInfo(c *context.Context, baseInfo *AccountBaseInfo) error {
	uid, err := a.getOperateUid(c)
	if err != nil {
		return nil
	}
	return a.model.UpdateAccountBaseInfo(uid, baseInfo)
}

func (a *AccountMgr) GetProfiles(c *context.Context) (map[KeyName]string, error) {
	uid, err := a.getOperateUid(c)
	if err != nil {
		return nil, err
	}
	return a.model.GetProfiles(uid)
}

func (a *AccountMgr) UpdateProfiles(c *context.Context, profiles map[KeyName]string) error {
	uid, err := a.getOperateUid(c)
	if err != nil {
		return err
	}
	return a.model.UpdateProfiles(uid, profiles)
}

func (a *AccountMgr) GetOthers(c *context.Context) (map[KeyName]string, error) {
	uid, err := a.getOperateUid(c)
	if err != nil {
		return nil, err
	}
	return a.model.GetOthers(uid)
}

func (a *AccountMgr) UpdateOthers(c *context.Context, others map[KeyName]string) error {
	uid, err := a.getOperateUid(c)
	if err != nil {
		return err
	}
	return a.model.UpdateOthers(uid, others)

}

// Operation on other account, needs permission check
func (a *AccountMgr) OtherAccount(c *context.Context, uid string) *AccountMgr {
	ss := c.Input.CruSession
	ss.Set("OtherUser", uid)
	return a
}
