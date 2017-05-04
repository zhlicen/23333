package users

import (
	"fmt"

	"github.com/astaxie/beego"
)

type LogoutController struct {
	beego.Controller
}

func (c *LogoutController) Post() {
	c.StartSession()
	err := accountMgr.LogoutSession(c.Ctx, c.CruSession.SessionID())
	if err != nil {
		fmt.Println(err.Error())
		c.Ctx.Redirect(302, "/")
		return
	}
	c.Ctx.Redirect(302, "/login")
}
