package idgen

// randomIdGenerator generate id with random characters
type uuidGenerator struct {
}

// NewUuidGenerator constructor of uuidGenerator
func NewUuidGenerator() *uuidGenerator {
	return &uuidGenerator{}
}

// Generate generate an uuid
// return is the uuid generated with format:
//
func (g *uuidGenerator) Generate(param ...interface{}) (string, error) {
	return "", nil
}
