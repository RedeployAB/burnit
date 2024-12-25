package session

import "errors"

var (
	// ErrNilStore is returned when no store is provided.
	ErrNilStore = errors.New("nil store")
	// ErrSessionNotFound is returned when the session is not found.
	ErrSessionNotFound = errors.New("session not found")
	// ErrSessionExpired is returned when the session has expired.
	ErrSessionExpired = errors.New("session expired")
)
