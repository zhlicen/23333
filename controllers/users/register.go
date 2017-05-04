package users

import (
	"23333/account"
	"fmt"

	"github.com/astaxie/beego"
)

type RegisterController struct {
	beego.Controller
}

func (c *RegisterController) Get() {
	c.StartSession()
	err := accountMgr.LogoutSession(c.Ctx, c.CruSession.SessionID())
	if err != nil {
		fmt.Println(err.Error())
	}
	c.TplName = "register/index.tpl"
}

func (c *RegisterController) Post() {
	accountInfo := account.NewAccountInfo()
	accountInfo.GenRandomUid()
	username := c.GetString("username")
	password := c.GetString("password")
	mobile := c.GetString("mobile")
	email := c.GetString("email")
	accountInfo.Domain = beego.BConfig.AppName
	accountInfo.Ids[account.UserName.Name] = account.NewAccountId(username)
	fmt.Println("UserName:" + username)
	accountInfo.Ids[account.Mobile.Name] = account.NewAccountId(mobile)
	fmt.Println("Mobile:" + mobile)
	accountInfo.Ids[account.Email.Name] = account.NewAccountId(email, false)
	fmt.Println("Email:" + email)

	accountInfo.Password.SetPwd("Password", password, accountInfo.Uid.String(), pwdEncryptorSalt)
	pwd, err := accountInfo.Password.GetPwd()
	if err == nil {
		fmt.Println("Password:" + pwd)
	} else {
		fmt.Println(err)
	}
	regErr := accountMgr.Register(c.Ctx, accountInfo)
	if regErr != nil {
		fmt.Println(regErr.Error())
		c.Ctx.Redirect(302, "/register")
		return
	}
	queryString := "uid=" + accountInfo.Uid.String() + "&email=" + email + "&username=" + username
	c.Ctx.Redirect(302, "/sendVerify?"+queryString)
}
