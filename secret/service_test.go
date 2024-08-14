package secret

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RedeployAB/burnit/db"
	dberrors "github.com/RedeployAB/burnit/db/errors"
	"github.com/RedeployAB/burnit/db/models"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewService(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			secrets db.SecretRepository
			options []ServiceOption
		}
		want    *service
		wantErr error
	}{
		{
			name: "new service",
			input: struct {
				secrets db.SecretRepository
				options []ServiceOption
			}{
				secrets: &mockSecretRepository{},
			},
			want: &service{
				secrets: &mockSecretRepository{},
				timeout: defaultTimeout,
			},
		},
		{
			name: "new service - with options",
			input: struct {
				secrets db.SecretRepository
				options []ServiceOption
			}{
				secrets: &mockSecretRepository{},
				options: []ServiceOption{
					func(s *service) {
						s.timeout = 30 * time.Second
					},
				},
			},
			want: &service{
				secrets: &mockSecretRepository{},
				timeout: 30 * time.Second,
			},
		},
		{
			name: "new service - nil repository",
			input: struct {
				secrets db.SecretRepository
				options []ServiceOption
			}{
				secrets: nil,
			},
			wantErr: ErrNilRepository,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewService(test.input.secrets, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(service{}, mockSecretRepository{})); diff != "" {
				t.Errorf("NewService() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("NewService() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			secrets db.SecretRepository
			id      string
			key     string
		}
		want    Secret
		wantErr error
	}{
		{
			name: "get secret",
			input: struct {
				secrets db.SecretRepository
				id      string
				key     string
			}{
				secrets: &mockSecretRepository{
					secrets: []models.Secret{
						{
							ID: "1",
							Value: func() string {
								v, _ := encrypt("secret", "key")
								return v
							}(),
							ExpiresAt: now().Add(1 * time.Hour),
						},
					},
				},
				id:  "1",
				key: "key",
			},
			want: Secret{
				ID:    "1",
				Value: "secret",
			},
		},
		{
			name: "get secret - not found",
			input: struct {
				secrets db.SecretRepository
				id      string
				key     string
			}{
				secrets: &mockSecretRepository{},
				id:      "1",
			},
			wantErr: ErrSecretNotFound,
		},
		{
			name: "get secret - expired",
			input: struct {
				secrets db.SecretRepository
				id      string
				key     string
			}{
				secrets: &mockSecretRepository{
					secrets: []models.Secret{
						{
							ID: "1",
							Value: func() string {
								v, _ := encrypt("secret", "key")
								return v
							}(),
							ExpiresAt: now().Add(-1 * time.Hour),
						},
					},
				},
				id:  "1",
				key: "key",
			},
			wantErr: ErrSecretNotFound,
		},
		{
			name: "get secret - error",
			input: struct {
				secrets db.SecretRepository
				id      string
				key     string
			}{
				secrets: &mockSecretRepository{
					err: dberrors.ErrSecretNotFound,
				},
				id: "1",
			},
			wantErr: ErrSecretNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svc := &service{
				secrets: test.input.secrets,
				timeout: defaultTimeout,
			}

			got, gotErr := svc.Get(test.input.id, test.input.key)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Secret{})); diff != "" {
				t.Errorf("Get() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Get() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestService_Create(t *testing.T) {
	now = func() time.Time {
		return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	n := now()

	var tests = []struct {
		name  string
		input struct {
			secrets db.SecretRepository
			secret  Secret
			id      string
		}
		want    Secret
		wantErr error
	}{
		{
			name: "create secret",
			input: struct {
				secrets db.SecretRepository
				secret  Secret
				id      string
			}{
				secrets: &mockSecretRepository{},
				secret: Secret{
					Value:      "secret",
					Passphrase: "key",
				},
				id: "2",
			},
			want: Secret{
				ID:        "2",
				TTL:       defaultTTL,
				ExpiresAt: n.Add(defaultTTL),
			},
		},
		{
			name: "create secret - no passphrase",
			input: struct {
				secrets db.SecretRepository
				secret  Secret
				id      string
			}{
				secrets: &mockSecretRepository{},
				secret: Secret{
					Value: "secret",
				},
				id: "2",
			},
			want: Secret{
				ID:        "2",
				TTL:       defaultTTL,
				ExpiresAt: n.Add(defaultTTL),
			},
		},
		{
			name: "create secret - error",
			input: struct {
				secrets db.SecretRepository
				secret  Secret
				id      string
			}{
				secrets: &mockSecretRepository{
					err: errCreateSecret,
				},
			},
			wantErr: errCreateSecret,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			newUUID = func() string {
				return test.input.id
			}

			svc := &service{
				secrets: test.input.secrets,
				timeout: defaultTimeout,
			}

			got, gotErr := svc.Create(test.input.secret)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Secret{})); diff != "" {
				t.Errorf("Create() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Create() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

type mockSecretRepository struct {
	secrets []models.Secret
	err     error
}

func (r mockSecretRepository) Generate(length int, specialCharacters bool) string {
	return ""
}

func (r mockSecretRepository) Get(ctx context.Context, id string) (models.Secret, error) {
	if r.err != nil {
		return models.Secret{}, r.err
	}

	for _, s := range r.secrets {
		if s.ID == id {
			return s, nil
		}
	}
	return models.Secret{}, nil
}

func (r *mockSecretRepository) Create(ctx context.Context, s models.Secret) (models.Secret, error) {
	if r.err != nil {
		return models.Secret{}, r.err
	}

	r.secrets = append(r.secrets, s)
	return models.Secret{
		ID:        s.ID,
		ExpiresAt: s.ExpiresAt,
	}, nil
}

func (r mockSecretRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r mockSecretRepository) DeleteExpired(ctx context.Context) error {
	return nil
}

func (r mockSecretRepository) GetSettings(ctx context.Context) (models.Settings, error) {
	return models.Settings{}, nil
}

func (r mockSecretRepository) CreateSettings(ctx context.Context, settings models.Settings) (models.Settings, error) {
	return models.Settings{}, nil
}

func (r mockSecretRepository) UpdateSettings(ctx context.Context, settings models.Settings) (models.Settings, error) {
	return models.Settings{}, nil
}

func (r mockSecretRepository) Close() error {
	return nil
}

var (
	errCreateSecret = errors.New("create secret error")
)
