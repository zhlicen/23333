package utilities

import (
	"crypto/rand"
	r "math/rand"
	"time"
)

type KeyGenerator interface {
	Generate(param ...interface{}) (string, error)
}

type randomKeyGenerator struct {
	len       int
	alphabets []byte
}

func NewRandomKeyGenerator(len int, alphabets ...byte) *randomKeyGenerator {
	if alphabets == nil {
		return &randomKeyGenerator{len,
			[]byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")}
	}
	return &randomKeyGenerator{len, alphabets}
}

func (g *randomKeyGenerator) Generate(param ...interface{}) (string, error) {
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
