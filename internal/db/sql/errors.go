package sql

import "errors"

var (
	// ErrNilDB is returned when the database is nil.
	ErrNilDB = errors.New("database is nil")
	// ErrDriverNotSupported is returned when the driver is not supported.
	ErrDriverNotSupported = errors.New("driver not supported")
)
