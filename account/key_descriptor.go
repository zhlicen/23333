package account

import (
	"errors"
	"regexp"
)

type KeyName string

type KeyDescriptor struct {
	Name                KeyName
	Format, Description string
	CaseSensitive       bool
}

func (desc *KeyDescriptor) Validate(id string) bool {
	matched, err := regexp.MatchString(desc.Format, id)
	if err != nil {
		// log
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
