package beeaccount

import (
	"23333/utils/encrypt"
	"23333/utils/idgen"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
)

func initAccountInfo() {
	gob.Register(UserId{})
}

// LoginId id for login
// member:Id the id value string
// member:Verified if the id is verified
type LoginId struct {
	Id       string
	Verified bool
}

// NewLoginId constructor of LoginId
// param id is the id string
// param verified
// if param verified is not specified, true will be used default
// return the new id constructed
func NewLoginId(id string, verified ...bool) LoginId {
	if verified == nil {
		return LoginId{id, true}
	}
	return LoginId{id, verified[0]}
}

// LoginPwd login password
// member:pwd the encrypted password
type LoginPwd struct {
	pwd string
}

// SetPwd set the password and encrypt methods
// descriptor is the key descriptor name for format checking
// pwd is the password of plain text
// param is the encrypt param of encryptor
// encryptor is the encrypt method of this password
// returns error if there's any
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

// GetPwd get the encrypted password
// an error will be returned if no password is set
func (LoginPwd *LoginPwd) GetPwd() (string, error) {
	if LoginPwd.pwd == "" {
		return "", errors.New("pwd not exist")
	}
	return LoginPwd.pwd, nil
}

// SetEncryptedPwd set encrypted password directly
func (LoginPwd *LoginPwd) SetEncryptedPwd(pwd string) {
	LoginPwd.pwd = pwd
}

// AccountStatus status of an account
// member:Activated user can not login if the account is not activated
// member:Locked user can not login if the account is not locked
// member:Sessions beego session ids of this account
type AccountStatus struct {
	Activated bool
	Locked    bool
	Sessions  []string
}

// UserId user identity
// member:Domain domain of this account
// member:Group group of this account
// member:Uid unique id in this domian, should be used as pk in db
type UserId struct {
	Domain string
	Group  string
	Uid    string
}

// AccountBasicInfo account basic info
// parent:UserId
// member:LoginIds ids for login, id should identify one account
// member:Password password for account
type AccountBasicInfo struct {
	UserId
	LoginIds map[IdName]LoginId
	Password LoginPwd
}

// NewAccountBasicInfo constructor of AccountBasicInfo
func NewAccountBasicInfo() *AccountBasicInfo {
	accountBasicInfo := new(AccountBasicInfo)
	accountBasicInfo.LoginIds = make(map[IdName]LoginId)
	return accountBasicInfo
}

// GenRandomUid generate a random uid
func (a *AccountBasicInfo) GenRandomUid() (string, error) {
	keyGen := idgen.NewRandomIdGenerator(16, []byte(`1234567890abcdefghijklmnopqrstuvwxyz`)...)
	uid, keyErr := keyGen.Generate()
	if keyErr != nil {
		return "", keyErr
	}
	a.Uid = uid
	return uid, nil
}

// AccountInfo account info
// parent:AccountBasicInfo
// member:Profiles profiles key-value pair
// member:Others other information
// member:Status status of account
type AccountInfo struct {
	AccountBasicInfo
	Profiles map[KeyName]string
	Others   map[KeyName]string
	Status   AccountStatus
}

// NewAccountInfo constructor of AccountInfo
func NewAccountInfo() *AccountInfo {
	accountInfo := new(AccountInfo)
	accountInfo.LoginIds = make(map[IdName]LoginId)
	accountInfo.Profiles = make(map[KeyName]string)
	accountInfo.Others = make(map[KeyName]string)
	return accountInfo
}

// Validate check if the account info is valid with descriptors defined
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
	}

	return nil
}
