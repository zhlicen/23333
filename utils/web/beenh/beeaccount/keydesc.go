package beeaccount

import (
	"errors"
	"regexp"
	"strings"
)

type KeyName string

func (k KeyName) String() string {
	return string(k)
}

type KeyDescriptor struct {
	Name                KeyName
	Format, Description string
	CaseSensitive       bool
}

func (desc *KeyDescriptor) Validate(key string) bool {
	if !desc.CaseSensitive {
		key = strings.ToLower(key)
	}
	matched, err := regexp.MatchString(desc.Format, key)
	if err != nil {
		return false
	}
	return matched
}

func NewKeyDescriptor(name KeyName, format string, caseSensortive bool, description string) (*KeyDescriptor, error) {
	descriptor := &KeyDescriptor{name, format, description, caseSensortive}
	err := GlobalKeyDescriptorRegistry.Register(*descriptor)
	return descriptor, err
}

type keyDescriptorRegistry struct {
	idd_map map[KeyName]KeyDescriptor
}

func NewKeyDescriptorRegistry() *keyDescriptorRegistry {
	return &keyDescriptorRegistry{make(map[KeyName]KeyDescriptor)}
}

func (r *keyDescriptorRegistry) Register(descriptor KeyDescriptor) error {
	r.idd_map[descriptor.Name] = descriptor
	return nil
}

func (r *keyDescriptorRegistry) Get(name KeyName) (*KeyDescriptor, error) {
	if descriptor, ok := r.idd_map[name]; ok {
		return &descriptor, nil
	}
	return nil, errors.New("not found")
}

var GlobalKeyDescriptorRegistry *keyDescriptorRegistry

func GetKeyDescriptor(name KeyName) (*KeyDescriptor, error) {
	return GlobalKeyDescriptorRegistry.Get(name)
}

func ValidateKey(name KeyName, value string) (bool, error) {
	descriptor, getErr := GlobalKeyDescriptorRegistry.Get(name)
	if getErr != nil {
		return false, getErr
	}
	return descriptor.Validate(value), nil
}

func initKeyDescRegistry() {
	GlobalKeyDescriptorRegistry = NewKeyDescriptorRegistry()
}
