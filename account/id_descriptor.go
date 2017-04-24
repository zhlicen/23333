package account

import (
	"errors"
	"fmt"
	"regexp"
)

type IdDescriptor struct {
	Name, Format, Description string
}

func (desc *IdDescriptor) Validate(id string) bool {
	matched, err := regexp.MatchString(desc.Format, id)
	if err != nil {
		// log
		return false
	}
	return matched
}

func NewIdDescriptor(name, format, description string) (*IdDescriptor, error) {
	descriptor := &IdDescriptor{name, format, description}
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
