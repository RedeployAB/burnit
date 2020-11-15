// Package security contains functions to handle
// encryption and encoding.
// Source: https://github.com/gtank/cryptopasta
// Modified to encode/decode from base64 and a key string
// and encode it to MD5. Some minor modifications in byte/string
// handling ass well. Encrypts output a base64 encoded string,
// Decrypt takes a base64 encoded string.
package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Encrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Encrypt(plaintext []byte, key string) (ciphertext []byte, err error) {
	md5key := ToMD5(key)
	block, err := aes.NewCipher([]byte(md5key))
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

	return []byte(ToBase64(gcm.Seal(nonce, nonce, plaintext, nil))), nil
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Decrypt(ciphertext []byte, key string) (plaintext []byte, err error) {
	md5Key := ToMD5(key)
	block, err := aes.NewCipher([]byte(md5Key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	decodedCipher := FromBase64(ciphertext)
	if len(decodedCipher) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		decodedCipher[:gcm.NonceSize()],
		decodedCipher[gcm.NonceSize():],
		nil,
	)
}
