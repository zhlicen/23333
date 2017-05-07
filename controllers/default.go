package controllers

import (
	"23333/controllers/users"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	userId := users.AccountMgr.GetLoginUserId(c.Ctx)
	if userId == nil {
		c.Ctx.Redirect(302, "/login")
		return
	}
	c.Data["LoginUser"] = userId.Uid
	c.TplName = "index.tpl"
}
