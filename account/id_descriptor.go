package account

import (
	"errors"
	"fmt"
	"regexp"
)

type IdDescriptor struct {
	KeyDescriptor
}

func NewIdDescriptor(name, format string, caseSensortive bool, description string) (*IdDescriptor, error) {
	descriptor := &IdDescriptor{KeyDescriptor{name, format, description, caseSensortive}}
	err := GlobalIdDescriptorRegistry.Register(*descriptor)
	return descriptor, err
}

type idDescriptorRegistry struct {
	idd_map map[string]IdDescriptor
}

func newIdDescriptorRegistry() *idDescriptorRegistry {
	return &idDescriptorRegistry{make(map[string]IdDescriptor)}
}

func (registry *idDescriptorRegistry) Register(descriptor IdDescriptor) error {
	registry.idd_map[descriptor.Name] = descriptor
	return nil
}

func (registry *idDescriptorRegistry) Match(id string) (*IdDescriptor, error) {
	for _, v := range registry.idd_map {
		matched, err := regexp.MatchString(v.Format, id)
		if err == nil && matched {
			return &v, nil
		}
	}
	fmt.Println("no match descriptor found")
	return nil, errors.New("no match descriptor found")
}

func (registry *idDescriptorRegistry) Get(name string) (*IdDescriptor, error) {
	descriptor := registry.idd_map[name]
	return &descriptor, nil
}

var GlobalIdDescriptorRegistry *idDescriptorRegistry

func initIdRegistry() {
	GlobalIdDescriptorRegistry = newIdDescriptorRegistry()
}

func ValidateId(name string, value string) (bool, error) {
	descriptor, getErr := GlobalIdDescriptorRegistry.Get(name)
	if getErr != nil {
		return false, getErr
	}
	return descriptor.Validate(value), nil
}

func GetIdDescriptor(name string) (*IdDescriptor, error) {
	return GlobalIdDescriptorRegistry.Get(name)
}

func MatchIdDescriptor(id string) (*IdDescriptor, error) {
	return GlobalIdDescriptorRegistry.Match(id)
}
