package sql

import "errors"

var (
	// ErrDriverNotSupported is returned when the driver is not supported.
	ErrDriverNotSupported = errors.New("driver not supported")
)
