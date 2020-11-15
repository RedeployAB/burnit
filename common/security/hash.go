// Package security contains functions to handle
// encryption and encoding.
package security

import (
	"crypto/md5"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// Bcrypt creates a hash with the help of bcrypt.
func Bcrypt(s string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(h)
}

// CompareBcryptHash compares a hash with a string. Returns true if match,
// false otherwise.
func CompareBcryptHash(hash, s string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
	if err != nil {
		return false
	}
	return true
}

// ToMD5 hashes a string to MD5.
func ToMD5(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

// CompareMD5Hash compares a hash with a string. Returns true if match
// false otherwise.
func CompareMD5Hash(hash, s string) bool {
	if ToMD5(s) == hash {
		return true
	}
	return false
}
