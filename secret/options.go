package secret

import "time"

// WithEncryptionKey sets the encryption key for the service.
func WithEncryptionKey(key string) ServiceOption {
	return func(s *service) {
		s.encryptionKey = key
	}
}

// WithTimeout sets the timeout for the service.
func WithTimeout(d time.Duration) ServiceOption {
	return func(s *service) {
		s.timeout = d
	}
}
