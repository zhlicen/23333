package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	if c.GetSession("LoginUser") == nil {
		c.Ctx.Redirect(302, "/login")
		return
	}
	c.Data["LoginUser"] = c.GetSession("LoginUser")
	c.TplName = "index.tpl"
}
