package secret

import "time"

// WithTimeout sets the timeout for the service.
func WithTimeout(d time.Duration) ServiceOption {
	return func(s *service) {
		s.timeout = d
	}
}
