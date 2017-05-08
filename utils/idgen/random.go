package idgen

import (
	"crypto/rand"
	r "math/rand"
	"time"
)

// randomIdGenerator generate id with random characters
type randomIdGenerator struct {
	len       int
	alphabets []byte
}

// NewRandomIdGenerator create a randomIdGenerator
// len is the default length of id to be generated
// alphabets is the characters of []bytes used to generate the id
// if alphabets is not specified, defalt alphabets with be used
func NewRandomIdGenerator(len int, alphabets ...byte) *randomIdGenerator {
	if alphabets == nil {
		return &randomIdGenerator{len,
			[]byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")}
	}
	return &randomIdGenerator{len, alphabets}
}

// Generate generate a id
// param is the length of id to be generated
// if param is not specified, defalt length will be used
// return is the id generated
func (g *randomIdGenerator) Generate(param ...interface{}) (string, error) {
	var n int
	if len(param) == 1 {
		n = param[0].(int)
	} else {
		n = g.len
	}
	var bytes = make([]byte, n)
	var randBy bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randBy = true
	}
	for i, b := range bytes {
		if randBy {
			bytes[i] = g.alphabets[r.Intn(len(g.alphabets))]
		} else {
			bytes[i] = g.alphabets[b%byte(len(g.alphabets))]
		}
	}
	return string(bytes), nil
}
