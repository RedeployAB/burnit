package dto

import (
	"encoding/json"
	"io"
	"time"
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

	exp := time.Now().AddDate(0, 0, 7)
	if s.TTL != 0 {
		exp = time.Now().Add(time.Minute * time.Duration(s.TTL))
	}
	s.ExpiresAt = exp

	return s, nil
}
