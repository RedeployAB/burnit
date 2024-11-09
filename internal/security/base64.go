package security

import (
	"encoding/base64"
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrInvalidHash is returned when the hash has an invalid length.
	ErrInvalidHashLength = errors.New("invalid hash length")
	// ErrInvalidBase64 is returned when the base64 encoding is invalid.
	ErrInvalidBase64 = errors.New("invalid base64 encoding")
)

// DecodeBase64SHA256 decodes a base64 SHA-256 hash and returns the decoded
// value as a byte slice. Supports both standard and URL (safe) encoding, both
// with and without padding.
func DecodeBase64SHA256(hash string) ([]byte, error) {
	if len(hash) < 43 || len(hash) > 44 {
		return nil, ErrInvalidHashLength
	}

	hash = strings.Replace(hash, "=", "", -1)
	var encoding *base64.Encoding
	re := regexp.MustCompile(`[/+]`)
	if !re.MatchString(hash) {
		encoding = base64.RawURLEncoding
	} else {
		encoding = base64.RawStdEncoding
	}

	dst := make([]byte, encoding.DecodedLen(len(hash)))
	n, err := encoding.Decode(dst, []byte(hash))
	if err != nil {
		switch err.(type) {
		case base64.CorruptInputError:
			return nil, ErrInvalidBase64
		}
		return nil, err
	}
	return dst[:n], nil
}
