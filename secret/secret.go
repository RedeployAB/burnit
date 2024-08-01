package secret

import "time"

// Secret contains the secret data.
type Secret struct {
	ID         string
	Value      string
	Passphrase string
	TTL        time.Duration
}
