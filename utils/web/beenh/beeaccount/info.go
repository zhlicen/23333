package beeaccount

import (
	"23333/utils/encrypt"
	"23333/utils/idgen"
	"encoding/gob"
	"errors"
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
func (LoginPwd *LoginPwd) SetPwd(domain string,
	pwd string, param interface{}, encryptor encrypt.Encryptor) error {
	var err error
	accountSchema, schemaErr := GetAccountSchema(domain)
	if schemaErr != nil {
		return errors.New("no schema for domian " + domain)
	}
	pwdSchema := accountSchema.GetPasswordSchema()
	if pwdSchema == nil {
		return errors.New("no password schema specified")
	}
	if pwdSchema.Validator.Validate(pwd) {
		return errors.New("invalid pwd format")
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
	LoginIds map[string]LoginId
	Password LoginPwd
}

// NewAccountBasicInfo constructor of AccountBasicInfo
func NewAccountBasicInfo(domain string) *AccountBasicInfo {
	accountBasicInfo := new(AccountBasicInfo)
	accountBasicInfo.Domain = domain
	accountBasicInfo.LoginIds = make(map[string]LoginId)
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
	Profiles map[string]interface{}
	Others   map[string]interface{}
	Status   AccountStatus
}

// NewAccountInfo constructor of AccountInfo
func NewAccountInfo(domain string) (*AccountInfo, error) {
	_, err := GetAccountSchema(domain)
	if err != nil {
		return nil, err
	}
	accountInfo := new(AccountInfo)
	accountInfo.Domain = domain
	accountInfo.LoginIds = make(map[string]LoginId)
	accountInfo.Profiles = make(map[string]interface{})
	accountInfo.Others = make(map[string]interface{})
	return accountInfo, nil
}

func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for _, mapOne := range maps {
		for k, v := range mapOne {
			merged[k] = v
		}
	}
	return merged
}

// Validate check if the account info is valid with descriptors defined
func (accountInfo *AccountInfo) Validate() error {
	accountSchema, _ := GetAccountSchema(accountInfo.Domain)
	if accountSchema == nil {
		return errors.New("schema undefined for domain " + accountInfo.Domain)
	}

	// Group
	if !accountSchema.IsGroupExist(accountInfo.Group) {
		return errors.New("unknown group " + accountInfo.Group)
	}

	// UserId
	if accountInfo.Uid == "" {
		return errors.New("uid can not be empty")
	}

	// LoginIds
	if len(accountInfo.LoginIds) == 0 {
		return errors.New("should have at least one login id")
	}
	requiredIds := accountSchema.getRequiredLogIds()
	for _, requiredId := range requiredIds {
		if _, ok := accountInfo.LoginIds[requiredId]; !ok {
			return errors.New("login id:" + requiredId + " is required but not specified")
		}
	}
	for k, v := range accountInfo.LoginIds {
		loginIdSchema, _ := accountSchema.GetLoginIdSchema(k)
		if loginIdSchema == nil {
			return errors.New("login id schema for " + k + " is not defined")
		} else {
			if !loginIdSchema.NeedVerified {
				v.Verified = true
				// accountInfo.LoginIds[k] = v
			}
			if !loginIdSchema.Validator.Validate(v.Id) {
				return errors.New("invalid format of login id " + k + ":" + v.Id)
			}
		}
	}

	// options
	optionsMap := mergeMaps(accountInfo.Profiles, accountInfo.Others)
	requiredOptions := accountSchema.getRequiredOptions()
	for _, requiredOption := range requiredOptions {
		if _, ok := optionsMap[requiredOption]; !ok {
			return errors.New("option:" + requiredOption + " is required but not specified")
		}
	}

	for k, v := range optionsMap {
		optionSchema, _ := accountSchema.GetOptionSchema(k)
		if optionSchema == nil {
			return errors.New("option schema for " + k + " is not defined")
		} else {
			if !optionSchema.Validator.Validate(v) {
				return errors.New("invalid format of option " + k)
			}
		}
	}

	return nil
}
