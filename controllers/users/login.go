package users

import (
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
	c.SetSession("LoginUser", username)
	c.Ctx.Redirect(302, "/")
}
