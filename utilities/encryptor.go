package utilities

import (
	"crypto/md5"
	"fmt"
	"io"
)

type Encryptor interface {
	Encrypt(content string, param interface{}) (encrypted string, err error)
}

type SaultEncryptor struct {
	salt1 string
	salt2 string
}

func NewSaultEncryptor(salt1 string, salt2 string) *SaultEncryptor {
	return &SaultEncryptor{salt1, salt2}
}

// Encrypt md5(salt1 + param + salt2 + md5(content))
func (encryptor *SaultEncryptor) Encrypt(content string, param interface{}) (string, error) {
	h := md5.New()
	io.WriteString(h, content)
	pwMd5 := fmt.Sprintf("%x", h.Sum(nil))

	//salt1+用户名+salt2+MD5拼接
	io.WriteString(h, encryptor.salt1)
	io.WriteString(h, param.(string))
	io.WriteString(h, encryptor.salt2)
	io.WriteString(h, pwMd5)
	last := fmt.Sprintf("%x", h.Sum(nil))
	return last, nil
}
