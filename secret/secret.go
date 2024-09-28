package secret

import (
	"math/rand"
	"strings"
	"time"
)

const (
	// charset is the standard letters used for generating a secret.
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUWXYZ"
	// specialCharset is the special characters used for generating a secret.
	specialCharset = "_-!?=()&%"
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

// Generate new secret. The length of the secret is set by the provided
// length argument (with a max of 512 characters, a longer length will be trimmed to this value).
// If specialCharacters is set to true, the secret will contain special characters.
func Generate(length int, specialCharacters bool) string {
	if length > maxLength {
		length = maxLength
	}

	chars := charset
	if specialCharacters {
		chars += specialCharset
	}

	var builder strings.Builder
	builder.Grow(length)
	len := len(chars)

	for i := 0; i < length; i++ {
		builder.WriteByte(chars[rand.Intn(len)])
	}

	return builder.String()
}

// generate new secret. The length of the secret is set by the provided
// length argument (with a max of 512 characters, a longer length will be trimmed to this value).
// If specialCharacters is set to true, the secret will contain special characters.
var generate = Generate
