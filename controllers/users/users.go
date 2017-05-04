package users

import (
	"23333/account"
	. "23333/models/users"
	"23333/utilities"
	"23333/utilities/verify"
	"net/mail"
)

var accountModel account.AccountModel
var accountMgr *account.AccountMgr
var pwdEncryptorSalt utilities.Encryptor
var accountVS *verify.SMTPVerifyService

func init() {
	accountModel = new(UserModel)
	accountMgr = account.NewAccountService("app", accountModel)
	pwdEncryptorSalt = utilities.NewSaultEncryptor("@!#!@", "12ws")
	// intialize verify service
	keyGen := utilities.NewRandomKeyGenerator(32)
	smtpConfig := verify.SMTPConfig{"smtp.sohu.com", "jjj_noreply@sohu.com", "232323",
		mail.Address{"23333", "jjj_noreply@sohu.com"}}
	accountVS = verify.NewSMTPVerifyService("views/mail.tpl", 60, keyGen, smtpConfig)
	accountVS.AddMailTplData("url", "https://localhost:8443/accountVerify")
}
