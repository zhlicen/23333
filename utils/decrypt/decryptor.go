package decrypt

type Decryptor interface {
	Decrypt(content string) (decrypted string, err error)
}
