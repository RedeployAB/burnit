package secret

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/RedeployAB/burnit/burnitdb/db"
	"github.com/RedeployAB/burnit/common/security"
)

var correctPassphrase = security.Bcrypt("passphrase")
var id1 = "507f1f77bcf86cd799439011"
var apiKey = "ABCDEF"

// Mock to handle repository answers in handler tests.
type mockSecretRepository struct {
	action string
	mode   string
}

func (r *mockSecretRepository) Find(id string) (*db.Secret, error) {
	// Return different results based on underlying structs
	// state.
	var secret *db.Secret
	var err error

	switch r.mode {
	case "find-success":
		secret = &db.Secret{ID: id1, Secret: "value"}
		err = nil
	case "find-not-found":
		secret = nil
		err = nil
	case "find-error":
		secret = nil
		err = errors.New("find error")
	}
	return secret, err
}

func (r *mockSecretRepository) Insert(s *db.Secret) (*db.Secret, error) {
	var secret *db.Secret
	var err error

	switch r.mode {
	case "insert-success":
		secret = &db.Secret{ID: id1, Secret: "value"}
		err = nil
	case "insert-error":
		secret = nil
		err = errors.New("insert error")
	}
	return secret, err
}

func (r *mockSecretRepository) Delete(id string) (int64, error) {
	var result int64
	var err error
	switch r.mode {
	case "delete-success":
		result = 1
		err = nil
	case "delete-not-found":
		result = 0
		err = nil
	case "delete-error":
		result = -10
		err = errors.New("delete error")
	}
	return result, err
}

func (r *mockSecretRepository) DeleteExpired() (int64, error) {
	var result int64
	var err error
	switch r.mode {
	case "delete-expired-success":
		result = 1
		err = nil
	case "delete-expired-not-found":
		result = 0
		err = nil
	case "delete-expired-error":
		result = -10
		err = errors.New("delete error")
	}
	return result, err
}

func SetupService(action, mode string) Service {
	repo := &mockSecretRepository{action: action, mode: mode}
	return NewService(repo)
}

func TestNewService(t *testing.T) {
	repo := &mockSecretRepository{action: "", mode: ""}
	service := NewService(repo)

	if service == nil {
		t.Errorf("error in creating service")
	}
}

func TestServiceGet(t *testing.T) {
	var tests = []struct {
		mode    string
		input   string
		want    *Secret
		wantErr error
	}{
		{mode: "find-success", input: id1, want: &Secret{ID: id1, Secret: "value"}, wantErr: nil},
		{mode: "find-not-found", input: id1, want: nil, wantErr: nil},
		{mode: "find-error", input: id1, want: nil, wantErr: errors.New("find error")},
	}

	for _, test := range tests {
		svc := SetupService("find", test.mode)
		sec, err := svc.Get(test.input)
		if sec != nil && sec.Secret != test.want.Secret {
			t.Errorf("incorrect value, got: %v, want: %v", sec.Secret, test.want.Secret)
		}
		if err != nil && err.Error() != test.wantErr.Error() {
			t.Errorf("incorrect value, got: %v, want: %v", err.Error(), test.wantErr.Error())
		}
	}
}

func TestServiceCreate(t *testing.T) {
	var tests = []struct {
		mode    string
		input   *Secret
		want    *Secret
		wantErr error
	}{
		{mode: "insert-success", input: &Secret{Secret: "value"}, want: &Secret{ID: id1}, wantErr: nil},
		{mode: "insert-error", input: &Secret{Secret: "value"}, want: nil, wantErr: errors.New("insert error")},
	}

	for _, test := range tests {
		svc := SetupService("create", test.mode)
		sec, err := svc.Create(test.input)
		if sec != nil && sec.ID != test.want.ID {
			t.Errorf("incorrect value, got: %v, want: %v", sec.ID, test.want.ID)
		}
		if err != nil && err.Error() != test.wantErr.Error() {
			t.Errorf("incorrect value, got: %v, want: %v", err.Error(), test.wantErr.Error())
		}
	}
}

func TestServiceDelete(t *testing.T) {
	var tests = []struct {
		mode    string
		input   string
		want    int64
		wantErr error
	}{
		{mode: "delete-success", input: id1, want: 1, wantErr: nil},
		{mode: "delete-not-found", input: id1, want: 0, wantErr: nil},
		{mode: "delete-error", input: id1, want: 0, wantErr: errors.New("delete error")},
	}

	for _, test := range tests {
		svc := SetupService("delete", test.mode)
		deleted, err := svc.Delete(test.input)
		if deleted != test.want {
			t.Errorf("incorrect value, got: %d, want: %d", deleted, test.want)
		}
		if err != nil && err.Error() != test.wantErr.Error() {
			t.Errorf("incorrect value, got: %v, want: %v", err.Error(), test.wantErr.Error())
		}
	}
}

