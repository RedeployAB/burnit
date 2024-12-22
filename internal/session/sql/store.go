package sql

import (
	"time"
)

const (
	// defaultStoreCleanupInterval is the default interval for cleaning up expired sessions
	// and CSRF tokens in the store.
	defaultStoreCleanupInterval = 5 * time.Second //time.Minute
	// defaultStoreTimeout is the default timeout for Store actions.
	defaultStoreTimeout = 5 * time.Second
)
