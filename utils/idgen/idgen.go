//copyright

/*
Package idgen defines the interface for generating a id string and
implemented some generators
*/
package idgen

// IdGenerator interface for generate id
type IdGenerator interface {
	Generate(param ...interface{}) (string, error)
}
