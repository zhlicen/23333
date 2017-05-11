package users

import (
	"23333/utils/web/beenh/beeaccount"
	"fmt"

	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	if AccountMgr.GetLoginUserID(c.Ctx) != nil {
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
	userID, getErr := AccountMgr.GetUserID(c.Ctx, username)
	if getErr != nil {
		fmt.Println(getErr.Error())
		c.Ctx.Redirect(302, "/login")
		return
	}
	userPwd := new(beeaccount.LoginPwd)
	userPwd.SetPwd("23333", password, userID.Uid, pwdEncryptorSalt)
	loginErr := AccountMgr.Login(c.Ctx, username, userPwd)

	if loginErr != nil {
		fmt.Println(loginErr.Error())
		c.Ctx.Redirect(302, "/login")
		return
	}
	accountInfo, _ := AccountMgr.CurrentAccount(c.Ctx).GetAccountBasicInfo()
	c.SetSession("UserName", accountInfo.LoginIDs[UserName].ID)
	c.Ctx.Redirect(302, "/")
}
