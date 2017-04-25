package account

var Email *IdDescriptor
var Mobile *IdDescriptor
var UserName *IdDescriptor

func initIds() {
	UserName, _ = NewIdDescriptor("UserName", `^[a-z]{1}[\w_]{3,15}`, false, "User Name")
	Mobile, _ = NewIdDescriptor("Mobile", `^1[\d]{10}`, false, "Mobile Number")
	Email, _ = NewIdDescriptor("Email", `^([\w\.\_]{2,20})@(\w{1,}).([a-z]{2,4})$`, false, "Email xxx@xxx.xxx")
}
