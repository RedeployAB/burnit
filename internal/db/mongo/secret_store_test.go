package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/RedeployAB/burnit/internal/db"
	dberrors "github.com/RedeployAB/burnit/internal/db/errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewSecretStore(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			client  Client
			options []SecretStoreOption
		}
		want    *secretStore
		wantErr error
	}{
		{
			name: "new secret store",
			input: struct {
				client  Client
				options []SecretStoreOption
			}{
				client: &stubMongoClient{},
			},
			want: &secretStore{
				client:     &stubMongoClient{},
				collection: defaultSecretStoreCollection,
				timeout:    defaultSecretStoreTimeout,
			},
		},
		{
			name: "new secret store - with options",
			input: struct {
				client  Client
				options []SecretStoreOption
			}{
				client: &stubMongoClient{},
				options: []SecretStoreOption{
					func(o *SecretStoreOptions) {
						o.Database = "test"
						o.Collection = "test"
					},
				},
			},
			want: &secretStore{
				client:     &stubMongoClient{},
				collection: "test",
				timeout:    defaultSecretStoreTimeout,
			},
		},
		{
			name: "new secret store - nil client",
			input: struct {
				client  Client
				options []SecretStoreOption
			}{
				client: nil,
			},
			wantErr: ErrNilClient,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewSecretStore(test.input.client, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(secretStore{}, stubMongoClient{}), cmpopts.IgnoreFields(secretStore{}, "createSecret")); diff != "" {
				t.Errorf("NewSecretStore() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("NewSecretStore() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestSecretStore_Get(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			secrets []db.Secret
			id      string
			err     error
		}
		want    db.Secret
		wantErr error
	}{
		{
			name: "get secret",
			input: struct {
				secrets []db.Secret
				id      string
				err     error
			}{
				secrets: []db.Secret{
					{
						ID:    "1",
						Value: "secret",
					},
				},
				id: "1",
			},
			want: db.Secret{
				ID:    "1",
				Value: "secret",
			},
		},
		{
			name: "get secret - not found",
			input: struct {
				secrets []db.Secret
				id      string
				err     error
			}{
				secrets: []db.Secret{},
				id:      "1",
			},
			wantErr: dberrors.ErrSecretNotFound,
		},
		{
			name: "get secret - error",
			input: struct {
				secrets []db.Secret
				id      string
				err     error
			}{
				secrets: []db.Secret{},
				id:      "1",
				err:     errFindOne,
			},
			wantErr: errFindOne,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := &secretStore{
				client: &stubMongoClient{
					secrets: test.input.secrets,
					err:     test.input.err,
				},
			}

			got, gotErr := store.Get(context.Background(), test.input.id)

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
	date := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	now = func() time.Time {
		return date
	}

	var tests = []struct {
		name  string
		input struct {
			secrets []db.Secret
			secret  db.Secret
			err     error
		}
		want    db.Secret
		wantErr error
	}{
		{
			name: "create secret",
			input: struct {
				secrets []db.Secret
				secret  db.Secret

				err error
			}{
				secrets: []db.Secret{},
				secret: db.Secret{
					ID:    "1",
					Value: "secret",
				},
			},
			want: db.Secret{
				ID:    "1",
				Value: "secret",
			},
		},
		{
			name: "create secret - error",
			input: struct {
				secrets []db.Secret
				secret  db.Secret
				err     error
			}{
				secrets: []db.Secret{},
				secret: db.Secret{
					ID:    "1",
					Value: "secret",
				},
				err: errInsertOne,
			},
			wantErr: errInsertOne,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := &secretStore{
				client: &stubMongoClient{
					secrets: test.input.secrets,
					err:     test.input.err,
				},
			}
			setCreateSecret(store)

			got, gotErr := store.Create(context.Background(), test.input.secret)

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
			secrets []db.Secret
			id      string
			err     error
		}
		wantErr error
	}{
		{
			name: "delete secret",
			input: struct {
				secrets []db.Secret
				id      string
				err     error
			}{
				secrets: []db.Secret{
					{
						ID:    "1",
						Value: "secret",
					},
				},
				id: "1",
			},
		},
		{
			name: "delete secret - not found",
			input: struct {
				secrets []db.Secret
				id      string
				err     error
			}{
				secrets: []db.Secret{},
				id:      "1",
			},
			wantErr: dberrors.ErrSecretNotFound,
		},
		{
			name: "delete secret - error",
			input: struct {
				secrets []db.Secret
				id      string
				err     error
			}{
				secrets: []db.Secret{},
				id:      "1",
				err:     errDeleteOne,
			},
			wantErr: errDeleteOne,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := &secretStore{
				client: &stubMongoClient{
					secrets: test.input.secrets,
					err:     test.input.err,
				},
			}

			gotErr := store.Delete(context.Background(), test.input.id)

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Delete() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestSecretStore_DeleteExpired(t *testing.T) {
	date := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	now = func() time.Time {
		return date
	}

	var tests = []struct {
		name  string
		input struct {
			secrets []db.Secret
			err     error
		}
		wantErr error
	}{
		{
			name: "delete expired secrets",
			input: struct {
				secrets []db.Secret
				err     error
			}{
				secrets: []db.Secret{
					{
						ID:        "1",
						Value:     "secret",
						ExpiresAt: date.Add(-time.Hour * 2),
					},
					{
						ID:        "2",
						Value:     "secret",
						ExpiresAt: date.Add(-time.Hour * 2),
					},
				},
			},
		},
		{
			name: "delete expired secrets - error",
			input: struct {
				secrets []db.Secret
				err     error
			}{
				err: errDeleteMany,
			},
			wantErr: errDeleteMany,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := &secretStore{
				client: &stubMongoClient{
					secrets: test.input.secrets,
					err:     test.input.err,
				},
			}

			gotErr := store.DeleteExpired(context.Background())

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("DeleteExpired() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}
