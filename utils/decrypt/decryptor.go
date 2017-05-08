//copyright

/*
Package decrypt defines the interface for decrypting string and
implemented some decryptors with some popular decrypt methods
*/

package decrypt

// Decryptor defines interface for decrypting string
type Decryptor interface {
	Decrypt(content string) (decrypted string, err error)
}
