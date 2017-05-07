package idgen

import (
	"crypto/rand"
	r "math/rand"
	"time"
)

type randomIdGenerator struct {
	len       int
	alphabets []byte
}

func NewRandomIdGenerator(len int, alphabets ...byte) *randomIdGenerator {
	if alphabets == nil {
		return &randomIdGenerator{len,
			[]byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")}
	}
	return &randomIdGenerator{len, alphabets}
}

func (g *randomIdGenerator) Generate(param ...interface{}) (string, error) {
	n := g.len
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
