package secrets

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/RedeployAB/burnit/burnit/db"
	"github.com/RedeployAB/burnit/burnit/security"
)

const (
	stdLetters         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUWXYZ"
	generatedMaxLength = 512
)

// Secret is to be used as the middle-hand between
// incoming JSON payload and the data model.
type Secret struct {
	ID         string
	Value      string
	Passphrase string
	CreatedAt  time.Time
	ExpiresAt  time.Time
	TTL        int
}

// Generate generates a random string. Argument
// l determines amount of characters in the
// resulting string. Argument sc determines if
// special characters should be used.
func Generate(l int, sc bool) string {
	if l > generatedMaxLength {
		l = generatedMaxLength
	}

	var strb strings.Builder
	strb.WriteString(stdLetters)
	if sc {
		strb.WriteString("_-!?=()&%")
	}
	bltrs := []byte(strb.String())

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, l)
	for i := range b {
		b[i] = bltrs[rand.Intn(len(bltrs))]
	}

	return string(b)
}

// NewFromJSON creates an incoming JSON payload
// and creates a Secret from it.
func NewFromJSON(b io.ReadCloser) (*Secret, error) {
	var s *Secret
	if err := json.NewDecoder(b).Decode(&s); err != nil {
		return nil, err
	}

	if len(s.Value) == 0 {
		return nil, errors.New("a provided secret (value) is missing")
	}

	if len(s.Value) > maxLength {
		return nil, errors.New("provided value is too large, character limit exceeds: " + strconv.Itoa(maxLength))
	}

	if s.TTL != 0 {
		s.ExpiresAt = time.Now().Add(time.Minute * time.Duration(s.TTL))
	} else {
		s.ExpiresAt = time.Now().AddDate(0, 0, 7)
	}

	return s, nil
}

// toModel transforms a Secret to the
// data model variant of Secret.
func toModel(s *Secret, passphrase string) *db.Secret {
	var createdAt, expiresAt time.Time
	if s.CreatedAt.IsZero() {
		createdAt = time.Now()
	} else {
		createdAt = s.CreatedAt
	}

	if s.ExpiresAt.IsZero() {
		expiresAt = time.Now().Add(time.Minute * time.Duration(10080))
	} else {
		expiresAt = s.ExpiresAt
	}

	val, err := security.Encrypt([]byte(s.Value), passphrase)
	if err != nil {
		return nil
	}

	return &db.Secret{
		Value:     string(val),
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}
}

// toSecret transforms the data model variant
// of Secret to a Secret.
func toSecret(s *db.Secret, passphrase string) *Secret {
	var val string
	if len(s.Value) > 0 {
		decrypted, err := security.Decrypt([]byte(s.Value), passphrase)
		if err != nil {
			return &Secret{}
		}
		val = string(decrypted)
	}

	return &Secret{
		ID:        s.ID,
		Value:     val,
		CreatedAt: s.CreatedAt,
		ExpiresAt: s.ExpiresAt,
	}
}
