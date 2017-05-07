package beaccount

import (
	"23333/utils/encrypt"
	"23333/utils/idgen"
	"errors"
	"fmt"
	"strings"
)

type LoginId struct {
	Id       string
	Verified bool
}

func NewLoginId(id string, verified ...bool) LoginId {
	if verified == nil {
		return LoginId{id, true}
	}
	return LoginId{id, verified[0]}
}

type LoginPwd struct {
	pwd string
}

func (LoginPwd *LoginPwd) SetPwd(descriptor KeyName,
	pwd string, param interface{}, encryptor encrypt.Encryptor) error {
	var err error
	desc, _ := GetKeyDescriptor(descriptor)
	if desc != nil && desc.Validate(pwd) {
		return errors.New("invalid pwd:" + desc.Description)
	}
	LoginPwd.pwd, err = encryptor.Encrypt(pwd, param)
	return err
}

func (LoginPwd *LoginPwd) GetPwd() (string, error) {
	if LoginPwd.pwd == "" {
		return "", errors.New("pwd not exist")
	}
	return LoginPwd.pwd, nil
}

func (LoginPwd *LoginPwd) SetEncryptedPwd(pwd string) {
	LoginPwd.pwd = pwd
}

type AccountStatus struct {
	Activated  bool
	Locked     bool
	LockExpire string
	Sessions   []string
}

type UserId struct {
	Domain string
	Group  string
	Uid    string
}

type AccountBaseInfo struct {
	UserId
	LoginIds map[IdName]LoginId
	Password LoginPwd
}

func NewAccountBaseInfo() *AccountBaseInfo {
	accountBaseInfo := new(AccountBaseInfo)
	accountBaseInfo.LoginIds = make(map[IdName]LoginId)
	return accountBaseInfo
}

func (a *AccountBaseInfo) GenRandomUid() (string, error) {
	keyGen := idgen.NewRandomIdGenerator(16, []byte(`1234567890abcdefghijklmnopqrstuvwxyz`)...)
	uid, keyErr := keyGen.Generate()
	if keyErr != nil {
		return "", keyErr
	}
	a.Uid = uid
	return uid, nil
}

type AccountInfo struct {
	AccountBaseInfo
	OAuth2Id map[KeyName]string
	Profiles map[KeyName]string
	Others   map[KeyName]string
	Status   AccountStatus
}

func NewAccountInfo() *AccountInfo {
	accountInfo := new(AccountInfo)
	accountInfo.LoginIds = make(map[IdName]LoginId)
	accountInfo.OAuth2Id = make(map[KeyName]string)
	accountInfo.Profiles = make(map[KeyName]string)
	accountInfo.Others = make(map[KeyName]string)
	return accountInfo
}

func (accountInfo *AccountInfo) Validate() error {
	// validate ids
	validIdCount := 0
	for k, v := range accountInfo.LoginIds {
		descriptor, err := GetIdDescriptor(k)
		if !descriptor.CaseSensitive {
			v.Id = strings.ToLower(v.Id)
			accountInfo.LoginIds[k] = v
		}
		fmt.Println("Checking " + string(k))
		if err == nil && !descriptor.Validate(v.Id) {
			return errors.New(string(k) + " do not match format, " + descriptor.Description)
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
		fmt.Println("Checking " + string(k))
		if err == nil && !descriptor.Validate(v) {
			return errors.New(string(k) + " do not match format, " + descriptor.Description)

		}
		fmt.Println("---OK")
	}

	for k, v := range accountInfo.Others {
		descriptor, err := GetKeyDescriptor(k)
		if !descriptor.CaseSensitive {
			v = strings.ToLower(v)
			accountInfo.Others[k] = v
		}
		fmt.Println("Checking " + string(k))
		if err == nil && !descriptor.Validate(v) {
			return errors.New(string(k) + " do not match format, " + descriptor.Description)

		}
		fmt.Println("---OK")
	}

	for k, v := range accountInfo.OAuth2Id {
		descriptor, err := GetKeyDescriptor(k)
		if !descriptor.CaseSensitive {
			v = strings.ToLower(v)
			accountInfo.OAuth2Id[k] = v
		}
		fmt.Println("Checking " + string(k))
		if err == nil && !descriptor.Validate(v) {
			return errors.New(string(k) + " do not match format, " + descriptor.Description)
		}
		fmt.Println("---OK")
	}

	return nil
}
