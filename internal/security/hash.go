package security

import "crypto/sha256"

// SHA256 hashes the given data using SHA-256.
func SHA256(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}
