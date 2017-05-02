package account

import (
	"23333/utilities"
	"errors"

	"github.com/astaxie/beego"
	context "github.com/astaxie/beego/context"
)

type AccountOperation int

const (
	Register AccountOperation = iota
	Login
	LogoutSession
)

type sessionKeys struct {
	loginUser string
}

type AccountService struct {
	domain string
	model  AccountModel
	ssKeys sessionKeys
}

func NewAccountService(domain string, model AccountModel) *AccountService {
	if model == nil {
		return nil
	}
	keyGen := utilities.NewRandomKeyGenerator(16, []byte(`1234567890`)...)
	loginUser, _ := keyGen.Generate()
	return &AccountService{domain, model, sessionKeys{loginUser}}
}

func (s *AccountService) Register(c *context.Context, info *AccountInfo) error {
	valErr := info.Validate()
	if valErr != nil {
		return valErr
	}
	return s.model.Add(info)
}

func (s *AccountService) Login(c *context.Context, userId string, pwd *AccountPwd) (string, error) {
	accountInfo, findErr := s.model.FindById(userId)
	if findErr != nil {
		return "", errors.New("invalid user id")
	}
	accountPwd, pwdErr := accountInfo.Password.GetPwd()
	if pwdErr == nil {
		userPwd, _ := pwd.GetPwd()
		if userPwd != accountPwd {
			return "", errors.New("invalid password")
		}
	} else {
		return "", pwdErr
	}
	ss := c.Input.CruSession
	ss.Set(s.ssKeys.loginUser, accountInfo.Uid)
	return accountInfo.Uid, nil
}

func (s *AccountService) GetUidById(c *context.Context, userId string) (string, error) {
	accountInfo, error := s.model.FindById(userId)
	return accountInfo.Uid, error
}

func (s *AccountService) LogoutSession(c *context.Context, sid string) error {
	curSs := c.Input.CruSession
	ss, ssErr := beego.GlobalSessions.GetSessionStore(sid)
	if ssErr != nil {
		return nil
	}
	ssUser := ss.Get("LoginUser")
	if ssUser == nil {
		return errors.New("no login user with this session")
	}
	curUser := curSs.Get("LoginUser")
	if curUser != ssUser {
		return errors.New("login usr not match")
	}
	// ****
	ss.SessionRelease(c.ResponseWriter)
	return nil
}

func (s *AccountService) ChangePwd(c *context.Context, oldPwd *AccountPwd, newPwd *AccountPwd) error {
	return nil
}

func (s *AccountService) GetProfiles(c *context.Context) (map[KeyName]string, error) {
	return nil, nil
}

func (s *AccountService) UpdateProfiles(c *context.Context, profiles map[KeyName]string) error {
	return nil
}

// Operation on other account, needs permission check
func (s *AccountService) OtherAccount(c *context.Context, uid string) (error, *AccountService) {
	return nil, nil
}
