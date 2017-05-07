package encrypt

type Encryptor interface {
	Encrypt(content string, param interface{}) (encrypted string, err error)
}
