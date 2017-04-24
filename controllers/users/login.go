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
	pwd := new(account.AccountPwd)
	pwd.SetPwd(password, accountInfo.PrimaryId, pwdEncryptorSalt)
	accountPwd, pwdErr := accountInfo.Password.GetPwd()
	if pwdErr != nil {
		fmt.Println("invalid password")
		c.Ctx.Redirect(302, "/login")
		return
	}
	userPwd, _ := pwd.GetPwd()
	if userPwd != accountPwd {
		fmt.Println("invalid password")
		c.Ctx.Redirect(302, "/login")
		return
	}
	c.SetSession("LoginUser", username)
	c.Ctx.Redirect(302, "/")
}
