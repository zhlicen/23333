package beeaccount

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
)

// UserDataBase type for user to store data
type UserDataBase struct {
	UserData interface{}
}

func initAccountSchema() {
	if AccountSchemaMgr == nil {
		AccountSchemaMgr = new(accountSchemaMgr)
		AccountSchemaMgr.schemas = make(map[string]*AccountSchema)
	}
}

// loginIdSchema login id schema
// member:isRequired if the id is required
// member:needVerfied if the id needs be verified
// member:validator validator for id string
type LoginIdSchema struct {
	UserDataBase
	Name         string
	IsRequired   bool
	NeedVerified bool
	Validator    *stringValidator
}

// Validator for options(profiles and others)
type OptionValidator interface {
	Validate(option interface{}) bool
}

type stringValidator struct {
	ingoreCase bool
	patten     string
}

// NewStringValidator constructor of stringValidator
func NewStringValidator(ignoreCase bool, patten string) *stringValidator {
	return &stringValidator{ignoreCase, patten}
}

func (s *stringValidator) Validate(option interface{}) bool {
	if value, ok := option.(string); ok {
		if s.ingoreCase {
			value = strings.ToLower(value)
		}
		matched, _ := regexp.MatchString(s.patten, value)
		return matched
	}
	return false
}

type intValidator struct {
	min int
	max int
}

// NewIntValidator constructor if intValidator
// parems: min and max of the int value
func NewIntValidator(minVal int, maxVal int) *intValidator {
	return &intValidator{minVal, maxVal}
}

func (i *intValidator) Validate(option interface{}) bool {
	if value, ok := option.(int); ok {
		if value >= i.min && value <= i.max {
			return true
		}
	}
	return false
}

// OptionSchema schema for profiles and others
type OptionSchema struct {
	UserDataBase
	Name       string
	IsRequired bool
	OptionType reflect.Type
	Validator  OptionValidator
}

// PasswordSchema
type PasswordSchema struct {
	UserDataBase
	Validator *stringValidator
}

// AccountSchema account schema
type AccountSchema struct {
	UserDataBase
	domain         string
	groups         map[string]interface{}
	passwordSchema *PasswordSchema
	loginIdSchemas map[string]*LoginIdSchema
	optionSchemas  map[string]*OptionSchema
}

// GetDomain get domain of this schema
func (a *AccountSchema) GetDomain() string {
	return a.domain
}

// IsGroupExist if the group is exist in schema
func (a *AccountSchema) IsGroupExist(group string) bool {
	_, exist := a.groups[group]
	return exist
}

// AddGroups add groups definitions to schema
// error returns if the group is exist
func (a *AccountSchema) AddGroups(groups ...string) error {
	for _, group := range groups {
		if a.IsGroupExist(group) {
			return errors.New("group " + group + " exist!")
		}
		a.groups[group] = nil
	}
	return nil
}

// SetGroupData set group data
func (a *AccountSchema) SetGroupData(group string, data interface{}) error {
	if !a.IsGroupExist(group) {
		return errors.New("group " + group + " not exist!")
	}
	a.groups[group] = data
	return nil
}

// GetGroupData get group data
func (a *AccountSchema) GetGroupData(group string) (interface{}, error) {
	if data, ok := a.groups[group]; ok {
		return data, nil
	}
	return nil, errors.New("group not exist")
}

// SetPasswordSchema set password schema
func (a *AccountSchema) SetPasswordSchema(patten string, userData ...interface{}) {
	newSchema := new(PasswordSchema)
	newSchema.Validator = NewStringValidator(false, patten)
	if len(userData) > 0 {
		newSchema.UserData = userData[0]
	}
	a.passwordSchema = newSchema
}

// GetPasswordSchema set password schema
func (a *AccountSchema) GetPasswordSchema() *PasswordSchema {
	return a.passwordSchema
}

// AddLoginIdSchema add login id schema
func (a *AccountSchema) AddLoginIdSchema(name string, isRequired bool, needVerified bool,
	ignoreCase bool, patten string, userData ...interface{}) error {
	if _, ok := a.loginIdSchemas[name]; ok {
		return errors.New("schema " + name + " exist")
	}
	newSchema := new(LoginIdSchema)
	newSchema.Name = name
	newSchema.IsRequired = isRequired
	newSchema.NeedVerified = needVerified
	newSchema.Validator = NewStringValidator(ignoreCase, patten)
	if len(userData) > 0 {
		newSchema.UserData = userData[0]
	}
	a.loginIdSchemas[name] = newSchema
	return nil
}

