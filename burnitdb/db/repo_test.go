package db

import (
	"errors"
	"testing"

	"github.com/RedeployAB/burnit/burnitdb/internal/models"
	"github.com/RedeployAB/burnit/common/security"
)

type mockClient struct{}

func (m *mockClient) Database(name string) *mockDatabase {
	return &mockDatabase{}
}

type mockDatabase struct{}

func (d *mockDatabase) Collection(name string) *mockCollection {
	return &mockCollection{}
}

type mockCollection struct {
	mode string
}

var id1 = "507f1f77bcf86cd799439011"
var encryptionKey = "abcdefg"
var encrypted, _ = security.Encrypt([]byte("value"), encryptionKey)

func (c *mockCollection) FindOne(id string) (*models.Secret, error) {
	var model *models.Secret
	var err error

	switch c.mode {
	case "find-success":
		model = &models.Secret{ID: id1, Secret: string(encrypted)}
		err = nil
	case "find-not-found":
		model = nil
		err = nil
	case "find-error":
		model = nil
		err = errors.New("error in db")
	}
	return model, err
}

func (c *mockCollection) InsertOne(m *models.Secret) (*models.Secret, error) {
	var model *models.Secret
	var err error

	switch c.mode {
	case "insert-success":
		model = &models.Secret{ID: id1}
		err = nil
	case "insert-error":
		model = nil
		err = errors.New("error in db")
	}

	return model, err
}

func (c *mockCollection) DeleteOne(id string) (int64, error) {
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

func (c *mockCollection) DeleteMany() (int64, error) {
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
		EncryptionKey: encryptionKey,
		HashMethod:    "bcrypt",
	}
	return &SecretRepository{
		collection: &mockCollection{mode: mode},
		options:    opts,
		hash:       security.Bcrypt,
	}
}

func TestSecretRepositoryFind(t *testing.T) {
	var tests = []struct {
		inputMode string
		input     string
		wanted    *models.Secret
		wantedErr error
	}{
		{inputMode: "find-success", input: id1, wanted: &models.Secret{ID: id1, Secret: "value"}, wantedErr: nil},
		{inputMode: "find-not-found", input: id1, wanted: nil, wantedErr: nil},
		{inputMode: "find-invalid-oid", input: "1234", wanted: nil, wantedErr: nil},
		{inputMode: "find-error", input: id1, wanted: nil, wantedErr: errors.New("error in db")},
	}

	for _, test := range tests {
		repo := SetupRepository(test.inputMode)

		res, err := repo.Find(test.input)
		if res != nil && res.Secret != test.wanted.Secret {
			t.Errorf("incorrect value, got: %v, want: %v", res.Secret, test.wanted.Secret)
		}
		if err != nil && err.Error() != "error in db" {
			t.Errorf("incorrect value, got: %v, want: %v", err, test.wantedErr)
		}
	}
}

func TestSecretRepositoryInsert(t *testing.T) {
	var tests = []struct {
		inputMode string
		input     *models.Secret
		wanted    *models.Secret
		wantedErr error
	}{
		{inputMode: "insert-success", input: &models.Secret{Secret: "value"}, wanted: &models.Secret{ID: id1}, wantedErr: nil},
		{inputMode: "insert-success", input: &models.Secret{Secret: "value", Passphrase: "passphrase"}, wanted: &models.Secret{ID: id1, Passphrase: security.ToMD5("passphrase")}, wantedErr: nil},
		{inputMode: "insert-error", input: &models.Secret{Secret: "value"}, wanted: nil, wantedErr: errors.New("error in db")},
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
			t.Errorf("incorrect value, got: %v, want: %v", err, test.wantedErr)
		}
	}
}

func TestSecretRepositoryDelete(t *testing.T) {
	var tests = []struct {
		inputMode string
		input     string
		wanted    int64
		wantedErr error
	}{
		{inputMode: "delete-success", input: id1, wanted: 1, wantedErr: nil},
		{inputMode: "delete-not-found", input: id1, wanted: 0, wantedErr: nil},
		{inputMode: "delete-error", input: id1, wanted: 0, wantedErr: errors.New("error in db")},
	}

	for _, test := range tests {
		repo := SetupRepository(test.inputMode)

		res, err := repo.Delete(test.input)
		if res != test.wanted {
			t.Errorf("incorrect value, got: %d, want: %d", res, test.wanted)
		}
		if err != nil && err.Error() != "error in db" {
			t.Errorf("incorrect value, got: %v, want: %v", err, test.wantedErr)
		}
	}
}

func TestSecretRepositoryDeleteMany(t *testing.T) {
	var tests = []struct {
		inputMode string
		wanted    int64
		wantedErr error
	}{
		{inputMode: "delete-many-success", wanted: 2, wantedErr: nil},
		{inputMode: "delete-many-not-found", wanted: 0, wantedErr: nil},
		{inputMode: "delete-many-error", wanted: 0, wantedErr: errors.New("error in db")},
	}

	for _, test := range tests {
		repo := SetupRepository(test.inputMode)

		res, err := repo.DeleteExpired()
		if res != test.wanted {
			t.Errorf("incorrect value, got: %d, want: %d", res, test.wanted)
		}
		if err != nil && err.Error() != "error in db" {
			t.Errorf("incorrect value, got: %v, want: %v", err, test.wantedErr)
		}
	}
}
