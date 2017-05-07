package idgen

type IdGenerator interface {
	Generate(param ...interface{}) (string, error)
}
