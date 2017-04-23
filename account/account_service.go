package account

import context "github.com/astaxie/beego/context"

type AccountService struct {
}

func (service *AccountService) Register(context *context.Context, info *AccountInfo) error {
	return nil
}
