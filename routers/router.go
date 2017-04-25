package routers

import (
	"23333/controllers"
	"23333/controllers/users"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &users.RegisterController{})
	beego.Router("/login", &users.LoginController{})
	beego.Router("/idCheck", &users.IdCheckController{})
}
