package beaccount

import (
	"errors"
	"regexp"
	"strings"
)

type IdName string

func (i IdName) String() string {
	return string(i)
}

type IdDescriptor struct {
	Name                IdName
	Format, Description string
	CaseSensitive       bool
}

func (desc *IdDescriptor) Validate(id string) bool {
	if !desc.CaseSensitive {
		id = strings.ToLower(id)
	}
	matched, err := regexp.MatchString(desc.Format, id)
	if err != nil {
		// log
		return false
	}
	return matched
}

func NewIdDescriptor(name IdName, format string, caseSensortive bool, description string) (*IdDescriptor, error) {
	descriptor := &IdDescriptor{name, format, description, caseSensortive}
	err := GlobalIdDescriptorRegistry.Register(*descriptor)
	return descriptor, err
}

type idDescriptorRegistry struct {
	idd_map map[IdName]IdDescriptor
}

func newIdDescriptorRegistry() *idDescriptorRegistry {
	return &idDescriptorRegistry{make(map[IdName]IdDescriptor)}
}

func (registry *idDescriptorRegistry) Register(descriptor IdDescriptor) error {
	registry.idd_map[IdName(descriptor.Name)] = descriptor
	return nil
}

func (registry *idDescriptorRegistry) Match(id string) (*IdDescriptor, error) {
	for _, v := range registry.idd_map {
		matched := v.Validate(id)
		if matched {
			return &v, nil
		}
	}
	return nil, errors.New("no match descriptor found")
}

func (registry *idDescriptorRegistry) Get(name IdName) (*IdDescriptor, error) {
	descriptor := registry.idd_map[name]
	return &descriptor, nil
}

var GlobalIdDescriptorRegistry *idDescriptorRegistry

func initIdRegistry() {
	GlobalIdDescriptorRegistry = newIdDescriptorRegistry()
}

func ValidateId(name IdName, value string) (bool, error) {
	descriptor, getErr := GlobalIdDescriptorRegistry.Get(name)
	if getErr != nil {
		return false, getErr
	}
	return descriptor.Validate(value), nil
}

func GetIdDescriptor(name IdName) (*IdDescriptor, error) {
	return GlobalIdDescriptorRegistry.Get(name)
}

func MatchIdDescriptor(id string) (*IdDescriptor, error) {
	return GlobalIdDescriptorRegistry.Match(id)
}
