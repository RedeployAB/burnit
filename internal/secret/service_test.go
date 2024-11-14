package secret

import (
	"context"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/RedeployAB/burnit/internal/db"
	dberrors "github.com/RedeployAB/burnit/internal/db/errors"
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
				secrets:         &mockSecretRepository{},
				timeout:         defaultTimeout,
				cleanupInterval: defaultCleanupInterval,
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
						s.cleanupInterval = 30 * time.Second
					},
				},
			},
			want: &service{
				secrets:         &mockSecretRepository{},
				timeout:         30 * time.Second,
				cleanupInterval: 30 * time.Second,
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

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(service{}, mockSecretRepository{}), cmpopts.IgnoreFields(service{}, "stopCh")); diff != "" {
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
					secrets: []db.Secret{
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
					secrets: []db.Secret{
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
					err: errGetSecret,
				},
				id: "1",
			},
			wantErr: errGetSecret,
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
				ID:         "2",
				Passphrase: "key",
				TTL:        time.Until(n.Add(defaultTTL)).Round(time.Minute),
				ExpiresAt:  n.Add(defaultTTL),
			},
		},
		{
			name: "create secret - base64 encoded value",
			input: struct {
				secrets db.SecretRepository
				secret  Secret
				id      string
			}{
				secrets: &mockSecretRepository{},
				secret: Secret{
					Value: func() string {
						return base64.StdEncoding.EncodeToString([]byte("secret"))
					}(),
					Passphrase: "key",
				},
				id: "2",
			},
			want: Secret{
				ID:         "2",
				Passphrase: "key",
				TTL:        time.Until(n.Add(defaultTTL)).Round(time.Minute),
				ExpiresAt:  n.Add(defaultTTL),
			},
		},
		{
			name: "create secret - base64 encoded value with invalid data",
			input: struct {
				secrets db.SecretRepository
				secret  Secret
				id      string
			}{
				secrets: &mockSecretRepository{},
				secret: Secret{
					Value: func() string {
						f, err := os.OpenFile("../../assets/burnit.png", os.O_RDONLY, 0644)
						if err != nil {
							panic(err)
						}
						b, _ := io.ReadAll(f)
						return base64.StdEncoding.EncodeToString(b)
					}(),
					Passphrase: "key",
				},
				id: "2",
			},
			wantErr: ErrSecretTooManyBytes,
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
				secret: Secret{
					Value:      "secret",
					Passphrase: "key",
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

func TestService_Delete(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			secrets db.SecretRepository
			id      string
		}
		wantErr error
	}{
		{
			name: "delete secret",
			input: struct {
				secrets db.SecretRepository
				id      string
			}{
				secrets: &mockSecretRepository{
					secrets: []db.Secret{
						{
							ID:    "1",
							Value: "secret",
						},
					},
				},
				id: "1",
			},
			wantErr: nil,
		},
		{
			name: "delete secret - not found",
			input: struct {
				secrets db.SecretRepository
				id      string
			}{
				secrets: &mockSecretRepository{},
				id:      "1",
			},
			wantErr: ErrSecretNotFound,
		},
		{
			name: "delete secret - error",
			input: struct {
				secrets db.SecretRepository
				id      string
			}{
				secrets: &mockSecretRepository{
					secrets: []db.Secret{
						{
							ID:    "1",
							Value: "secret",
						},
					},
					err: errDeleteSecret,
				},
				id: "1",
			},
			wantErr: errDeleteSecret,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svc := &service{
				secrets: test.input.secrets,
				timeout: defaultTimeout,
			}

			gotErr := svc.Delete(test.input.id)

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Delete() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestService_DeleteExpired(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			secrets db.SecretRepository
		}
		wantErr error
	}{
		{
			name: "delete expired secrets",
			input: struct {
				secrets db.SecretRepository
			}{
				secrets: &mockSecretRepository{
					secrets: []db.Secret{
						{
							ID:        "1",
							Value:     "secret",
							ExpiresAt: pastHour,
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "delete expired secrets - no secrets",
			input: struct {
				secrets db.SecretRepository
			}{
				secrets: &mockSecretRepository{
					secrets: []db.Secret{},
				},
			},
			wantErr: nil,
		},
		{
			name: "delete expired secrets - error",
			input: struct {
				secrets db.SecretRepository
			}{
				secrets: &mockSecretRepository{
					secrets: []db.Secret{
						{
							ID:        "1",
							Value:     "secret",
							ExpiresAt: pastHour,
						},
					},
					err: errDeleteManySecrets,
				},
			},
			wantErr: errDeleteManySecrets,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svc := &service{
				secrets: test.input.secrets,
				timeout: defaultTimeout,
			}

			gotErr := svc.DeleteExpired()

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("DeleteExpired() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestValidValue(t *testing.T) {
	var tests = []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid value",
			input:   "value",
			wantErr: nil,
		},
		{
			name: "valid value - base64 encoded",
			input: func() string {
				return base64.StdEncoding.EncodeToString([]byte("value"))
			}(),
			wantErr: nil,
		},
		{
			name: "invalid value - base64 encoded",
			input: func() string {
				return base64.StdEncoding.EncodeToString([]byte{0})
			}(),
			wantErr: ErrInvalidSecretValue,
		},
		{
			name: "valid value - base32 encoded",
			input: func() string {
				return base32.StdEncoding.EncodeToString([]byte("value"))
			}(),
			wantErr: nil,
		},
		{
			name: "invalid value - base32 encoded",
			input: func() string {
				return base32.StdEncoding.EncodeToString([]byte{0, 1, 2, 3})
			}(),
			wantErr: ErrInvalidSecretValue,
		},
		{
			name: "valid value - base32 encoded (hex)",
			input: func() string {
				return base32.HexEncoding.EncodeToString([]byte("value"))
			}(),
			wantErr: nil,
		},
		{
			name: "invalid value - base32 encoded (hex)",
			input: func() string {
				return base32.HexEncoding.EncodeToString([]byte{0, 1, 2, 3})
			}(),
			wantErr: ErrInvalidSecretValue,
		},
		{
			name: "valid value - hex encoded",
			input: func() string {
				return hex.EncodeToString([]byte("value"))
			}(),
			wantErr: nil,
		},
		{
			name: "invalid value - hex encoded",
			input: func() string {
				return hex.EncodeToString([]byte{0, 1, 2, 3})
			}(),
			wantErr: ErrInvalidSecretValue,
		},
		{
			name: "invalid length",
			input: func() string {
				return base64.StdEncoding.EncodeToString([]byte(strings.Repeat("value", 100)))
			}(),
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := validValue(test.input)

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("validValue() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

type mockSecretRepository struct {
	secrets []db.Secret
	err     error
}

func (r mockSecretRepository) Get(ctx context.Context, id string) (db.Secret, error) {
	if r.err != nil {
		return db.Secret{}, r.err
	}

	for _, s := range r.secrets {
		if s.ID == id {
			return s, nil
		}
	}
	return db.Secret{}, dberrors.ErrSecretNotFound
}

func (r *mockSecretRepository) Create(ctx context.Context, s db.Secret) (db.Secret, error) {
	if r.err != nil {
		return db.Secret{}, r.err
	}

	r.secrets = append(r.secrets, s)
	return db.Secret{
		ID:        s.ID,
		ExpiresAt: s.ExpiresAt,
	}, nil
}

func (r *mockSecretRepository) Delete(ctx context.Context, id string) error {
	if r.err != nil && errors.Is(r.err, errDeleteSecret) {
		return r.err
	}

	for i, s := range r.secrets {
		if s.ID == id {
			r.secrets = append(r.secrets[:i], r.secrets[i+1:]...)
			return nil
		}
	}

	return dberrors.ErrSecretNotFound
}

func (r *mockSecretRepository) DeleteExpired(ctx context.Context) error {
	if r.err != nil && errors.Is(r.err, errDeleteManySecrets) {
		return r.err
	}

	for i, s := range r.secrets {
		if s.ExpiresAt.Before(time.Now().UTC()) {
			r.secrets = append(r.secrets[:i], r.secrets[i+1:]...)
		}
	}
	return dberrors.ErrSecretsNotDeleted
}

func (r mockSecretRepository) Close() error {
	return nil
}

var (
	errGetSecret         = errors.New("get secret error")
	errCreateSecret      = errors.New("create secret error")
	errDeleteSecret      = errors.New("delete secret error")
	errDeleteManySecrets = errors.New("delete many secrets error")
)

var (
	pastHour = time.Now().UTC().Add(-1 * time.Hour)
)
