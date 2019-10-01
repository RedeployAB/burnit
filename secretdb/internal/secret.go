package internal

import (
	"github.com/RedeployAB/redeploy-secrets/common/security"
)

// VerifyPassphrase compares a hash with a string,
// if no hash is passed, it always return true.
func VerifyPassphrase(hash, passphrase string) bool {
	if len(hash) > 0 && !security.CompareHash(hash, passphrase) {
		return false
	}
	return true
}
