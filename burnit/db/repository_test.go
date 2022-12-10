package db

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// mockClient is a struct to test the different
// methods on SecretRepository.
type mockClient struct {
	mode string
}

var id1 = "507f1f77bcf86cd799439011"

func (c *mockClient) Connect(ctx context.Context) error {
	return nil
}

func (c *mockClient) Disconnect(ctx context.Context) error {
	return nil
}

func (c *mockClient) Find(ctx context.Context, id string) (*Secret, error) {
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

func (c *mockClient) Insert(ctx context.Context, m *Secret) (*Secret, error) {
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

func (c *mockClient) Delete(ctx context.Context, id string) (int64, error) {
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

func (c *mockClient) DeleteMany(ctx context.Context) (int64, error) {
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
	return &SecretRepository{
		client: &mockClient{mode: mode},
	}
}

func TestNewSecretRepository(t *testing.T) {
	want := &SecretRepository{
		client:  &mockClient{},
		timeout: time.Second * 5,
	}
	got := NewSecretRepository(&mockClient{}, &SecretRepositoryOptions{})

	if diff := cmp.Diff(want, got, cmp.AllowUnexported(mockClient{}, SecretRepository{})); diff != "" {
		t.Errorf("NewSecretRepository(%q, %q) = unexpected result (-want, +got)\n%s\n", &mockClient{}, &SecretRepositoryOptions{}, diff)
	}
}

func TestSecretRepositoryGet(t *testing.T) {
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

		res, err := repo.Get(test.input)
		if res != nil && res.Value != test.want.Value {
			t.Errorf("incorrect value, got: %v, want: %v", res.Value, test.want.Value)
		}
		if err != nil && err.Error() != "error in db" {
			t.Errorf("incorrect value, got: %v, want: %v", err, test.wantErr)
		}
	}
}

func TestSecretRepositoryCreate(t *testing.T) {
	var tests = []struct {
		inputMode string
		input     *Secret
		want      *Secret
		wantErr   error
	}{
		{inputMode: "insert-success", input: &Secret{Value: "value"}, want: &Secret{ID: id1}, wantErr: nil},
		{inputMode: "insert-success", input: &Secret{Value: "value"}, want: &Secret{ID: id1}, wantErr: nil},
		{inputMode: "insert-error", input: &Secret{Value: "value"}, want: nil, wantErr: errors.New("error in db")},
	}

	for _, test := range tests {
		repo := SetupRepository(test.inputMode)

		res, err := repo.Create(test.input)
		if res != nil && res.ID != id1 {
			t.Errorf("incorrect value, got: %v, want: %v", res.ID, id1)
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
