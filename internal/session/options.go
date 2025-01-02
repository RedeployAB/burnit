package session

import "time"

// WithTimeout sets the timeout for the service.
func WithTimeout(d time.Duration) ServiceOption {
	return func(s *service) {
		s.timeout = d
	}
}

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

// GetWithID a session by its ID.
func GetWithID(id string) GetOption {
	return func(o *GetOptions) {
		o.ID = id
	}
}

// GetWithCSRFToken a session by its CSRF token.
func GetWithCSRFToken(token string) GetOption {
	return func(o *GetOptions) {
		o.CSRFToken = token
	}
}

// DeleteWithID deletes a session by its ID.
func DeleteWithID(id string) DeleteOption {
	return func(o *DeleteOptions) {
		o.ID = id
	}
}

// DeleteWithCSRFToken deletes a session by its CSRF token.
func DeleteWithCSRFToken(token string) DeleteOption {
	return func(o *DeleteOptions) {
		o.CSRFToken = token
	}
}
