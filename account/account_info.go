package account

import (
	"23333/utilities"
	"errors"
	"fmt"
	"strings"
)

type AccountId struct {
	Id       string
	Verified bool
}

func NewAccountId(id string, verified ...bool) AccountId {
	if verified == nil {
		return AccountId{id, true}
	}
	return AccountId{id, verified[0]}
}

type AccountPwd struct {
	pwd string
}

func (accountPwd *AccountPwd) SetPwd(descriptor string,
	pwd string, param interface{}, encryptor utilities.Encryptor) error {
	var err error
	desc, _ := GetKeyDescriptor(descriptor)
	if desc != nil && desc.Validate(pwd) {
		return errors.New("invalid pwd:" + desc.Description)
	}
	accountPwd.pwd, err = encryptor.Encrypt(pwd, param)
	return err
}

func (accountPwd *AccountPwd) GetPwd() (string, error) {
	if accountPwd.pwd == "" {
		return "", errors.New("pwd not exist")
	}
	return accountPwd.pwd, nil
}

func (accountPwd *AccountPwd) SetEncryptedPwd(pwd string) {
	accountPwd.pwd = pwd
}

type AccountStatus struct {
	Activated  bool
	Locked     bool
	LockExpire string
}

type AccountInfo struct {
	Domain   string
	Group    string
	Uid      string
	Password AccountPwd
	Ids      map[string]AccountId
	OAuth2Id map[string]string
	Profiles map[string]string
	Others   map[string]string
	Status   AccountStatus
}

func NewAccountInfo() *AccountInfo {
	accountInfo := new(AccountInfo)
	keyGen := utilities.NewRandomKeyGenerator(16, []byte(`1234567890abcdefghijklmnopqrstuvwxyz`)...)
	var keyErr error
	accountInfo.Uid, keyErr = keyGen.Generate()
	if keyErr != nil {
		return nil
	}
	accountInfo.Ids = make(map[string]AccountId)
	accountInfo.OAuth2Id = make(map[string]string)
	accountInfo.Profiles = make(map[string]string)
	accountInfo.Others = make(map[string]string)
	return accountInfo
}

func (accountInfo *AccountInfo) Validate() error {
	// validate ids
	validIdCount := 0
	for k, v := range accountInfo.Ids {
		descriptor, err := GetIdDescriptor(k)
		if !descriptor.CaseSensitive {
			v.Id = strings.ToLower(v.Id)
			accountInfo.Ids[k] = v
		}
		fmt.Println("Checking " + k)
		if err == nil && !descriptor.Validate(v.Id) {
			return errors.New(k + " do not match format, " + descriptor.Description)
		}
		fmt.Println("---OK")
		validIdCount++
	}
	if validIdCount == 0 {
		return errors.New("no valid id")
	}

	for k, v := range accountInfo.Profiles {
		descriptor, err := GetKeyDescriptor(k)
		if !descriptor.CaseSensitive {
			v = strings.ToLower(v)
			accountInfo.Profiles[k] = v
		}
		fmt.Println("Checking " + k)
		if err == nil && !descriptor.Validate(v) {
			return errors.New(k + " do not match format, " + descriptor.Description)

		}
		fmt.Println("---OK")
	}

	for k, v := range accountInfo.Others {
		descriptor, err := GetKeyDescriptor(k)
		if !descriptor.CaseSensitive {
			v = strings.ToLower(v)
			accountInfo.Others[k] = v
		}
		fmt.Println("Checking " + k)
		if err == nil && !descriptor.Validate(v) {
			return errors.New(k + " do not match format, " + descriptor.Description)

		}
		fmt.Println("---OK")
	}

	for k, v := range accountInfo.OAuth2Id {
		descriptor, err := GetKeyDescriptor(k)
		if !descriptor.CaseSensitive {
			v = strings.ToLower(v)
			accountInfo.OAuth2Id[k] = v
		}
		fmt.Println("Checking " + k)
		if err == nil && !descriptor.Validate(v) {
			return errors.New(k + " do not match format, " + descriptor.Description)
		}
		fmt.Println("---OK")
	}

	return nil
}
