package security

import (
	"crypto/md5"
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
