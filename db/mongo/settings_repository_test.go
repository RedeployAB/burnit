package mongo

import (
	"context"
	"testing"

	dberrors "github.com/RedeployAB/burnit/db/errors"
	"github.com/RedeployAB/burnit/db/models"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewSettingsRepository(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			client  Client
			options []SettingsRepositoryOption
		}
		want    *SettingsRepository
		wantErr error
	}{
		{
			name: "new settings repository",
			input: struct {
				client  Client
				options []SettingsRepositoryOption
			}{
				client: &mockMongoClient{},
			},
			want: &SettingsRepository{
				client:     &mockMongoClient{},
				collection: defaultSettingsRepositoryCollection,
				timeout:    defaultSettingsRepositoryTimeout,
			},
		},
		{
			name: "new settings repository - with options",
			input: struct {
				client  Client
				options []SettingsRepositoryOption
			}{
				client: &mockMongoClient{},
				options: []SettingsRepositoryOption{
					func(o *SettingsRepositoryOptions) {
						o.Database = "test"
						o.Collection = "test"
					},
				},
			},
			want: &SettingsRepository{
				client:     &mockMongoClient{},
				collection: "test",
				timeout:    defaultSettingsRepositoryTimeout,
			},
		},
		{
			name: "new settings repository - nil client",
			input: struct {
				client  Client
				options []SettingsRepositoryOption
			}{
				client: nil,
			},
			want:    nil,
			wantErr: ErrNilClient,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewSettingsRepository(test.input.client, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(SettingsRepository{}, mockMongoClient{})); diff != "" {
				t.Errorf("NewSettingsRepository() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("NewSettingsRepository() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestSettingsRepository_Get(t *testing.T) {
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
			repo := &SettingsRepository{
				client: &mockMongoClient{
					settings: test.input.settings,
					err:      test.input.err,
				},
			}

			got, gotErr := repo.Get(context.Background())

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Get() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Get() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestSettingsRepository_Create(t *testing.T) {
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
			repo := &SettingsRepository{
				client: &mockMongoClient{
					settings: test.input.settings,
					err:      test.input.err,
				},
			}

			got, gotErr := repo.Create(context.Background(), test.input.inSettings)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Create() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Create() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestSettingsRepository(t *testing.T) {
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
			repo := &SettingsRepository{
				client: &mockMongoClient{
					settings: test.input.settings,
					err:      test.input.err,
				},
			}

			got, gotErr := repo.Update(context.Background(), test.input.inSettings)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Update() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Update() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}
