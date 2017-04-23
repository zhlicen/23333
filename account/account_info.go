package account

import (
	"23333/utilities"
	"errors"
)

type AccountPwd struct {
	pwd string
}

func (accountPwd *AccountPwd) SetPwd(pwd string, param interface{}, encryptor utilities.Encryptor) error {
	var err error
	accountPwd.pwd, err = encryptor.Encrypt(pwd, param)
	return err
}

func (accountPwd *AccountPwd) GetPwd() (string, error) {
	if accountPwd.pwd == "" {
		return "", errors.New("Pwd Not Exist")
	}
	return accountPwd.pwd, nil
}

type AccountStatus struct {
	Activated  bool
	Locked     bool
	LockExpire string
}

type AccountInfo struct {
	Domain        string
	Group         string
	PrimaryId     string
	Ids           map[string]string
	OAuth2Id      map[string]string
	Password      AccountPwd
	Profiles      map[string]string
	FixedProfiles map[string]string
	Status        AccountStatus
	Others        map[string]string
}

func (accountInfo *AccountInfo) Validate() error {
	return nil
}
