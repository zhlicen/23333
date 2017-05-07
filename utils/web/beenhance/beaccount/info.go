package beaccount

import (
	"errors"
	"fmt"
	"strings"
	"23333/utils/encrypt"
	"23333/utils/idgen"
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

func (accountPwd *AccountPwd) SetPwd(descriptor KeyName,
	pwd string, param interface{}, encryptor encrypt.Encryptor) error {
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
type AccountUid string

func NewAccountUid(uid string) AccountUid {
	return AccountUid(uid)
}

func (a AccountUid) String() string {
	return string(a)
}

func (a *AccountUid) SetVal(val string) {
	*a = AccountUid(val)
}

type AccountBaseInfo struct {
	Domain   string
	Group    string
	Uid      AccountUid
	Ids      map[IdName]AccountId
	Password AccountPwd
}

func NewAccountBaseInfo() *AccountBaseInfo {
	accountBaseInfo := new(AccountBaseInfo)
	accountBaseInfo.Ids = make(map[IdName]AccountId)
	return accountBaseInfo
}

func (a *AccountBaseInfo) GenRandomUid() (string, error) {
	keyGen := idgen.NewRandomIdGenerator(16, []byte(`1234567890abcdefghijklmnopqrstuvwxyz`)...)
	uid, keyErr := keyGen.Generate()
	if keyErr != nil {
		return "", keyErr
	}
	a.Uid.SetVal(uid)
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
	accountInfo.Ids = make(map[IdName]AccountId)
	accountInfo.OAuth2Id = make(map[KeyName]string)
	accountInfo.Profiles = make(map[KeyName]string)
	accountInfo.Others = make(map[KeyName]string)
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
