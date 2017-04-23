package account

import (
	"errors"
)

type IdDescriptor struct {
	Name, Format, Description string
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
