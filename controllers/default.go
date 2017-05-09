package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	userName := c.GetSession("UserName")
	if userName == nil {
		c.Ctx.Redirect(302, "/login")
		return
	}
	c.Data["LoginUser"] = userName.(string)
	c.TplName = "index.tpl"
}
