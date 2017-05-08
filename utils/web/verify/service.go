package verify

// Verifier verifier
type Verifier struct {
	vs    VerifyService
	key   string
	token string
}

// NewVerifier create a verifier
func NewVerifier(vs VerifyService, key string, token string) *Verifier {
	if vs == nil {
		return nil
	}
	return &Verifier{vs, key, token}
}

// Verify verify the condition
func (v *Verifier) Verify() error {
	return v.vs.Verify(v.key, v.token)
}

// VerifyService service interface for verification
type VerifyService interface {
	SendToken(key string, params ...interface{}) error
	Verify(key string, token string) error
}
