package account

import (
	"regexp"
)

type KeyDescriptor struct {
	Name, Format, Description string
	CaseSensitive             bool
}

func (desc *KeyDescriptor) Validate(id string) bool {
	matched, err := regexp.MatchString(desc.Format, id)
	if err != nil {
		// log
		return false
	}
	return matched
}

func NewKeyDescriptor(name, format string, caseSensortive bool, description string) (*KeyDescriptor, error) {
	descriptor := &KeyDescriptor{name, format, description, caseSensortive}
	err := GlobalKeyDescriptorRegistry.Register(*descriptor)
	return descriptor, err
}

type keyDescriptorRegistry struct {
	idd_map map[string]KeyDescriptor
}

func NewKeyDescriptorRegistry() *keyDescriptorRegistry {
	return &keyDescriptorRegistry{make(map[string]KeyDescriptor)}
}

func (registry *keyDescriptorRegistry) Register(descriptor KeyDescriptor) error {
	registry.idd_map[descriptor.Name] = descriptor
	return nil
}

func (registry *keyDescriptorRegistry) Get(name string) (*KeyDescriptor, error) {
	descriptor := registry.idd_map[name]
	return &descriptor, nil
}

var GlobalKeyDescriptorRegistry *keyDescriptorRegistry

func GetKeyDescriptor(name string) (*KeyDescriptor, error) {
	return GlobalKeyDescriptorRegistry.Get(name)
}

func ValidateKey(name string, value string) (bool, error) {
	descriptor, getErr := GlobalKeyDescriptorRegistry.Get(name)
	if getErr != nil {
		return false, getErr
	}
	return descriptor.Validate(value), nil
}

func initKeyDescRegistry() {
	GlobalKeyDescriptorRegistry = NewKeyDescriptorRegistry()
}
