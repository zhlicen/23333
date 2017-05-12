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

// loginIDSchema login id schema
// member:Name name of this loginID
// member:IsRequired if the id is required
// member:NeedVerfied if the id needs be verified
// member:Validator validator for id string
type LoginIDSchema struct {
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
	loginIDSchemas map[string]*LoginIDSchema
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

// AddLoginIDSchema add login id schema
func (a *AccountSchema) AddLoginIDSchema(name string, isRequired bool, needVerified bool,
	ignoreCase bool, patten string, userData ...interface{}) error {
	if _, ok := a.loginIDSchemas[name]; ok {
		return errors.New("schema " + name + " exist")
	}
	newSchema := new(LoginIDSchema)
	newSchema.Name = name
	newSchema.IsRequired = isRequired
	newSchema.NeedVerified = needVerified
	newSchema.Validator = NewStringValidator(ignoreCase, patten)
	if len(userData) > 0 {
		newSchema.UserData = userData[0]
	}
	a.loginIDSchemas[name] = newSchema
	return nil
}

func (a *AccountSchema) getRequiredLogIDs() []string {
	var loginIDKeys []string
	for _, v := range a.loginIDSchemas {
		if v.IsRequired {
			loginIDKeys = append(loginIDKeys, v.Name)
		}
	}
	return loginIDKeys
}

// GetLoginIDSchema get login id schema
func (a *AccountSchema) GetLoginIDSchema(name string) (*LoginIDSchema, error) {
	if v, ok := a.loginIDSchemas[name]; ok {
		return v, nil
	}
	return nil, errors.New("schema not found")
}

// MatchLoginID get login id name
func (a *AccountSchema) MatchLoginID(id string) (string, error) {
	for k, v := range a.loginIDSchemas {
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

// AccountSchemaMgr account schema manager
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
		loginIDSchemas: make(map[string]*LoginIDSchema),
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
