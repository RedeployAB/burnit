package session

import (
	"time"

	"github.com/google/uuid"
)

const (
	// defaultSessionDuration is the default duration for a session.
	defaultSessionDuration = 24 * time.Hour
)

// Session is a struct that holds the session data.
type Session struct {
	id        string
	csrf      CSRF
	expiresAt time.Time
}

// SessionOptions is a struct that holds the session options.
type SessionOptions struct {
	ID          string
	ExpiresAt   time.Time
	CSRF        CSRF
	CSRFOptions []CSRFOption
}

// SessionOption is a function that sets a session option.
type SessionOption func(o *SessionOptions)

// NewSession creates a new session.
func NewSession(options ...SessionOption) Session {
	opts := SessionOptions{
		ExpiresAt: now().Add(defaultSessionDuration),
	}
	for _, option := range options {
		option(&opts)
	}

	var csrf CSRF
	if opts.CSRFOptions != nil {
		csrf = NewCSRF(opts.CSRFOptions...)
	} else {
		csrf = opts.CSRF
	}

	var id string
	if len(opts.ID) > 0 {
		id = opts.ID
	} else {
		id = newUUID()
	}

	return Session{
		id:        id,
		expiresAt: opts.ExpiresAt,
		csrf:      csrf,
	}
}

// IsEmpty returns true if the session is empty.
func (s *Session) IsEmpty() bool {
	return s.id == "" && s.expiresAt.IsZero()
}

// SetCSRF sets the CSRF token.
func (s *Session) SetCSRF(csrf CSRF) Session {
	s.csrf = csrf
	return *s
}

// ID returns the session ID.
func (s *Session) ID() string {
	return s.id
}

// ExpiresAt returns the expiration time of the session.
func (s Session) ExpiresAt() time.Time {
	return s.expiresAt
}

// Expired returns true if the session has expired.
func (s Session) Expired() bool {
	return s.expiresAt.Before(now())
}

// CSRF returns the CSRF token.
func (s Session) CSRF() CSRF {
	return s.csrf
}

// DeleteCSRF deletes the CSRF token.
func (s *Session) DeleteCSRF() Session {
	s.csrf = CSRF{}
	return *s
}

// newUUID is a function that returns a new UUID.
var newUUID = func() string {
	return uuid.New().String()
}

var now = func() time.Time {
	return time.Now().UTC()
}
