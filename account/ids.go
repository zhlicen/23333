package account

var Email *IdDescriptor
var Mobile *IdDescriptor

func initIds() {
	Mobile, _ = NewIdDescriptor("Mobile", "number-number", "Mobile Number(xxx)xxxx")
	Email, _ = NewIdDescriptor("Email", "xxx@xxx.xxx", "Email xxx@xxx.xxx")
}
