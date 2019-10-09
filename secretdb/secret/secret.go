package secret

import (
	"encoding/json"
	"io"
	"time"

	"github.com/RedeployAB/burnit/common/security"
)

// Secret is to be used as the middle hand between
// incoming requests and the data model.
type Secret struct {
	ID         string
	Secret     string
	Passphrase string
	CreatedAt  time.Time
	ExpiresAt  time.Time
	TTL        int
}

// NewSecret takes an io.Reader and attempt to decode it to
// a Secret.
func NewSecret(b io.ReadCloser) (*Secret, error) {
	var s *Secret
	if err := json.NewDecoder(b).Decode(&s); err != nil {
		return &Secret{}, err
	}

	if len(s.Passphrase) > 0 {
		s.HashPassphrase()
	}

	exp := time.Now().AddDate(0, 0, 7)
	if s.TTL != 0 {
		exp = time.Now().Add(time.Minute * time.Duration(s.TTL))
	}
	s.ExpiresAt = exp

	return s, nil
}

// VerifyPassphrase compares a hash with a string,
// if no hash is passed, it always return true.
func (s *Secret) VerifyPassphrase(hash string) bool {
	if len(s.Passphrase) > 0 && !security.CompareHash(s.Passphrase, hash) {
		return false
	}
	return true
}

// HashPassphrase hashes the field passphrase.
func (s *Secret) HashPassphrase() *Secret {
	s.Passphrase = security.Hash(s.Passphrase)
	return s
}
