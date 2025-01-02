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
	// defaultGenerateSecretLength is the default length of a generated secret.
	defaultGenerateSecretLength = 16
	// maxGenerateSecretLength is the maximum length of a generated secret.
	maxGenerateSecretLength = 512
)

// Secret contains the secret data.
type Secret struct {
	ID         string
	Value      string
	Passphrase string
	TTL        time.Duration
	ExpiresAt  time.Time
}

// GenerateOptions contains the options for generating a new secret.
type GenerateOptions struct {
	Length            int
	SpecialCharacters bool
}

// GenerateOption is a function that sets options for generating a new secret.
type GenerateOption func(o *GenerateOptions)

// Generate new secret. The length of the secret is set by the provided
// length argument (with a max of 512 characters, a longer length will be trimmed to this value).
// If specialCharacters is set to true, the secret will contain special characters.
func Generate(options ...GenerateOption) string {
	opts := GenerateOptions{
		Length: defaultGenerateSecretLength,
	}
	for _, option := range options {
		option(&opts)
	}
	if opts.Length > maxGenerateSecretLength {
		opts.Length = maxGenerateSecretLength
	}

	chars := charset
	if opts.SpecialCharacters {
		chars += specialCharset
	}

	var builder strings.Builder
	builder.Grow(opts.Length * 3)
	len := len(chars)

	for i := 0; i < opts.Length; i++ {
		builder.WriteByte(chars[rand.Intn(len)])
	}

	return builder.String()
}

// generate new secret. The length of the secret is set by the provided
// length argument (with a max of 512 characters, a longer length will be trimmed to this value).
// If specialCharacters is set to true, the secret will contain special characters.
var generate = Generate
