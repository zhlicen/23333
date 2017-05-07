package beeaccount

import (
	"errors"

	"github.com/astaxie/beego"
	context "github.com/astaxie/beego/context"
	"github.com/astaxie/beego/session"
)

type accountInvoker struct {
	model        AccountModel
	loginUserId  *UserId
	invokeUserId *UserId
	context      *context.Context
	mgr          *AccountMgr
}

type AccountAction int

const (
	Account_LogoutSession AccountAction = iota
	Account_ChangePwd
	Account_GetAccountBaseInfo
	Account_UpdateAccountBaseInfo
	Account_GetProfiles
	Account_UpdateProfiles
	Account_GetOthers
	Account_UpdateOthers
)

func (a *accountInvoker) checkPermission(action AccountAction) error {
	if a.mgr.pc != nil {
		return a.mgr.pc.Check(a.loginUserId, a.invokeUserId, action)
	}
	return nil
}

func (a *accountInvoker) LogoutSession(sid string) error {
	perErr := a.checkPermission(Account_LogoutSession)
	if perErr != nil {
		return perErr
	}
	curUser := a.invokeUserId
	if curUser == nil {
		return errors.New("Illgal invoke")
	}

	curSs := a.context.Input.CruSession
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

	if a.loginUserId != nil && a.loginUserId.Uid == curUser.Uid {
		if curUser.Uid != ssUser.(UserId).Uid {
			return errors.New("login user not match")
		}
	}

	ss.Flush()
	ss.SessionRelease(a.context.ResponseWriter)

	return nil
}

func (a *accountInvoker) ChangePwd(oldPwd *LoginPwd, newPwd *LoginPwd) error {
	perErr := a.checkPermission(Account_ChangePwd)
	if perErr != nil {
		return perErr
	}
	uid := a.invokeUserId.Uid
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

func (a *accountInvoker) GetAccountBaseInfo() (*AccountBaseInfo, error) {
	perErr := a.checkPermission(Account_GetAccountBaseInfo)
	if perErr != nil {
		return nil, perErr
	}
	uid := a.invokeUserId.Uid
	return a.model.GetAccountBaseInfo(uid)
}

func (a *accountInvoker) UpdateAccountBaseInfo(baseInfo *AccountBaseInfo) error {
	perErr := a.checkPermission(Account_UpdateAccountBaseInfo)
	if perErr != nil {
		return perErr
	}
	uid := a.invokeUserId.Uid
	return a.model.UpdateAccountBaseInfo(uid, baseInfo)
}

func (a *accountInvoker) GetProfiles() (map[KeyName]string, error) {
	perErr := a.checkPermission(Account_GetProfiles)
	if perErr != nil {
		return nil, perErr
	}
	uid := a.invokeUserId.Uid
	return a.model.GetProfiles(uid)
}

func (a *accountInvoker) UpdateProfiles(profiles map[KeyName]string) error {
	perErr := a.checkPermission(Account_UpdateProfiles)
	if perErr != nil {
		return perErr
	}
	uid := a.invokeUserId.Uid
	return a.model.UpdateProfiles(uid, profiles)
}

func (a *accountInvoker) GetOthers() (map[KeyName]string, error) {
	perErr := a.checkPermission(Account_GetOthers)
	if perErr != nil {
		return nil, perErr
	}
	uid := a.invokeUserId.Uid
	return a.model.GetOthers(uid)
}

func (a *accountInvoker) UpdateOthers(others map[KeyName]string) error {
	perErr := a.checkPermission(Account_UpdateOthers)
	if perErr != nil {
		return perErr
	}
	uid := a.invokeUserId.Uid
	return a.model.UpdateOthers(uid, others)

}
