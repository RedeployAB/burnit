package session

import "time"

const (
	// defaultCSRFDuration is the default duration for a CSRF token.
	defaultCSRFDuration = 15 * time.Minute
)

// CSRF is a struct that holds the CSRF token and its expiration time.
type CSRF struct {
	token     string
	expiresAt time.Time
}

// CSRFOption is a function that sets a CSRF option.
type CSRFOption func(c *CSRF)

// NewCSRF creates a new CSRF.
func NewCSRF(options ...CSRFOption) CSRF {
	c := CSRF{
		token:     randomString(),
		expiresAt: now().Add(defaultCSRFDuration),
	}
	for _, option := range options {
		option(&c)
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
