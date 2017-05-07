package users

import (
	"23333/utils/web/beenh/beeaccount"
	"fmt"

	"github.com/astaxie/beego"
)

type RegisterController struct {
	beego.Controller
}

func (c *RegisterController) Get() {
	c.StartSession()
	err := AccountMgr.CurrentAccount(c.Ctx).LogoutSession(c.CruSession.SessionID())
	if err != nil {
		fmt.Println(err.Error())
	}
	c.TplName = "register/index.tpl"
}

func (c *RegisterController) Post() {
	accountInfo := beeaccount.NewAccountInfo()
	accountInfo.GenRandomUid()
	username := c.GetString("username")
	password := c.GetString("password")
	mobile := c.GetString("mobile")
	email := c.GetString("email")
	accountInfo.Domain = beego.BConfig.AppName
	accountInfo.LoginIds[beeaccount.UserName.Name] = beeaccount.NewLoginId(username)
	fmt.Println("UserName:" + username)
	accountInfo.LoginIds[beeaccount.Mobile.Name] = beeaccount.NewLoginId(mobile)
	fmt.Println("Mobile:" + mobile)
	accountInfo.LoginIds[beeaccount.Email.Name] = beeaccount.NewLoginId(email, false)
	fmt.Println("Email:" + email)

	accountInfo.Password.SetPwd("Password", password, accountInfo.Uid, pwdEncryptorSalt)
	pwd, err := accountInfo.Password.GetPwd()
	if err == nil {
		fmt.Println("Password:" + pwd)
	} else {
		fmt.Println(err)
	}
	regErr := AccountMgr.Register(c.Ctx, accountInfo)
	if regErr != nil {
		fmt.Println(regErr.Error())
		c.Ctx.Redirect(302, "/register")
		return
	}
	queryString := "uid=" + accountInfo.Uid + "&email=" + email + "&username=" + username
	c.Ctx.Redirect(302, "/sendVerify?"+queryString)
}
