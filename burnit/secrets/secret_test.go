package secrets

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/RedeployAB/burnit/burnit/db"
	"github.com/RedeployAB/burnit/burnit/security"
)

func TestGenerate(t *testing.T) {
	var tests = []struct {
		x int
		y bool
		n int
	}{
		{8, false, 8},
		{16, true, 16},
	}

	for _, test := range tests {
		secret := Generate(test.x, test.y)
		count := utf8.RuneCountInString(secret)
		if count != test.n {
			t.Errorf("number of characters in generated string incorrect, got: %d, want: %d", count, test.n)
		}
	}
}

func TestNewFromJSON(t *testing.T) {
	str1 := []byte(`{"value":"value1"}`)
	str2 := []byte(`{"value":"value2","passphrase":"1234"}`)
	str3 := []byte(`{"value":"value3","ttl":4320}`)
	strMalformed1 := []byte(`{"value":`)

	expectedDay1 := time.Now().AddDate(0, 0, 7)
	expectedDay2 := time.Now().AddDate(0, 0, 3)

	var tests = []struct {
		input   []byte
		want    *Secret
		wantErr error
	}{
		{input: str1, want: &Secret{Value: "value1", Passphrase: "", ExpiresAt: expectedDay1}, wantErr: nil},
		{input: str2, want: &Secret{Value: "value2", Passphrase: "1234", ExpiresAt: expectedDay1}, wantErr: nil},
		{input: str3, want: &Secret{Value: "value3", Passphrase: "", ExpiresAt: expectedDay2}, wantErr: nil},
		{input: strMalformed1, want: nil, wantErr: errors.New("unexpected EOF")},
	}

	for _, test := range tests {
		got, err := NewFromJSON(io.NopCloser(bytes.NewBuffer(test.input)))
		if got != nil && got.Value != test.want.Value {
			t.Errorf("incorrect value, got: %s, want: %s", got.Value, test.want.Value)
		}
		if got != nil && got.Passphrase != test.want.Passphrase {
			t.Errorf("incorrect value, got: %s, want: %s", got.Passphrase, test.want.Passphrase)
		}
		if err != nil && err.Error() != test.wantErr.Error() {
			t.Errorf("incorrect value, got: %v, want: %v", err.Error(), test.wantErr.Error())
		}
	}
}

func TestNewFromJSONMalformed(t *testing.T) {
	str1 := []byte(`{}`)
	str2 := []byte(`{"value":""}`)

	var str strings.Builder
	maxLength := 5000
	for i := 0; i < maxLength+1; i++ {
		str.WriteString("a")
	}
	str3 := []byte(`{"value":"` + str.String() + `"}`)

	wantErr1 := "a provided secret (value) is missing"
	wantErr2 := "provided value is too large, character limit exceeds: " + strconv.Itoa(maxLength)

	var tests = []struct {
		input   []byte
		wantErr string
	}{
		{input: str1, wantErr: wantErr1},
		{input: str2, wantErr: wantErr1},
		{input: str3, wantErr: wantErr2},
	}

	for _, test := range tests {
		_, err := NewFromJSON(io.NopCloser(bytes.NewBuffer(test.input)))

		if err == nil {
			t.Errorf("incorrect value, should return an error")
		}

		if err.Error() != test.wantErr {
			t.Errorf("incorrect value, got: %s, want: %s", err.Error(), test.wantErr)
		}
	}
}

func TestToModel(t *testing.T) {
	createdAt := time.Now()
	expiresAt := time.Now()

	encrypted1, _ := security.Encrypt([]byte("value1"), encryptionKey)
	encrypted2, _ := security.Encrypt([]byte("value2"), encryptionKey)
	encrypted3, _ := security.Encrypt([]byte("value3"), encryptionKey)

	var tests = []struct {
		input *Secret
		want  *db.Secret
	}{
		{input: &Secret{Value: "value1"}, want: &db.Secret{Value: string(encrypted1)}},
		{input: &Secret{Value: "value2", Passphrase: "1234"}, want: &db.Secret{Value: string(encrypted2)}},
		{input: &Secret{Value: "value3", CreatedAt: createdAt, ExpiresAt: expiresAt}, want: &db.Secret{Value: string(encrypted3), CreatedAt: createdAt, ExpiresAt: expiresAt}},
	}

	for _, test := range tests {
		got := toModel(test.input, encryptionKey)

		decryptedGot, _ := security.Decrypt([]byte(got.Value), encryptionKey)
		decryptedWant, _ := security.Decrypt([]byte(test.want.Value), encryptionKey)
		if string(decryptedGot) != string(decryptedWant) {
			t.Errorf("incorrect value, got: %s, want: %s", string(decryptedGot), string(decryptedWant))
		}
		if !test.input.CreatedAt.IsZero() && got.CreatedAt != test.want.CreatedAt {
			t.Errorf("incorrect value, got: %v, want: %v", got.CreatedAt, test.want.CreatedAt)
		}
		if !test.input.ExpiresAt.IsZero() && got.ExpiresAt != test.want.ExpiresAt {
			t.Errorf("incorrect value, got: %v, want: %v", got.ExpiresAt, test.want.ExpiresAt)
		}
	}
}
