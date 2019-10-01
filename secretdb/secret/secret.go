package secret

import (
	"time"

	"github.com/RedeployAB/redeploy-secrets/common/security"
)

// Secret is to be used as the middle hand between
// incoming requests and the data model.
type Secret struct {
	ID         string
	Secret     string
	Passphrase string
	CreatedAt  time.Time
	ExpiresAt  time.Time
}

// VerifyPassphrase compares a hash with a string,
// if no hash is passed, it always return true.
func VerifyPassphrase(hash, passphrase string) bool {
	if len(hash) > 0 && !security.CompareHash(hash, passphrase) {
		return false
	}
	return true
}

// Hash creates a hash with the help of bcrypt and
// returns it.
func Hash(s string) string {
	return security.Hash(s)
}

// Encrypt the field secret and return it.
func Encrypt(plaintext, passphrase string) string {
	encrypted, err := security.Encrypt([]byte(plaintext), passphrase)
	if err != nil {
		panic(err)
	}
	return string(encrypted)
}

// Decrypt the field secret end return it.
func Decrypt(ciphertext, passphrase string) string {
	decrypted, err := security.Decrypt([]byte(ciphertext), passphrase)
	if err != nil {
		panic(err)
	}
	return string(decrypted)
}