func (a *AccountSchema) getRequiredLogIds() []string {
	var loginIdKeys []string
	for _, v := range a.loginIdSchemas {
		if v.IsRequired {
			loginIdKeys = append(loginIdKeys, v.Name)
		}
	}
	return loginIdKeys
}

// GetLoginIdSchema get login id schema
func (a *AccountSchema) GetLoginIdSchema(name string) (*LoginIdSchema, error) {
	if v, ok := a.loginIdSchemas[name]; ok {
		return v, nil
	}
	return nil, errors.New("schema not found")
}

// MatchLoginId get login id name
func (a *AccountSchema) MatchLoginId(id string) (string, error) {
	for k, v := range a.loginIdSchemas {
		if v.Validator.Validate(id) {
			return k, nil
		}
	}
	return "", errors.New("login id not found")
}

// AddOptionSchema add option schema
func (a *AccountSchema) AddOptionSchema(name string, isRequired bool,
	optionType reflect.Type, validator OptionValidator, userData ...interface{}) error {
	if _, ok := a.optionSchemas[name]; ok {
		return errors.New("schema " + name + " exist")
	}
	newSchema := new(OptionSchema)
	newSchema.Name = name
	newSchema.IsRequired = isRequired
	newSchema.OptionType = optionType
	newSchema.Validator = validator
	if len(userData) > 0 {
		newSchema.UserData = userData[0]
	}
	a.optionSchemas[name] = newSchema
	return nil
}

// GetOptionSchema get option schema
func (a *AccountSchema) GetOptionSchema(name string) (*OptionSchema, error) {
	if v, ok := a.optionSchemas[name]; ok {
		return v, nil
	}
	return nil, errors.New("schema not found")
}

func (a *AccountSchema) getRequiredOptions() []string {
	var options []string
	for _, v := range a.optionSchemas {
		if v.IsRequired {
			options = append(options, v.Name)
		}
	}
	return options
}

var AccountSchemaMgr *accountSchemaMgr

// accountSchemaMgr account schema manager
type accountSchemaMgr struct {
	schemas map[string]*AccountSchema
}

// AddAccountSchema add
func (a *accountSchemaMgr) AddAccountSchema(domain string) (*AccountSchema, error) {
	if _, ok := a.schemas[domain]; ok {
		return nil, errors.New("existed domain schema")
	}
	accountSchema := &AccountSchema{domain: domain,
		groups:         make(map[string]interface{}),
		loginIdSchemas: make(map[string]*LoginIdSchema),
		optionSchemas:  make(map[string]*OptionSchema)}
	a.schemas[domain] = accountSchema
	return accountSchema, nil
}

// AddAccountSchema quick func of add
func AddAccountSchema(domain string) (*AccountSchema, error) {
	if AccountSchemaMgr == nil {
		initAccountSchema()
	}
	return AccountSchemaMgr.AddAccountSchema(domain)
}

// GetAccountSchema get
func (a *accountSchemaMgr) GetAccountSchema(domain string) (*AccountSchema, error) {
	if v, ok := a.schemas[domain]; ok {
		return v, nil
	}
	return nil, errors.New(domain + " not found")
}

// GetAccountSchema quick func of
func GetAccountSchema(domain string) (*AccountSchema, error) {
	if AccountSchemaMgr == nil {
		initAccountSchema()
	}
	return AccountSchemaMgr.GetAccountSchema(domain)
}

// DelAccountSchema del
func (a *accountSchemaMgr) DelAccountSchema(domain string) error {
	if _, ok := a.schemas[domain]; ok {
		delete(a.schemas, domain)
		return nil
	}
	return errors.New("not found")
}

// DelAccountSchema quick func of del
func DelAccountSchema(domain string) error {
	if AccountSchemaMgr == nil {
		initAccountSchema()
	}
	return AccountSchemaMgr.DelAccountSchema(domain)
}
