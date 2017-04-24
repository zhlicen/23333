package users

import (
	"23333/account"
	. "23333/models/users"
	"23333/utilities"
)

var accountModel account.AccountModel
var accountService *account.AccountService
var pwdEncryptorSalt utilities.Encryptor

func init() {
	accountModel = new(UserModel)
	accountService = account.NewAccountService(accountModel)
	pwdEncryptorSalt = utilities.NewSaultEncryptor("@!#!@", "12ws")
}
