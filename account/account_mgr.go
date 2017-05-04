package account

import (
	"23333/utilities/verify"
	"errors"

	"fmt"

	"github.com/astaxie/beego"
	context "github.com/astaxie/beego/context"
	"github.com/astaxie/beego/session"
)

type AccountAction int

const (
	Register AccountAction = iota
	Login
	LogoutSession
)

type AccountActionFilter interface {
	onActionStart(action AccountAction, c *context.Context, params ...interface{}) bool
	onActionResult(action AccountAction, c *context.Context, params ...interface{})
}

type sessionKeys struct {
	loginUser string
}

type AccountMgr struct {
	domain string
	model  AccountModel
}

func NewAccountService(domain string, model AccountModel) *AccountMgr {
	if model == nil {
		return nil
	}
	return &AccountMgr{domain, model}
}

func (s *AccountMgr) getOperateUid(c *context.Context) (AccountUid, error) {
	ss := c.Input.CruSession
	ssUser := ss.Get("LoginUser")
	otherUser := ss.Get("OtherUser")

	if ssUser != nil {
		return AccountUid(ssUser.(string)), nil
	}
	if otherUser != nil {
		return AccountUid(otherUser.(string)), nil
	}
	return "", errors.New("account is logged in")
}

// Can not be called when logged in
func (s *AccountMgr) Register(c *context.Context, info *AccountInfo) error {
	ss := c.Input.CruSession
	ssUser := ss.Get("LoginUser")
	if ssUser != nil {
		return errors.New("account is logged in")
	}
	valErr := info.Validate()
	if valErr != nil {
		return valErr
	}
	return s.model.Add(info)
}

func (s *AccountMgr) VerifyId(c *context.Context, v *verify.Verifier, id string) error {
	if v == nil {
		return errors.New("invalid verifier")
	}
	vErr := v.Verify()
	if vErr == nil {
		uid, uidErr := s.GetUidById(c, id)
		if uidErr != nil {
			return uidErr
		}
		desc, _ := MatchIdDescriptor(id)
		baseInfo, getAccountErr := s.model.GetAccountBaseInfo(uid)
		if getAccountErr != nil {
			return getAccountErr
		}
		baseInfo.Ids[desc.Name] = NewAccountId(baseInfo.Ids[desc.Name].Id, true)
		return s.model.UpdateAccountBaseInfo(uid, baseInfo)
	}
	return errors.New("invalid token")
}

// Require user verify
func (s *AccountMgr) ResetPwd(c *context.Context, v *verify.Verifier, id string, newPwd *AccountPwd) error {
	if v == nil {
		return errors.New("invalid verifier")
	}
	vErr := v.Verify()
	if vErr == nil {
		uid, uidErr := s.GetUidById(c, id)
		if uidErr != nil {
			return errors.New("unknown account")
		}

		accountBaseInfo, accountErr := s.model.GetAccountBaseInfo(uid)
		if accountErr != nil {
			return accountErr
		}

		newPwdRaw, newPwdErr := newPwd.GetPwd()
		if newPwdErr == nil {
			accountBaseInfo.Password.SetEncryptedPwd(newPwdRaw)
			return s.model.UpdateAccountBaseInfo(uid, accountBaseInfo)
		}
	}
	return errors.New("invalid token")
}

// Universal Interface, can be called without login
func (s *AccountMgr) GetUidById(c *context.Context, userId string) (AccountUid, error) {
	desc, matchErr := MatchIdDescriptor(userId)
	if matchErr != nil {
		return "", matchErr
	}
	return s.model.GetUidById(desc.Name, userId)
}

func (s *AccountMgr) Login(c *context.Context, userId string, pwd *AccountPwd) error {
	uid, findErr := s.GetUidById(c, userId)
	if findErr != nil {
		return errors.New("invalid user id")
	}
	baseInfo, accountErr := s.model.GetAccountBaseInfo(uid)
	if accountErr != nil {
		return accountErr
	}
	desc, _ := MatchIdDescriptor(userId)
	if !baseInfo.Ids[desc.Name].Verified {
		return errors.New("account id not verified")
	}
	accountPwd, pwdErr := baseInfo.Password.GetPwd()
	if pwdErr == nil {
		userPwd, _ := pwd.GetPwd()
		if userPwd != accountPwd {
			return errors.New("invalid password")
		}
	} else {
		return pwdErr
	}
	ss := c.Input.CruSession
	ss.Set("LoginUser", uid.String())
	return nil
}

func (s *AccountMgr) LogoutSession(c *context.Context, sid string) error {
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
	curUser, _ := s.getOperateUid(c)
	fmt.Println(curUser.String() + " : " + ssUser.(string))
	if curUser.String() != ssUser.(string) {
		return errors.New("login user not match")
	}
	fmt.Println("Logged out")
	// ****
	ss.Flush()
	ss.SessionRelease(c.ResponseWriter)
	fmt.Println(ss.Get("LoginUser"))

	return nil
}

func (s *AccountMgr) ChangePwd(c *context.Context, oldPwd *AccountPwd, newPwd *AccountPwd) error {
	uid, err := s.getOperateUid(c)
	if err != nil {
		return err
	}
	accountBaseInfo, accountErr := s.model.GetAccountBaseInfo(uid)
	if accountErr != nil {
		return accountErr
	}
	oldPwdSaved, oldPwdSavedErr := accountBaseInfo.Password.GetPwd()
	oldPwdUser, oldPwdUserErr := oldPwd.GetPwd()

	if oldPwdSavedErr == nil && oldPwdUserErr == nil && oldPwdSaved == oldPwdUser {
		newPwdRaw, newPwdErr := newPwd.GetPwd()
		if newPwdErr == nil {
			accountBaseInfo.Password.SetEncryptedPwd(newPwdRaw)
			return s.model.UpdateAccountBaseInfo(uid, accountBaseInfo)
		}
	}
	return nil
}

func (s *AccountMgr) GetAccountBaseInfo(c *context.Context) (*AccountBaseInfo, error) {
	uid, err := s.getOperateUid(c)
	if err != nil {
		return nil, err
	}
	return s.model.GetAccountBaseInfo(uid)
}

func (s *AccountMgr) UpdateAccountBaseInfo(c *context.Context, baseInfo *AccountBaseInfo) error {
	uid, err := s.getOperateUid(c)
	if err != nil {
		return nil
	}
	return s.model.UpdateAccountBaseInfo(uid, baseInfo)
}

func (s *AccountMgr) GetProfiles(c *context.Context) (map[KeyName]string, error) {
	uid, err := s.getOperateUid(c)
	if err != nil {
		return nil, err
	}
	return s.model.GetProfiles(uid)
}

func (s *AccountMgr) UpdateProfiles(c *context.Context, profiles map[KeyName]string) error {
	uid, err := s.getOperateUid(c)
	if err != nil {
		return err
	}
	return s.model.UpdateProfiles(uid, profiles)
}

func (s *AccountMgr) GetOthers(c *context.Context) (map[KeyName]string, error) {
	uid, err := s.getOperateUid(c)
	if err != nil {
		return nil, err
	}
	return s.model.GetOthers(uid)
}

func (s *AccountMgr) UpdateOthers(c *context.Context, others map[KeyName]string) error {
	uid, err := s.getOperateUid(c)
	if err != nil {
		return err
	}
	return s.model.UpdateOthers(uid, others)

}

// Operation on other account, needs permission check
func (s *AccountMgr) OtherAccount(c *context.Context, uid string) *AccountMgr {
	ss := c.Input.CruSession
	ss.Set("OtherUser", uid)
	return s
}
