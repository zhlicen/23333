package beeaccount

import (
	"errors"
	"regexp"
	"strings"
)

// IdName type of login id name
type IdName string

func (i IdName) String() string {
	return string(i)
}

// IdDescriptor descriptor of login id
// member:Name login id type name, eg: UserName, Mobile
// member:Format regex patten of the id string
// member:Description description information of this id for debug
// member:CaseSensitive if this id is case sensitive, if not,
// the id value will be convert to lower case automatically
type IdDescriptor struct {
	Name                IdName
	Format, Description string
	CaseSensitive       bool
}

// Validate check if the id provided is valid
// id is the content to be validate
// return if the id is valid
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

// NewIdDescriptor constructor of IdDescriptor
// this function will register this new descriptor to GlobalIdDescriptorRegistry
// return the new descriptor and error if there's any
func NewIdDescriptor(name IdName, format string, caseSensortive bool, description string) (*IdDescriptor, error) {
	descriptor := &IdDescriptor{name, format, description, caseSensortive}
	err := GlobalIdDescriptorRegistry.Register(*descriptor)
	return descriptor, err
}

// idDescriptorRegistry container for storing IdDescriptor(s)
type idDescriptorRegistry struct {
	idd_map map[IdName]IdDescriptor
}

// newIdDescriptorRegistry constructor of idDescriptorRegistry
func newIdDescriptorRegistry() *idDescriptorRegistry {
	return &idDescriptorRegistry{make(map[IdName]IdDescriptor)}
}

// Register register new descriptor to registry
func (registry *idDescriptorRegistry) Register(descriptor IdDescriptor) error {
	registry.idd_map[IdName(descriptor.Name)] = descriptor
	return nil
}

// Match match the descriptor with login id string
// if no match found, an error will be returned
func (registry *idDescriptorRegistry) Match(id string) (*IdDescriptor, error) {
	for _, v := range registry.idd_map {
		matched := v.Validate(id)
		if matched {
			return &v, nil
		}
	}
	return nil, errors.New("no match descriptor found")
}

// Get get a descriptor with id name
func (registry *idDescriptorRegistry) Get(name IdName) (*IdDescriptor, error) {
	descriptor := registry.idd_map[name]
	return &descriptor, nil
}

// GlobalIdDescriptorRegistry global descriptor registry
var GlobalIdDescriptorRegistry *idDescriptorRegistry

func initIdRegistry() {
	GlobalIdDescriptorRegistry = newIdDescriptorRegistry()
}

// ValidateId quick method for validate id
func ValidateId(name IdName, value string) (bool, error) {
	descriptor, getErr := GlobalIdDescriptorRegistry.Get(name)
	if getErr != nil {
		return false, getErr
	}
	return descriptor.Validate(value), nil
}

// GetIdDescriptor quick method for getting an IdDescriptor
func GetIdDescriptor(name IdName) (*IdDescriptor, error) {
	return GlobalIdDescriptorRegistry.Get(name)
}

// MatchIdDescriptor quick method for matching an IdDescriptor
func MatchIdDescriptor(id string) (*IdDescriptor, error) {
	return GlobalIdDescriptorRegistry.Match(id)
}
