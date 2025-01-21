package inmem

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/RedeployAB/burnit/internal/db"
	dberrors "github.com/RedeployAB/burnit/internal/db/errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewSecretStore(t *testing.T) {
	got := NewSecretStore()
	if diff := cmp.Diff(&secretStore{secrets: map[string]db.Secret{}, mu: sync.RWMutex{}}, got, cmp.AllowUnexported(secretStore{}), cmpopts.IgnoreFields(secretStore{}, "mu")); diff != "" {
		t.Errorf("NewSecretStore() = unexpected result (-want +got)\n%s\n", diff)
	}
}

func TestSecretStore_Get(t *testing.T) {
	n := now()
	var tests = []struct {
		name  string
		input struct {
			secrets map[string]db.Secret
			id      string
		}
		want    db.Secret
		wantErr error
	}{
		{
			name: "Get secret",
			input: struct {
				secrets map[string]db.Secret
				id      string
			}{
				secrets: map[string]db.Secret{
					"test": {
						ID:        "test",
						Value:     "secret",
						ExpiresAt: n.Add(1),
					},
				},
				id: "test",
			},
			want: db.Secret{
				ID:        "test",
				Value:     "secret",
				ExpiresAt: n.Add(1),
			},
		},
		{
			name: "Secret not found",
			input: struct {
				secrets map[string]db.Secret
				id      string
			}{
				secrets: map[string]db.Secret{},
				id:      "test",
			},
			wantErr: dberrors.ErrSecretNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &secretStore{
				secrets: test.input.secrets,
				mu:      sync.RWMutex{},
			}

			got, gotErr := s.Get(context.Background(), test.input.id)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Get() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Get() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestSecretStore_Create(t *testing.T) {
	n := now()
	var tests = []struct {
		name  string
		input struct {
			secrets map[string]db.Secret
			secret  db.Secret
		}
		want    db.Secret
		wantErr error
	}{
		{
			name: "Create secret",
			input: struct {
				secrets map[string]db.Secret
				secret  db.Secret
			}{
				secrets: map[string]db.Secret{},
				secret: db.Secret{
					ID:        "test",
					Value:     "secret",
					ExpiresAt: n.Add(1),
				},
			},
			want: db.Secret{
				ID:        "test",
				Value:     "secret",
				ExpiresAt: n.Add(1),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &secretStore{
				secrets: test.input.secrets,
				mu:      sync.RWMutex{},
			}

			got, gotErr := s.Create(context.Background(), test.input.secret)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Create() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Create() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestSecretStore_Delete(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			secrets map[string]db.Secret
			id      string
		}
		wantErr error
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &secretStore{
				secrets: test.input.secrets,
				mu:      sync.RWMutex{},
			}

			gotErr := s.Delete(context.Background(), test.input.id)

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Delete() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestSecretStore_DeleteExpired(t *testing.T) {
	n := now()
	var tests = []struct {
		name  string
		input struct {
			secrets map[string]db.Secret
		}
		want    map[string]db.Secret
		wantErr error
	}{
		{
			name: "Delete expired secrets",
			input: struct {
				secrets map[string]db.Secret
			}{
				secrets: map[string]db.Secret{
					"test1": {
						ID:        "test1",
						Value:     "secret1",
						ExpiresAt: n.Add(-1),
					},
					"test2": {
						ID:        "test2",
						Value:     "secret2",
						ExpiresAt: n.Add(time.Hour),
					},
					"test3": {
						ID:        "test3",
						Value:     "secret3",
						ExpiresAt: n.Add(-1),
					},
				},
			},
			want: map[string]db.Secret{
				"test2": {
					ID:        "test2",
					Value:     "secret2",
					ExpiresAt: n.Add(time.Hour),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &secretStore{
				secrets: test.input.secrets,
				mu:      sync.RWMutex{},
			}

			gotErr := s.DeleteExpired(context.Background())

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("DeleteExpired() = unexpected error (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.want, s.secrets); diff != "" {
				t.Errorf("DeleteExpired() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}
