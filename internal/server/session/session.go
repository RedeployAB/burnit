package session

import (
	"encoding/base64"
	"time"

	"github.com/RedeployAB/burnit/internal/secret"
	"github.com/google/uuid"
)

const (
	// defaultSessionDuration is the default duration for a session.
	defaultSessionDuration = 24 * time.Hour
	// defaultCSFRDuration is the default duration for a CSFR token.
	defaultCSFRDuration = 15 * time.Minute
)

// Session is a struct that holds the session data.
type Session struct {
	id        string
	csfr      CSFR
	expiresAt time.Time
}

// SessionOptions is a struct that holds the session options.
type SessionOptions struct {
	ExpiresAt   time.Time
	CSFR        CSFR
	CSFROptions []CSFROption
}

// SessionOption is a function that sets a session option.
type SessionOption func(s *SessionOptions)

// NewSession creates a new session.
func NewSession(options ...SessionOption) Session {
	opts := SessionOptions{
		ExpiresAt: now().Add(defaultSessionDuration),
	}
	for _, option := range options {
		option(&opts)
	}

	var csfr CSFR
	if opts.CSFROptions != nil {
		csfr = NewCSFR(opts.CSFROptions...)
	} else {
		csfr = opts.CSFR
	}

	return Session{
		id:        newUUID(),
		expiresAt: opts.ExpiresAt,
		csfr:      csfr,
	}
}

// SetCSFR sets the CSFR token.
func (s *Session) SetCSFR(csfr CSFR) Session {
	if csfr != (CSFR{}) {
		s.csfr = csfr
	}
	return *s
}

// ID returns the session ID.
func (s *Session) ID() string {
	return s.id
}

// ExpiresAt returns the expiration time of the session.
func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}

// CSFR returns the CSFR token.
func (s *Session) CSFR() CSFR {
	return s.csfr
}

// DeleteCSFR deletes the CSFR token.
func (s *Session) DeleteCSFR() Session {
	s.csfr = CSFR{}
	return *s
}

// CSFR is a struct that holds the CSFR token and its expiration time.
type CSFR struct {
	token     string
	expiresAt time.Time
}

// CSFROption is a function that sets a CSFR option.
type CSFROption func(c *CSFR)

// NewCSFR creates a new CSFR.
func NewCSFR(options ...CSFROption) CSFR {
	c := CSFR{
		token:     randomString(),
		expiresAt: now().Add(defaultCSFRDuration),
	}
	for _, option := range options {
		option(&c)
	}
	return c
}

// Token returns the CSFR token.
func (c CSFR) Token() string {
	return c.token
}

// ExpiresAt returns the expiration time of the CSFR token.
func (c CSFR) ExpiresAt() time.Time {
	return c.expiresAt
}

// newUUID is a function that returns a new UUID.
var newUUID = func() string {
	return uuid.New().String()
}

var now = func() time.Time {
	return time.Now().UTC()
}

var randomString = func() string {
	return base64.RawURLEncoding.EncodeToString([]byte(secret.Generate(32, true)))
}
