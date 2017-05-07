package verify

type Verifier struct {
	vs    VerifyService
	key   string
	token string
}

func NewVerifier(vs VerifyService, key string, token string) *Verifier {
	if vs == nil {
		return nil
	}
	return &Verifier{vs, key, token}
}

func (v *Verifier) Verify() error {
	return v.vs.Verify(v.key, v.token)
}

type VerifyService interface {
	SendToken(key string, params ...interface{}) error
	Verify(key string, token string) error
}
