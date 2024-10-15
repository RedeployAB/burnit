package redis

import "errors"

var (
	// ErrNilClient is returned when the client is nil.
	ErrNilClient = errors.New("client is nil")
	// ErrKeyNotFound is returned when the key is not found.
	ErrKeyNotFound = errors.New("key not found")
)
