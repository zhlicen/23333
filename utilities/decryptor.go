package utilities

type Decryptor interface {
	Decrypt(content string) (decrypted string, err error)
}

type AesDecryptor struct {
}
