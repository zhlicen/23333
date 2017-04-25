package users

import (
	"23333/account"
	"fmt"

	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	if c.GetSession("LoginUser") != nil {
		c.Ctx.Redirect(302, "/")
		return
	}
	c.TplName = "login/index.tpl"
}

func (c *LoginController) Post() {
	username := c.GetString("username")
	password := c.GetString("password")
	fmt.Println("username:" + username)
	fmt.Println("password:" + password)
	accountInfo, getErr := accountService.GetAccountInfoById(c.Ctx, username)
	if getErr != nil {
		fmt.Println(getErr.Error())
		c.Ctx.Redirect(302, "/login")
		return
	}
	userPwd := new(account.AccountPwd)
	userPwd.SetPwd("Password", password, accountInfo.Uid, pwdEncryptorSalt)
	loginErr := accountService.Login(c.Ctx, username, userPwd)

	if loginErr != nil {
		fmt.Println(loginErr.Error())
		c.Ctx.Redirect(302, "/login")
		return
	}
	c.Ctx.Redirect(302, "/")
}
