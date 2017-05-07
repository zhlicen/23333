package users

import (
	"fmt"

	"github.com/astaxie/beego"
)

type IdCheckController struct {
	beego.Controller
}

func (c *IdCheckController) Get() {
	id := c.GetString("id")
	fmt.Println("id:" + id)
	_, getErr := accountMgr.GetUserId(c.Ctx, id)
	if getErr != nil {
		c.Data["json"] = map[string]interface{}{"code": 1, "message": getErr.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"code": 0, "message": "exist"}
	}
	c.ServeJSON()
}
