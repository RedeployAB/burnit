package db

import (
	"errors"
	"testing"

	"github.com/RedeployAB/burnit/common/security"
)

// mockClient is a struct to test the different
// methods on SecretRepository.
type mockClient struct {
	mode string
}

var id1 = "507f1f77bcf86cd799439011"

func (c *mockClient) FindOne(id string) (*Secret, error) {
	var secret *Secret
	var err error

	switch c.mode {
	case "find-success":
		secret = &Secret{ID: id1, Value: "value"}
		err = nil
	case "find-not-found":
		secret = nil
		err = nil
	case "find-error":
		secret = nil
		err = errors.New("error in db")
	}
	return secret, err
}

func (c *mockClient) InsertOne(m *Secret) (*Secret, error) {
	var secret *Secret
	var err error

	switch c.mode {
	case "insert-success":
		secret = &Secret{ID: id1}
		err = nil
	case "insert-error":
		secret = nil
		err = errors.New("error in db")
	}

	return secret, err
}

func (c *mockClient) DeleteOne(id string) (int64, error) {
	var deleted int64
	var err error

	switch c.mode {
	case "delete-success":
		deleted = 1
		err = nil
	case "delete-not-found":
		deleted = 0
		err = nil
	case "delete-error":
		deleted = 0
		err = errors.New("error in db")
	}

	return deleted, err
}

func (c *mockClient) DeleteMany() (int64, error) {
	var deleted int64
	var err error

	switch c.mode {
	case "delete-many-success":
		deleted = 2
		err = nil
	case "delete-many-not-found":
		deleted = 0
		err = nil
	case "delete-many-error":
		deleted = 0
		err = errors.New("error in db")
	}
	return deleted, err
}

func SetupRepository(mode string) *SecretRepository {
	opts := &SecretRepositoryOptions{
		HashMethod: "bcrypt",
	}
	return &SecretRepository{
		db:      &mockClient{mode: mode},
		options: opts,
		hash:    security.Bcrypt,
	}
}

func TestNewSecretRepository(t *testing.T) {
	// Test with redis driver and md5.
	opts := &SecretRepositoryOptions{
		Driver:     "redis",
		HashMethod: "md5",
	}

	rClient := &redisClient{}
	repo := NewSecretRepository(rClient, opts)

	expectedHashMethod := "md5"
	expectedDriver := "redis"

	if repo.options.HashMethod != expectedHashMethod {
		t.Errorf("incorrect hash method, got: %s, want: %s", repo.options.HashMethod, expectedHashMethod)
	}
	if repo.options.Driver != expectedDriver {
		t.Errorf("incorrect driver, got: %s, want: %s", repo.options.Driver, expectedDriver)
	}
	// Test with mongo driver and bcrypt.
	opts = &SecretRepositoryOptions{
		Driver:     "redis",
		HashMethod: "bcrypt",
	}

	repo = NewSecretRepository(rClient, opts)

	expectedHashMethod = "bcrypt"
	if repo.options.HashMethod != expectedHashMethod {
		t.Errorf("incorrect hash method, got: %s, want: %s", repo.options.HashMethod, expectedHashMethod)
	}
	if repo.options.Driver != expectedDriver {
		t.Errorf("incorrect driver, got: %s, want: %s", repo.options.Driver, expectedDriver)
	}
}

func TestSecretRepositoryFind(t *testing.T) {
	var tests = []struct {
		inputMode string
		input     string
		want      *Secret
		wantErr   error
	}{
		{inputMode: "find-success", input: id1, want: &Secret{ID: id1, Value: "value"}, wantErr: nil},
		{inputMode: "find-not-found", input: id1, want: nil, wantErr: nil},
		{inputMode: "find-invalid-oid", input: "1234", want: nil, wantErr: nil},
		{inputMode: "find-error", input: id1, want: nil, wantErr: errors.New("error in db")},
	}

	for _, test := range tests {
		repo := SetupRepository(test.inputMode)

		res, err := repo.Find(test.input)
		if res != nil && res.Value != test.want.Value {
			t.Errorf("incorrect value, got: %v, want: %v", res.Value, test.want.Value)
		}
		if err != nil && err.Error() != "error in db" {
			t.Errorf("incorrect value, got: %v, want: %v", err, test.wantErr)
		}
	}
}

func TestSecretRepositoryInsert(t *testing.T) {
	var tests = []struct {
		inputMode string
		input     *Secret
		want      *Secret
		wantErr   error
	}{
		{inputMode: "insert-success", input: &Secret{Value: "value"}, want: &Secret{ID: id1}, wantErr: nil},
		{inputMode: "insert-success", input: &Secret{Value: "value", Passphrase: "passphrase"}, want: &Secret{ID: id1, Passphrase: security.ToMD5("passphrase")}, wantErr: nil},
		{inputMode: "insert-error", input: &Secret{Value: "value"}, want: nil, wantErr: errors.New("error in db")},
	}

	for _, test := range tests {
		repo := SetupRepository(test.inputMode)

		res, err := repo.Insert(test.input)
		if res != nil && res.ID != id1 {
			t.Errorf("incorrect value, got: %v, want: %v", res.ID, id1)
		}
		if res != nil && res.Passphrase != "" && res.Passphrase != security.ToMD5("passphrase") {
			t.Errorf("incorrect value, got: %s, want: %s", res.Passphrase, security.ToMD5("passphrase"))
		}
		if err != nil && err.Error() != "error in db" {
			t.Errorf("incorrect value, got: %v, want: %v", err, test.wantErr)
		}
	}
}

func TestSecretRepositoryDelete(t *testing.T) {
	var tests = []struct {
		inputMode string
		input     string
		want      int64
		wantErr   error
	}{
		{inputMode: "delete-success", input: id1, want: 1, wantErr: nil},
		{inputMode: "delete-not-found", input: id1, want: 0, wantErr: nil},
		{inputMode: "delete-error", input: id1, want: 0, wantErr: errors.New("error in db")},
	}

	for _, test := range tests {
		repo := SetupRepository(test.inputMode)

		res, err := repo.Delete(test.input)
		if res != test.want {
			t.Errorf("incorrect value, got: %d, want: %d", res, test.want)
		}
		if err != nil && err.Error() != "error in db" {
			t.Errorf("incorrect value, got: %v, want: %v", err, test.wantErr)
		}
	}
}

func TestSecretRepositoryDeleteMany(t *testing.T) {
	var tests = []struct {
		inputMode string
		want      int64
		wantErr   error
	}{
		{inputMode: "delete-many-success", want: 2, wantErr: nil},
		{inputMode: "delete-many-not-found", want: 0, wantErr: nil},
		{inputMode: "delete-many-error", want: 0, wantErr: errors.New("error in db")},
	}

	for _, test := range tests {
		repo := SetupRepository(test.inputMode)

		res, err := repo.DeleteExpired()
		if res != test.want {
			t.Errorf("incorrect value, got: %d, want: %d", res, test.want)
		}
		if err != nil && err.Error() != "error in db" {
			t.Errorf("incorrect value, got: %v, want: %v", err, test.wantErr)
		}
	}
}
