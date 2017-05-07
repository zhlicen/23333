package users

import (
	"23333/utils/web/beenhance/beaccount"
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
	Uid, getErr := accountMgr.GetUidById(c.Ctx, username)
	if getErr != nil {
		fmt.Println(getErr.Error())
		c.Ctx.Redirect(302, "/login")
		return
	}
	userPwd := new(beaccount.AccountPwd)
	userPwd.SetPwd("Password", password, Uid.String(), pwdEncryptorSalt)
	loginErr := accountMgr.Login(c.Ctx, username, userPwd)

	if loginErr != nil {
		fmt.Println(loginErr.Error())
		c.Ctx.Redirect(302, "/login")
		return
	}
	c.Ctx.Redirect(302, "/")
}
