package security

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// Hash creates a hash with the help of bcrypt.
func Hash(str string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(h)
}

// CompareHash compares a hash with a string. Returns true if match,
// false otherwise.
func CompareHash(hash string, str string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	if err != nil {
		return false
	}
	return true
}

// ToMD5 hashes a string to MD5.
func ToMD5(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

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
