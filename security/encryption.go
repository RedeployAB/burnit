package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
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

// toMD5 hashes the given key using MD5.
func ToMD5(key []byte) []byte {
	hasher := md5.New()
	hasher.Write(key)
	return hasher.Sum(nil)
}
