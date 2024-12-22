package session

import "time"

// WithExpiresAt sets the expiration time of the session.
func WithExpiresAt(exp time.Time) SessionOption {
	return func(o *SessionOptions) {
		o.ExpiresAt = exp
	}
}

// WithCSRF sets the CSRF token.
func WithCSRF(csrf CSRF) SessionOption {
	return func(o *SessionOptions) {
		if csrf != (CSRF{}) {
			o.CSRF = csrf
		}
	}
}

// WithCSRFOptions sets the options for the CSRF token for the session.
func WithCSRFOptions(options ...CSRFOption) SessionOption {
	return func(o *SessionOptions) {
		o.CSRFOptions = options
	}
}

// WithCSRFExpiresAt sets the expiration time of the CSRF token.
func WithCSRFExpiresAt(exp time.Time) CSRFOption {
	return func(o *CSRFOptions) {
		o.ExpiresAt = exp
	}
}

// Get a session by its ID.
func WithID(id string) GetOption {
	return func(o *GetOptions) {
		o.ID = id
	}
}

// Get a session by its CSRF token.
func WithCSRFToken(token string) GetOption {
	return func(o *GetOptions) {
		o.CSRFToken = token
	}
}
