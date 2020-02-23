package dto

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/RedeployAB/burnit/common/security"
)

func TestNewSecret(t *testing.T) {
	expectedDay1 := time.Now().AddDate(0, 0, 7)
	expectedDay2 := time.Now().AddDate(0, 0, 3)
	tests := []struct {
		input []byte
		want  *Secret
	}{
		{
			[]byte(`{"secret":"value1"}`),
			&Secret{Secret: "value1", ExpiresAt: expectedDay1},
		},
		{
			[]byte(`{"secret":"value2","ttl":4320}`),
			&Secret{Secret: "value2", ExpiresAt: expectedDay2},
		},
	}

	for _, test := range tests {
		b := bytes.NewBuffer(test.input)
		dto, err := NewSecret(ioutil.NopCloser(b))
		if err != nil {
			t.Errorf("error in test: %v", err)
		}

		if dto.Secret != test.want.Secret {
			t.Errorf("incorrect value, got: %s, want: %s", dto.Secret, test.want.Secret)
		}
		expectedDay := test.want.ExpiresAt.Day()
		if dto.ExpiresAt.Day() != expectedDay {
			t.Errorf("incorrect value, got: %d, want: %d", dto.ExpiresAt.Day(), expectedDay)
		}
	}
}

func TestNewSecretFail(t *testing.T) {
	// Malformed JSON.
	b := bytes.NewBuffer([]byte(`{"secret":"value}`))
	_, err := NewSecret(ioutil.NopCloser(b))
	if err == nil {
		t.Errorf("error in test, test should fail.")
	}
}

func TestVerify(t *testing.T) {
	passphrase := "passphrase"
	hash := security.Hash(passphrase)

	dto := &Secret{
		Passphrase: hash,
	}

	verify1 := dto.Verify(passphrase)
	if !verify1 {
		t.Errorf("incorrect result, got: %v, want: true", verify1)
	}

	verify2 := dto.Verify("notpassphrase")
	if verify2 {
		t.Errorf("incorrect result, got: %v, want: false", verify2)
	}
}
