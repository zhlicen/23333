package encrypt

import (
	"crypto/md5"
	"fmt"
	"io"
)

type saultEncryptor struct {
	salt1 string
	salt2 string
}

func NewSaultEncryptor(salt1 string, salt2 string) *saultEncryptor {
	return &saultEncryptor{salt1, salt2}
}

// Encrypt md5(salt1 + param + salt2 + md5(content))
func (encryptor *saultEncryptor) Encrypt(content string, param interface{}) (string, error) {
	h := md5.New()
	io.WriteString(h, content)
	pwMd5 := fmt.Sprintf("%x", h.Sum(nil))
	io.WriteString(h, encryptor.salt1)
	io.WriteString(h, param.(string))
	io.WriteString(h, encryptor.salt2)
	io.WriteString(h, pwMd5)
	last := fmt.Sprintf("%x", h.Sum(nil))
	return last, nil
}
