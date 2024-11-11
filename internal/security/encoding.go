package security

import (
	"encoding/base64"
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrInvalidBase64 is returned when the base64 encoding is invalid.
	ErrInvalidBase64 = errors.New("invalid base64 encoding")
)

// DecodeBase64 decodes a base64 string.
func DecodeBase64(s string) (string, error) {
	s = strings.Replace(s, "=", "", -1)

	var encoding *base64.Encoding
	re := regexp.MustCompile(`[/+]`)
	if !re.MatchString(s) {
		encoding = base64.RawURLEncoding
	} else {
		encoding = base64.RawStdEncoding
	}

	b, err := encoding.DecodeString(s)
	if err != nil {
		return "", ErrInvalidBase64
	}
	return string(b), nil
}
