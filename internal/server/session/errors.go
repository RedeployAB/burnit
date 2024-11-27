package session

import "errors"

var (
	// ErrSessionNotFound is returned when the session is not found.
	ErrSessionNotFound = errors.New("session not found")
)