func TestServiceDeleteExpired(t *testing.T) {
	var tests = []struct {
		mode    string
		want    int64
		wantErr error
	}{
		{mode: "delete-expired-success", want: 1, wantErr: nil},
		{mode: "delete-expired-not-found", want: 0, wantErr: nil},
		{mode: "delete-expired-error", want: 0, wantErr: errors.New("delete error")},
	}

	for _, test := range tests {
		svc := SetupService("delete-expired", test.mode)
		deleted, err := svc.DeleteExpired()
		if deleted != test.want {
			t.Errorf("incorrect value, got: %d, want: %d", deleted, test.want)
		}
		if err != nil && err.Error() != test.wantErr.Error() {
			t.Errorf("incorrect value, got: %v, want: %v", err.Error(), test.wantErr.Error())
		}
	}
}

func TestNewFromJSON(t *testing.T) {
	str1 := []byte(`{"secret":"value1"}`)
	str2 := []byte(`{"secret":"value2","passphrase":"1234"}`)
	str3 := []byte(`{"secret":"value3","ttl":4320}`)
	strMalformed1 := []byte(`{"secret":`)

	expectedDay1 := time.Now().AddDate(0, 0, 7)
	expectedDay2 := time.Now().AddDate(0, 0, 3)

	var tests = []struct {
		input   []byte
		want    *Secret
		wantErr error
	}{
		{input: str1, want: &Secret{Secret: "value1", Passphrase: "", ExpiresAt: expectedDay1}, wantErr: nil},
		{input: str2, want: &Secret{Secret: "value2", Passphrase: "1234", ExpiresAt: expectedDay1}, wantErr: nil},
		{input: str3, want: &Secret{Secret: "value3", Passphrase: "", ExpiresAt: expectedDay2}, wantErr: nil},
		{input: strMalformed1, want: nil, wantErr: errors.New("unexpected EOF")},
	}

	for _, test := range tests {
		got, err := NewFromJSON(ioutil.NopCloser(bytes.NewBuffer(test.input)))
		if got != nil && got.Secret != test.want.Secret {
			t.Errorf("incorrect value, got: %s, want: %s", got.Secret, test.want.Secret)
		}
		if got != nil && got.Passphrase != test.want.Passphrase {
			t.Errorf("incorrect value, got: %s, want: %s", got.Passphrase, test.want.Passphrase)
		}
		if err != nil && err.Error() != test.wantErr.Error() {
			t.Errorf("incorrect value, got: %v, want: %v", err.Error(), test.wantErr.Error())
		}
	}
}

func TestToModel(t *testing.T) {
	createdAt := time.Now()
	expiresAt := time.Now()

	var tests = []struct {
		input *Secret
		want  *db.Secret
	}{
		{input: &Secret{Secret: "value1"}, want: &db.Secret{Secret: "value1"}},
		{input: &Secret{Secret: "value2", Passphrase: "1234"}, want: &db.Secret{Secret: "value2", Passphrase: "1234"}},
		{input: &Secret{Secret: "value3", CreatedAt: createdAt, ExpiresAt: expiresAt}, want: &db.Secret{Secret: "value3", CreatedAt: createdAt, ExpiresAt: expiresAt}},
	}

	for _, test := range tests {
		got := toModel(test.input)
		if got.Secret != test.want.Secret {
			t.Errorf("incorrect value, got: %s, want: %s", got.Secret, test.want.Secret)
		}
		if got.Passphrase != test.want.Passphrase {
			t.Errorf("incorrect value, got: %s, want: %s", got.Passphrase, test.want.Passphrase)
		}
		if !test.input.CreatedAt.IsZero() && got.CreatedAt != test.want.CreatedAt {
			t.Errorf("incorrect value, got: %v, want: %v", got.CreatedAt, test.want.CreatedAt)
		}
		if !test.input.ExpiresAt.IsZero() && got.ExpiresAt != test.want.ExpiresAt {
			t.Errorf("incorrect value, got: %v, want: %v", got.ExpiresAt, test.want.ExpiresAt)
		}
	}
}
