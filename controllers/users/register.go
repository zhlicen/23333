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
	if c.GetSession("LoginUser") != nil {
		fmt.Println("Logout current user to register")
		c.DelSession("LoginUser")
	}
	c.TplName = "register/index.tpl"
}

func (c *RegisterController) Post() {
	accountInfo := account.NewAccountInfo()
	username := c.GetString("username")
	password := c.GetString("password")
	mobile := c.GetString("mobile")
	email := c.GetString("email")
	accountInfo.Domain = "customer"
	accountInfo.Ids[account.UserName.Name] = username
	fmt.Println("UserName:" + username)
	accountInfo.Ids[account.Mobile.Name] = mobile
	fmt.Println("Mobile:" + mobile)
	accountInfo.Ids[account.Email.Name] = email
	fmt.Println("Email:" + email)

	accountInfo.Password.SetPwd(password, accountInfo.PrimaryId, pwdEncryptorSalt)
	pwd, err := accountInfo.Password.GetPwd()
	if err == nil {
		fmt.Println("Password:" + pwd)
	}

	regErr := accountService.Register(c.Ctx, accountInfo)
	if regErr != nil {
		fmt.Println(regErr.Error())
		c.Ctx.Redirect(302, "/register")
		return
	}
	c.Ctx.Redirect(302, "/login")
}
