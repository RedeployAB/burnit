package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
)

var (
	// ErrInvalidKey is returned when the provided key is invalid.
	ErrInvalidKey = errors.New("invalid key")
)

// Encrypt data with 256-bit AES-GCM encryption using the given key.
func Encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Decrypt data encrypted with 256-bit AES-GCM encryption using the given key.
func Decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, errors.New("malformed data")
	}

	decrypted, err := gcm.Open(
		nil,
		data[:gcm.NonceSize()],
		data[gcm.NonceSize():],
		nil,
	)
	if err != nil {
		return nil, ErrInvalidKey
	}

	return decrypted, nil
}

// ToSHA256 hashes the given data using SHA-256.
func ToSHA256(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

// SHA256ToHex encodes a SHA-256 hash to a hex string.
func SHA256ToHex(hash []byte) string {
	return hex.EncodeToString(hash)
}

// DecodeSHA256HexString decodes a SHA-256 hash from a hex string.
func DecodeSHA256HexString(hash string) ([]byte, error) {
	if len(hash) != 64 {
		return nil, ErrInvalidKey
	}
	return hex.DecodeString(hash)
}
