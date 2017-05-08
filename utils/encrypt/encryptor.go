//copyright

/*
Package encrypt defines the interface for encrypting string and
implemented some encryptors with some popular encrypt methods
*/
package encrypt

// Encryptor defines interface for encrypting string
type Encryptor interface {
	Encrypt(content string, param interface{}) (encrypted string, err error)
}
