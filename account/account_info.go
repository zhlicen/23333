package account

import (
	"23333/utilities"
	"errors"
	"fmt"
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
		return "", errors.New("pwd not exist")
	}
	return accountPwd.pwd, nil
}

func (accountPwd *AccountPwd) SetRawPwd(pwd string) {
	accountPwd.pwd = pwd
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

func NewAccountInfo() *AccountInfo {
	accountInfo := new(AccountInfo)
	keyGen := utilities.NewRandomKeyGenerator(16, []byte(`1234567890abcdefghijklmnopqrstuvwxyz`)...)
	var keyErr error
	accountInfo.PrimaryId, keyErr = keyGen.Generate()
	if keyErr != nil {
		return nil
	}
	accountInfo.Ids = make(map[string]string)
	accountInfo.OAuth2Id = make(map[string]string)
	accountInfo.Profiles = make(map[string]string)
	accountInfo.FixedProfiles = make(map[string]string)
	accountInfo.Others = make(map[string]string)
	return accountInfo
}

func (accountInfo *AccountInfo) Validate() error {
	// validate ids
	for k, v := range accountInfo.Ids {
		descriptor, err := GlobalIdDescriptorRegistry.Get(k)
		fmt.Print("Checking " + k)
		if err != nil {
			return errors.New("can not recognize id type " + k)
		}
		if !descriptor.Validate(v) {
			return errors.New(k + " do not match format, " + descriptor.Description)
		}
		fmt.Println("---OK")
	}
	return nil
}
