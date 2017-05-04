package users

import (
	"23333/utilities/verify"

	"github.com/astaxie/beego"
)

type AccountVerifyController struct {
	beego.Controller
}

func (c *AccountVerifyController) Get() {
	id := c.GetString("id")
	key := c.GetString("key")
	token := c.GetString("token")
	var code int
	var message string = "account verified!"
	verifier := verify.NewVerifier(accountVS, key, token)
	err := accountMgr.VerifyId(c.Ctx, verifier, id)
	if err != nil {
		code = 1
		message = err.Error()
	}
	c.Data["json"] = map[string]interface{}{"code": code, "message": message}
	c.ServeJSON()
}
