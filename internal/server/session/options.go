package session

import "time"

// WithExpiresIn sets the expiration time of the session.
func WithExpiresAt(exp time.Time) SessionOption {
	return func(o *SessionOptions) {
		o.ExpiresAt = exp
	}
}

// WithCSFR sets the CSFR token.
func WithCSFR(csfr CSFR) SessionOption {
	return func(o *SessionOptions) {
		if csfr != (CSFR{}) {
			o.CSFR = csfr
		}
	}
}

// WithCSFROptions sets the options for the CSFR token for the session.
func WithCSFROptions(options ...CSFROption) SessionOption {
	return func(o *SessionOptions) {
		o.CSFROptions = options
	}
}

// WithCSFRExpiresAt sets the expiration time of the CSFR token.
func WithCSFRExpiresAt(exp time.Time) CSFROption {
	return func(c *CSFR) {
		c.expiresAt = exp
	}
}
