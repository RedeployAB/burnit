package security

import (
	"encoding/base64"
	"errors"
	"regexp"
)

var (
	// ErrInvalidBase64 is returned when the base64 encoding is invalid.
	ErrInvalidBase64 = errors.New("invalid base64 encoding")
)

// DecodeBase64 decodes a base64 string.
func DecodeBase64(s string) ([]byte, error) {
	reStd := regexp.MustCompile(`^[a-zA-Z0-9+/]*={0,2}$`)
	reURL := regexp.MustCompile(`^[a-zA-Z0-9-_]*={0,2}$`)

	var encoding *base64.Encoding
	if reStd.MatchString(s) {
		encoding = base64.StdEncoding
	} else if reURL.MatchString(s) {
		encoding = base64.URLEncoding
	} else {
		return nil, ErrInvalidBase64
	}

	b, err := encoding.DecodeString(s)
	if err != nil {
		return nil, ErrInvalidBase64
	}

	return b, nil
}
