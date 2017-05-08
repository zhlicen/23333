package beeaccount

import (
	"errors"
	"regexp"
	"strings"
)

// KeyName key name
type KeyName string

func (k KeyName) String() string {
	return string(k)
}

// KeyDescriptor key descriptor
type KeyDescriptor struct {
	Name                KeyName
	Format, Description string
	CaseSensitive       bool
}

// Validate validate
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

// NewKeyDescriptor constructor of KeyDescriptor
func NewKeyDescriptor(name KeyName, format string, caseSensortive bool, description string) (*KeyDescriptor, error) {
	descriptor := &KeyDescriptor{name, format, description, caseSensortive}
	err := GlobalKeyDescriptorRegistry.Register(*descriptor)
	return descriptor, err
}

// keyDescriptorRegistry key descriptor registry
type keyDescriptorRegistry struct {
	idd_map map[KeyName]KeyDescriptor
}

// NewKeyDescriptorRegistry new key descriptor registry
func NewKeyDescriptorRegistry() *keyDescriptorRegistry {
	return &keyDescriptorRegistry{make(map[KeyName]KeyDescriptor)}
}

// Register register
func (r *keyDescriptorRegistry) Register(descriptor KeyDescriptor) error {
	r.idd_map[descriptor.Name] = descriptor
	return nil
}

// Get get
func (r *keyDescriptorRegistry) Get(name KeyName) (*KeyDescriptor, error) {
	if descriptor, ok := r.idd_map[name]; ok {
		return &descriptor, nil
	}
	return nil, errors.New("not found")
}

// GlobalKeyDescriptorRegistry golbal key descriptor registry
var GlobalKeyDescriptorRegistry *keyDescriptorRegistry

// GetKeyDescriptor get key descriptor
func GetKeyDescriptor(name KeyName) (*KeyDescriptor, error) {
	return GlobalKeyDescriptorRegistry.Get(name)
}

// ValidateKey validate key
func ValidateKey(name KeyName, value string) (bool, error) {
	descriptor, getErr := GlobalKeyDescriptorRegistry.Get(name)
	if getErr != nil {
		return false, getErr
	}
	return descriptor.Validate(value), nil
}

// initKeyDescRegistry
func initKeyDescRegistry() {
	GlobalKeyDescriptorRegistry = NewKeyDescriptorRegistry()
}
