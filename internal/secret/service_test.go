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
			secrets db.SecretStore
			options []ServiceOption
		}
		want    *service
		wantErr error
	}{
		{
			name: "new service",
			input: struct {
				secrets db.SecretStore
				options []ServiceOption
			}{
				secrets: &stubSecretStore{},
			},
			want: &service{
				secrets:                 &stubSecretStore{},
				timeout:                 defaultTimeout,
				cleanupInterval:         defaultCleanupInterval,
				valueMaxCharacters:      defaultValueMaxCharacters,
				passphraseMinCharacters: defaultPassphraseMinCharacters,
				passphraseMaxCharacters: defaultPassphraseMaxCharacters,
			},
		},
		{
			name: "new service - with options",
			input: struct {
				secrets db.SecretStore
				options []ServiceOption
			}{
				secrets: &stubSecretStore{},
				options: []ServiceOption{
					func(s *service) {
						s.timeout = 30 * time.Second
						s.cleanupInterval = 30 * time.Second
						s.valueMaxCharacters = 4000
						s.passphraseMinCharacters = 3
						s.passphraseMaxCharacters = 8
					},
				},
			},
			want: &service{
				secrets:                 &stubSecretStore{},
				timeout:                 30 * time.Second,
				cleanupInterval:         30 * time.Second,
				valueMaxCharacters:      4000,
				passphraseMinCharacters: 3,
				passphraseMaxCharacters: 8,
			},
		},
		{
			name: "new service - nil store",
			input: struct {
				secrets db.SecretStore
				options []ServiceOption
			}{
				secrets: nil,
			},
			wantErr: ErrNilStore,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewService(test.input.secrets, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(service{}, stubSecretStore{}), cmpopts.IgnoreFields(service{}, "stopCh")); diff != "" {
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
			secrets db.SecretStore
			id      string
			key     string
		}
		want    Secret
		wantErr error
	}{
		{
			name: "get secret",
			input: struct {
				secrets db.SecretStore
				id      string
				key     string
			}{
				secrets: &stubSecretStore{
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
				secrets db.SecretStore
				id      string
				key     string
			}{
				secrets: &stubSecretStore{},
				id:      "1",
			},
			wantErr: ErrSecretNotFound,
		},
		{
			name: "get secret - expired",
			input: struct {
				secrets db.SecretStore
				id      string
				key     string
			}{
				secrets: &stubSecretStore{
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
				secrets db.SecretStore
				id      string
				key     string
			}{
				secrets: &stubSecretStore{
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
			secrets db.SecretStore
			secret  Secret
			id      string
		}
		want    Secret
		wantErr error
	}{
		{
			name: "create secret",
			input: struct {
				secrets db.SecretStore
				secret  Secret
				id      string
			}{
				secrets: &stubSecretStore{},
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
				secrets db.SecretStore
				secret  Secret
				id      string
			}{
				secrets: &stubSecretStore{},
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
				secrets db.SecretStore
				secret  Secret
				id      string
			}{
				secrets: &stubSecretStore{},
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
			wantErr: ErrValueInvalid,
		},
		{
			name: "create secret - error",
			input: struct {
				secrets db.SecretStore
				secret  Secret
				id      string
			}{
				secrets: &stubSecretStore{
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
				secrets:                 test.input.secrets,
				valueMaxCharacters:      40000,
				passphraseMinCharacters: 3,
				passphraseMaxCharacters: 8,
				timeout:                 defaultTimeout,
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
			secrets db.SecretStore
			id      string
		}
		wantErr error
	}{
		{
			name: "delete secret",
			input: struct {
				secrets db.SecretStore
				id      string
			}{
				secrets: &stubSecretStore{
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
				secrets db.SecretStore
				id      string
			}{
				secrets: &stubSecretStore{},
				id:      "1",
			},
			wantErr: ErrSecretNotFound,
		},
		{
			name: "delete secret - error",
			input: struct {
				secrets db.SecretStore
				id      string
			}{
				secrets: &stubSecretStore{
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

func TestValidValue(t *testing.T) {
	var tests = []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid value - input 'value'",
			input:   "value",
			wantErr: nil,
		},
		{
			name:    "valid value - input 'secret'",
			input:   "secret",
			wantErr: nil,
		},
		{
			name:    "valid value - input base64 like",
			input:   "testinglocalstringtestinglocalstringtestinglocalstringtestinglocalstringtestinglocalstringtestinglocalstringtestinglocalstring",
			wantErr: nil,
		},
		{
			name: "valid value - base64 encoded",
			input: func() string {
				return base64.StdEncoding.EncodeToString([]byte("secret"))
			}(),
			wantErr: nil,
		},
		{
			name: "invalid value - base64 encoded",
			input: func() string {
				f, err := os.OpenFile("../../assets/burnit.png", os.O_RDONLY, 0644)
				if err != nil {
					panic(err)
				}
				b, _ := io.ReadAll(f)
				return base64.StdEncoding.EncodeToString(b)
			}(),
			wantErr: ErrValueInvalid,
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
				f, err := os.OpenFile("../../assets/burnit.png", os.O_RDONLY, 0644)
				if err != nil {
					panic(err)
				}
				b, _ := io.ReadAll(f)
				return base32.StdEncoding.EncodeToString(b)
			}(),
			wantErr: ErrValueInvalid,
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
				f, err := os.OpenFile("../../assets/burnit.png", os.O_RDONLY, 0644)
				if err != nil {
					panic(err)
				}
				b, _ := io.ReadAll(f)
				return base32.HexEncoding.EncodeToString(b)
			}(),
			wantErr: ErrValueInvalid,
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
				f, err := os.OpenFile("../../assets/burnit.png", os.O_RDONLY, 0644)
				if err != nil {
					panic(err)
				}
				b, _ := io.ReadAll(f)
				return hex.EncodeToString(b)
			}(),
			wantErr: ErrValueInvalid,
		},
		{
			name: "invalid amount of characters",
			input: func() string {
				return base64.StdEncoding.EncodeToString([]byte(strings.Repeat("value", 40000+1)))
			}(),
			wantErr: ErrValueTooManyCharacters,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := validValue(test.input, 40000)

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("validValue() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestValidPassphrase(t *testing.T) {
	var tests = []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid passphrase",
			input:   "test",
			wantErr: nil,
		},
		{
			name:    "invalid passphrase - too few characters",
			input:   "te",
			wantErr: ErrPassphraseTooFewCharacters,
		},
		{
			name:    "invalid passphrase - too many characters",
			input:   "testtesttest",
			wantErr: ErrPassphraseTooManyCharacters,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := validPassphrase(test.input, 4, 8)

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("validPassphrase() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

type stubSecretStore struct {
	secrets []db.Secret
	err     error
}

func (r stubSecretStore) Get(ctx context.Context, id string) (db.Secret, error) {
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

func (r *stubSecretStore) Create(ctx context.Context, s db.Secret) (db.Secret, error) {
	if r.err != nil {
		return db.Secret{}, r.err
	}

	r.secrets = append(r.secrets, s)
	return db.Secret{
		ID:        s.ID,
		ExpiresAt: s.ExpiresAt,
	}, nil
}

func (r *stubSecretStore) Delete(ctx context.Context, id string) error {
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

func (r *stubSecretStore) DeleteExpired(ctx context.Context) error {
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

func (r stubSecretStore) Close() error {
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
