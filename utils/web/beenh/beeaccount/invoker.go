package beeaccount

import (
	"errors"

	"github.com/astaxie/beego"
	context "github.com/astaxie/beego/context"
	"github.com/astaxie/beego/session"
)

// accountInvoker invoker for account functions
// member:model model implementation
// member:loginUserId current login user id
// member:invokeUserId invoke user id
// member:context beego context
// member:mgr account manager
type accountInvoker struct {
	model        AccountModel
	loginUserId  *UserId
	invokeUserId *UserId
	context      *context.Context
	mgr          *AccountMgr
}

// AccountAction account action
type AccountAction int

// Account Actions definition
const (
	Account_LogoutSession AccountAction = iota
	Account_ChangePwd
	Account_GetAccountBasicInfo
	Account_UpdateAccountBasicInfo
	Account_GetProfiles
	Account_UpdateProfiles
	Account_GetOthers
	Account_UpdateOthers
)

// checkPermission check permission of this action
func (a *accountInvoker) checkPermission(action AccountAction) error {
	if a.mgr.pc != nil {
		return a.mgr.pc.Check(a.loginUserId, a.invokeUserId, action)
	}
	return nil
}

// LogoutSession logout session with session id
// sid is the id of session to be logged out
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

// ChangePwd change password
// oldPwd is the old password
// newPwd is the new password
// returns error
func (a *accountInvoker) ChangePwd(oldPwd *LoginPwd, newPwd *LoginPwd) error {
	perErr := a.checkPermission(Account_ChangePwd)
	if perErr != nil {
		return perErr
	}
	uid := a.invokeUserId.Uid
	AccountBasicInfo, accountErr := a.model.GetAccountBasicInfo(uid)
	if accountErr != nil {
		return accountErr
	}
	oldPwdSaved, oldPwdSavedErr := AccountBasicInfo.Password.GetPwd()
	oldPwdUser, oldPwdUserErr := oldPwd.GetPwd()

	if oldPwdSavedErr == nil && oldPwdUserErr == nil && oldPwdSaved == oldPwdUser {
		newPwdRaw, newPwdErr := newPwd.GetPwd()
		if newPwdErr == nil {
			AccountBasicInfo.Password.SetEncryptedPwd(newPwdRaw)
			return a.model.UpdateAccountBasicInfo(uid, AccountBasicInfo)
		}
	}
	return nil
}

// GetAccountBasicInfo get account basic info
func (a *accountInvoker) GetAccountBasicInfo() (*AccountBasicInfo, error) {
	perErr := a.checkPermission(Account_GetAccountBasicInfo)
	if perErr != nil {
		return nil, perErr
	}
	uid := a.invokeUserId.Uid
	return a.model.GetAccountBasicInfo(uid)
}

// UpdateAccountBasicInfo update account basic info
func (a *accountInvoker) UpdateAccountBasicInfo(basicInfo *AccountBasicInfo) error {
	perErr := a.checkPermission(Account_UpdateAccountBasicInfo)
	if perErr != nil {
		return perErr
	}
	uid := a.invokeUserId.Uid
	return a.model.UpdateAccountBasicInfo(uid, basicInfo)
}

// GetProfiles get profiles
func (a *accountInvoker) GetProfiles() (map[string]string, error) {
	perErr := a.checkPermission(Account_GetProfiles)
	if perErr != nil {
		return nil, perErr
	}
	uid := a.invokeUserId.Uid
	return a.model.GetProfiles(uid)
}

// UpdateProfiles update profiles
func (a *accountInvoker) UpdateProfiles(profiles map[string]string) error {
	perErr := a.checkPermission(Account_UpdateProfiles)
	if perErr != nil {
		return perErr
	}
	uid := a.invokeUserId.Uid
	return a.model.UpdateProfiles(uid, profiles)
}

// GetOthers get others
func (a *accountInvoker) GetOthers() (map[string]string, error) {
	perErr := a.checkPermission(Account_GetOthers)
	if perErr != nil {
		return nil, perErr
	}
	uid := a.invokeUserId.Uid
	return a.model.GetOthers(uid)
}

// UpdateOthers update others
func (a *accountInvoker) UpdateOthers(others map[string]string) error {
	perErr := a.checkPermission(Account_UpdateOthers)
	if perErr != nil {
		return perErr
	}
	uid := a.invokeUserId.Uid
	return a.model.UpdateOthers(uid, others)

}
