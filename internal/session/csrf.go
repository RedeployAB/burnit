package session

import (
	"encoding/base64"
	"time"

	"github.com/RedeployAB/burnit/internal/secret"
)

const (
	// defaultCSRFDuration is the default duration for a CSRF token.
	defaultCSRFDuration = 15 * time.Minute
)

// CSRF is a struct that holds the CSRF token and its expiration time.
type CSRF struct {
	token     string
	expiresAt time.Time
}

// CSRFOptions is a struct that holds the CSRF options.
type CSRFOptions struct {
	Token     string
	ExpiresAt time.Time
}

// CSRFOption is a function that sets a CSRF option.
type CSRFOption func(o *CSRFOptions)

// NewCSRF creates a new CSRF.
func NewCSRF(options ...CSRFOption) CSRF {
	opts := CSRFOptions{}
	for _, option := range options {
		option(&opts)
	}
	if len(opts.Token) == 0 {
		opts.Token = randomString()
	}
	if opts.ExpiresAt.IsZero() {
		opts.ExpiresAt = now().Add(defaultCSRFDuration)
	}

	c := CSRF{
		token:     opts.Token,
		expiresAt: opts.ExpiresAt,
	}
	return c
}

// IsEmpty returns true if the CSRF token is empty.
func (c CSRF) IsEmpty() bool {
	return len(c.token) == 0 && c.expiresAt.IsZero()
}

// Token returns the CSRF token.
func (c CSRF) Token() string {
	return c.token
}

// ExpiresAt returns the expiration time of the CSRF token.
func (c CSRF) ExpiresAt() time.Time {
	return c.expiresAt
}

// Expired returns true if the CSRF token has expired.
func (c CSRF) Expired() bool {
	return c.expiresAt.Before(now())
}

var randomString = func() string {
	return base64.RawURLEncoding.EncodeToString([]byte(secret.Generate(func(o *secret.GenerateOptions) {
		o.Length = 32
		o.SpecialCharacters = true
	})))
}
