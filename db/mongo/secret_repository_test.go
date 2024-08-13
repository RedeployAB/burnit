package mongo

import (
	"context"
	"testing"
	"time"

	dberrors "github.com/RedeployAB/burnit/db/errors"
	"github.com/RedeployAB/burnit/db/models"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewSecretRepository(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			client  Client
			options []SecretRepositoryOption
		}
		want    *SecretRepository
		wantErr error
	}{
		{
			name: "new secret repository",
			input: struct {
				client  Client
				options []SecretRepositoryOption
			}{
				client: &mockMongoClient{},
			},
			want: &SecretRepository{
				client:             &mockMongoClient{},
				collection:         defaultSecretRepositoryCollection,
				settingsCollection: defaultSettingsCollection,
				timeout:            defaultSecretRepositoryTimeout,
			},
		},
		{
			name: "new secret repository - with options",
			input: struct {
				client  Client
				options []SecretRepositoryOption
			}{
				client: &mockMongoClient{},
				options: []SecretRepositoryOption{
					func(o *SecretRepositoryOptions) {
						o.Database = "test"
						o.Collection = "test"
						o.SettingsCollection = "test"
					},
				},
			},
			want: &SecretRepository{
				client:             &mockMongoClient{},
				collection:         "test",
				settingsCollection: "test",
				timeout:            defaultSecretRepositoryTimeout,
			},
		},
		{
			name: "new secret repository - nil client",
			input: struct {
				client  Client
				options []SecretRepositoryOption
			}{
				client: nil,
			},
			wantErr: ErrNilClient,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewSecretRepository(test.input.client, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(SecretRepository{}, mockMongoClient{})); diff != "" {
				t.Errorf("NewSecretRepository() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("NewSecretRepository() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestSecretRepository_Get(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			secrets []models.Secret
			id      string
			err     error
		}
		want    models.Secret
		wantErr error
	}{
		{
			name: "get secret",
			input: struct {
				secrets []models.Secret
				id      string
				err     error
			}{
				secrets: []models.Secret{
					{
						ID:    "1",
						Value: "secret",
					},
				},
				id: "1",
			},
			want: models.Secret{
				ID:    "1",
				Value: "secret",
			},
		},
		{
			name: "get secret - not found",
			input: struct {
				secrets []models.Secret
				id      string
				err     error
			}{
				secrets: []models.Secret{},
				id:      "1",
			},
			wantErr: dberrors.ErrSecretNotFound,
		},
		{
			name: "get secret - error",
			input: struct {
				secrets []models.Secret
				id      string
				err     error
			}{
				secrets: []models.Secret{},
				id:      "1",
				err:     errFindOne,
			},
			wantErr: errFindOne,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &SecretRepository{
				client: &mockMongoClient{
					secrets: test.input.secrets,
					err:     test.input.err,
				},
			}

			got, gotErr := repo.Get(context.Background(), test.input.id)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Get() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Get() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestSecretRepository_Create(t *testing.T) {
	date := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	now = func() time.Time {
		return date
	}

	var tests = []struct {
		name  string
		input struct {
			secrets []models.Secret
			secret  models.Secret
			err     error
		}
		want    models.Secret
		wantErr error
	}{
		{
			name: "create secret",
			input: struct {
				secrets []models.Secret
				secret  models.Secret

				err error
			}{
				secrets: []models.Secret{},
				secret: models.Secret{
					ID:    "1",
					Value: "secret",
				},
			},
			want: models.Secret{
				ID:    "1",
				Value: "secret",
			},
		},
		{
			name: "create secret - error",
			input: struct {
				secrets []models.Secret
				secret  models.Secret
				err     error
			}{
				secrets: []models.Secret{},
				secret: models.Secret{
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
			repo := &SecretRepository{
				client: &mockMongoClient{
					secrets: test.input.secrets,
					err:     test.input.err,
				},
			}

			got, gotErr := repo.Create(context.Background(), test.input.secret)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Create() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Create() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestSecretRepository_Delete(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			secrets []models.Secret
			id      string
			err     error
		}
		wantErr error
	}{
		{
			name: "delete secret",
			input: struct {
				secrets []models.Secret
				id      string
				err     error
			}{
				secrets: []models.Secret{
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
				secrets []models.Secret
				id      string
				err     error
			}{
				secrets: []models.Secret{},
				id:      "1",
			},
			wantErr: dberrors.ErrSecretNotFound,
		},
		{
			name: "delete secret - error",
			input: struct {
				secrets []models.Secret
				id      string
				err     error
			}{
				secrets: []models.Secret{},
				id:      "1",
				err:     errDeleteOne,
			},
			wantErr: errDeleteOne,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &SecretRepository{
				client: &mockMongoClient{
					secrets: test.input.secrets,
					err:     test.input.err,
				},
			}

			gotErr := repo.Delete(context.Background(), test.input.id)

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Delete() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestSecretRepository_DeleteExpired(t *testing.T) {
	date := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	now = func() time.Time {
		return date
	}

	var tests = []struct {
		name  string
		input struct {
			secrets []models.Secret
			err     error
		}
		wantErr error
	}{
		{
			name: "delete expired secrets",
			input: struct {
				secrets []models.Secret
				err     error
			}{
				secrets: []models.Secret{
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
				secrets []models.Secret
				err     error
			}{
				err: errDeleteMany,
			},
			wantErr: errDeleteMany,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &SecretRepository{
				client: &mockMongoClient{
					secrets: test.input.secrets,
					err:     test.input.err,
				},
			}

			gotErr := repo.DeleteExpired(context.Background())

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("DeleteExpired() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestSecretRepository_GetSettings(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			settings []models.Settings
			err      error
		}
		want    models.Settings
		wantErr error
	}{
		{
			name: "get settings",
			input: struct {
				settings []models.Settings
				err      error
			}{
				settings: []models.Settings{
					{
						Security: models.Security{
							ID:            "security",
							EncryptionKey: "test",
						},
					},
				},
			},
			want: models.Settings{
				Security: models.Security{
					ID:            "security",
					EncryptionKey: "test",
				},
			},
		},
		{
			name: "get settings - not found",
			input: struct {
				settings []models.Settings
				err      error
			}{},
			wantErr: dberrors.ErrSettingsNotFound,
		},
		{
			name: "get settings - error",
			input: struct {
				settings []models.Settings
				err      error
			}{
				err: errFindOne,
			},
			wantErr: errFindOne,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &SecretRepository{
				client: &mockMongoClient{
					settings: test.input.settings,
					err:      test.input.err,
				},
			}

			got, gotErr := repo.GetSettings(context.Background())

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetSettings() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("GetSettings() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestSecretRepository_CreateSettings(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			settings   []models.Settings
			inSettings models.Settings
			err        error
		}
		want    models.Settings
		wantErr error
	}{
		{
			name: "create settings",
			input: struct {
				settings   []models.Settings
				inSettings models.Settings
				err        error
			}{
				settings: []models.Settings{},
				inSettings: models.Settings{
					Security: models.Security{
						EncryptionKey: "test",
					},
				},
			},
			want: models.Settings{
				Security: models.Security{
					ID:            "security",
					EncryptionKey: "test",
				},
			},
		},
		{
			name: "create settings - error",
			input: struct {
				settings   []models.Settings
				inSettings models.Settings
				err        error
			}{
				settings: []models.Settings{},
				inSettings: models.Settings{
					Security: models.Security{
						EncryptionKey: "test",
					},
				},
				err: errInsertOne,
			},
			wantErr: errInsertOne,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &SecretRepository{
				client: &mockMongoClient{
					settings: test.input.settings,
					err:      test.input.err,
				},
			}

			got, gotErr := repo.CreateSettings(context.Background(), test.input.inSettings)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CreateSettings() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("CreateSettings() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestSecretRepository_UpdateSettings(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			settings   []models.Settings
			inSettings models.Settings
			err        error
		}
		want    models.Settings
		wantErr error
	}{
		{
			name: "update settings",
			input: struct {
				settings   []models.Settings
				inSettings models.Settings
				err        error
			}{
				settings: []models.Settings{
					{
						Security: models.Security{
							ID:            "security",
							EncryptionKey: "test",
						},
					},
				},
				inSettings: models.Settings{
					Security: models.Security{
						ID:            "security",
						EncryptionKey: "test-updated",
					},
				},
			},
			want: models.Settings{
				Security: models.Security{
					ID:            "security",
					EncryptionKey: "test-updated",
				},
			},
		},
		{
			name: "update settings - error",
			input: struct {
				settings   []models.Settings
				inSettings models.Settings
				err        error
			}{
				settings: []models.Settings{},
				inSettings: models.Settings{
					Security: models.Security{
						ID:            "security",
						EncryptionKey: "test-updated",
					},
				},
				err: errUpdateOne,
			},
			wantErr: errUpdateOne,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &SecretRepository{
				client: &mockMongoClient{
					settings: test.input.settings,
					err:      test.input.err,
				},
			}

			got, gotErr := repo.UpdateSettings(context.Background(), test.input.inSettings)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("UpdateSettings() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("UpdateSettings() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}
