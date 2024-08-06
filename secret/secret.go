package secret

import (
	"math/rand"
	"strings"
	"time"
)

const (
	// stdCharacters is the standard letters used for generating a secret.
	stdCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUWXYZ"
	// spcCharacters is the special characters used for generating a secret.
	spcCharacters = "_-!?=()&%"
	// maxLength is the maximum length of a secret.
	maxLength = 512
)

// Secret contains the secret data.
type Secret struct {
	ID         string
	Value      string
	Passphrase string
	TTL        time.Duration
	ExpiresAt  time.Time
}

// Generate a new secret. The length of the secret is set by the provided
// length argument (with a max of 512 characters, a longer length will be trimmed to this value).
// If specialCharacters is set to true, the secret will contain special characters.
func Generate(length int, specialCharacters bool) string {
	if length > maxLength {
		length = maxLength
	}

	var strb strings.Builder
	strb.WriteString(stdCharacters)
	if specialCharacters {
		strb.WriteString(spcCharacters)
	}
	bltrs := []byte(strb.String())
	b := make([]byte, length)
	for i := range b {
		b[i] = bltrs[rand.Intn(len(bltrs))]
	}

	return string(b)
}
