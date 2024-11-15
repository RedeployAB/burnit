package secret

import "time"

// WithTimeout sets the timeout for the service.
func WithTimeout(d time.Duration) ServiceOption {
	return func(s *service) {
		s.timeout = d
	}
}

// WithValueMaxCharacters sets the maximum number of characters
// a secret value can have.
func WithValueMaxCharacters(max int) ServiceOption {
	return func(s *service) {
		s.valueMaxCharacters = max
	}
}
