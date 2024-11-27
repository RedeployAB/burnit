package session

import "time"

// WithExpiresIn sets the expiration time of the session.
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
	return func(c *CSRF) {
		c.expiresAt = exp
	}
}
