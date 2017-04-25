package account

import (
	"errors"

	context "github.com/astaxie/beego/context"
)

type AccountService struct {
	domain string
	model  AccountModel
}

func NewAccountService(domain string, model AccountModel) *AccountService {
	if model == nil {
		return nil
	}
	return &AccountService{domain, model}
}

func (s *AccountService) Register(c *context.Context, info *AccountInfo) error {
	valErr := info.Validate()
	if valErr != nil {
		return valErr
	}
	return s.model.Add(info)
}

func (s *AccountService) Login(c *context.Context, userId string, pwd *AccountPwd) error {
	accountInfo, findErr := s.model.FindById(userId)
	if findErr != nil {
		return errors.New("invalid user id")
	}
	accountPwd, pwdErr := accountInfo.Password.GetPwd()
	if pwdErr == nil {
		userPwd, _ := pwd.GetPwd()
		if userPwd != accountPwd {
			return errors.New("invalid password")
		}
	}
	ss := c.Input.CruSession
	ss.Set("LoginUser", accountInfo.Uid)
	return nil
}

func (s *AccountService) GetAccountInfoById(c *context.Context, userId string) (*AccountInfo, error) {
	return s.model.FindById(userId)
}
