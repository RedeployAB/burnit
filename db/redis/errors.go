package redis

import "errors"

var (
	// ErrKeyNotFound is returned when the key is not found.
	ErrKeyNotFound = errors.New("key not found")
)
