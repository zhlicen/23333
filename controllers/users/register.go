package users

import (
	"23333/account"
	"23333/utilities"

	"fmt"

	"github.com/astaxie/beego"
)

type RegisterController struct {
	beego.Controller
}

func (c *RegisterController) Get() {
	c.TplName = "register/index.tpl"
}

func (c *RegisterController) Post() {
	accountInfo := new(account.AccountInfo)
	username := c.GetString("username")
	password := c.GetString("password")
	encryptor := utilities.NewSaultEncryptor("@!#!@", "12ws")
	accountInfo.Password.SetPwd(password, username, encryptor)
	pwd, err := accountInfo.Password.GetPwd()
	if err == nil {
		fmt.Println("Password:" + pwd)
	}
	service := new(account.AccountService)
	service.Register(c.Ctx, accountInfo)
	c.TplName = "register/index.tpl"
}
