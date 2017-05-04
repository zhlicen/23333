package users

import (
	"net/mail"

	"github.com/astaxie/beego"
)

type SendVerifyController struct {
	beego.Controller
}

func (c *SendVerifyController) Get() {
	username := c.GetString("username")
	email := c.GetString("email")
	uid := c.GetString("uid")
	to := mail.Address{username, email}
	data := make(map[string]string)
	data["uid"] = uid
	data["id"] = email
	tokenErr := accountVS.SendToken(uid, to, "Account Verification", data)
	if tokenErr != nil {
		c.Data["json"] = map[string]interface{}{"code": 1, "message": tokenErr.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"code": 0, "message": "Email sent"}
	}

	c.ServeJSON()
}
