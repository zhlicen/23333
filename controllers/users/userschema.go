package users

import (
	"23333/utils/web/beenh/beeaccount"
	"reflect"
)

const (
	UserName = "UserName"
	Email    = "Email"
	Mobile   = "Mobile"
	NickName = "NickName"
)

func initSchema() {
	accountSchema, _ := beeaccount.AddAccountSchema("23333")
	accountSchema.AddGroups("admin", "visitor", "customer")
	accountSchema.AddLoginIDSchema(UserName, false, false, true, `^[a-z]{1}[\w_]{3,15}$`)
	accountSchema.AddLoginIDSchema(Mobile, false, false, true, `^1[\d]{10}$`)
	accountSchema.AddLoginIDSchema(Email, false, false, true, `^([\w\.\_-]+)@([\w\.\_-]+)(\.[\w\.\_-]+)+$`)
	accountSchema.SetPasswordSchema(`^[\w+-_#*]{6, 16}$`)
	nickNameValidator := beeaccount.NewStringValidator(false, `^[\w_]{3,16}$`)
	accountSchema.AddOptionSchema(NickName, false, reflect.TypeOf(""), nickNameValidator)
	emailSchema, _ := accountSchema.GetLoginIDSchema(Email)
	emailSchema.UserData = "userdata"

}
