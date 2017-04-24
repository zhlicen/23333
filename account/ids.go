package account

var Email *IdDescriptor
var Mobile *IdDescriptor
var UserName *IdDescriptor

func initIds() {
	UserName, _ = NewIdDescriptor("UserName", `^[a-z]{1}[\w_]{3,15}`, "User Name")
	Mobile, _ = NewIdDescriptor("Mobile", `^1[\d]{10}`, "Mobile Number")
	Email, _ = NewIdDescriptor("Email", `^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, "Email xxx@xxx.xxx")
}
