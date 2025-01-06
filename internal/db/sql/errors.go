package sql

import "errors"

var (
	// ErrNilClient is returned when the client is nil.
	ErrNilClient = errors.New("client is nil")
	// ErrDriverNotSupported is returned when the driver is not supported.
	ErrDriverNotSupported = errors.New("driver not supported")
)
