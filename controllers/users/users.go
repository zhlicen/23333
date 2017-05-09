package users

import (
	. "23333/models/users"
	"23333/utils/encrypt"
	"23333/utils/idgen"
	"23333/utils/web/beenh/beeaccount"
	"23333/utils/web/verify"
	"net/mail"
)

var accountModel beeaccount.AccountModel
var AccountMgr *beeaccount.AccountMgr
var pwdEncryptorSalt encrypt.Encryptor
var accountVS *verify.SMTPVerifyService

func init() {
	accountModel = new(UserModel)
	AccountMgr = beeaccount.NewAccountMgr("23333", accountModel, nil)
	pwdEncryptorSalt = encrypt.NewSaultEncryptor("@!#!@", "12ws")
	// intialize verify service
	idGen := idgen.NewRandomIdGenerator(32)
	smtpConfig := verify.SMTPConfig{"smtp.sohu.com", "jjj_noreply@sohu.com", "232323",
		mail.Address{"23333", "jjj_noreply@sohu.com"}}
	accountVS = verify.NewSMTPVerifyService("views/mail.tpl", 60, idGen, smtpConfig)
	accountVS.AddMailTplData("url", "https://localhost:8443/accountVerify")
	initSchema()
}
