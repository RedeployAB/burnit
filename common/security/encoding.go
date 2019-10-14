package security

import "encoding/base64"

// ToBase64 encodes a byte slice containing a string
// and returns a base64 encoded string.
func ToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// FromBase64 accepts byte slice containing a base64 value and
// decodes it. Returns empty string if it fails.
func FromBase64(b []byte) []byte {
	data, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		panic(err)
	}
	return data
}
