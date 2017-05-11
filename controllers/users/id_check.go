package users

import (
	"fmt"

	"github.com/astaxie/beego"
)

type IDCheckController struct {
	beego.Controller
}

func (c *IDCheckController) Get() {
	id := c.GetString("id")
	fmt.Println("id:" + id)
	_, getErr := AccountMgr.GetUserID(c.Ctx, id)
	if getErr != nil {
		c.Data["json"] = map[string]interface{}{"code": 1, "message": getErr.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"code": 0, "message": "exist"}
	}
	c.ServeJSON()
}
