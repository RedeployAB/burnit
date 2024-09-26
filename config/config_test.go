package config

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			envs map[string]string
			args []string
		}
		want    *Configuration
		wantErr error
	}{
		{
			name: "new configuration - default",
			input: struct {
				envs map[string]string
				args []string
			}{},
			want: &Configuration{
				Server: Server{
					Host: "0.0.0.0",
					Port: 3000,
				},
				Services: Services{
					Secret: Secret{
						Timeout: defaultSecretServiceTimeout,
					},
					Database: Database{
						Database:       defaultDatabaseName,
						Timeout:        defaultDatabaseTimeout,
						ConnectTimeout: defaultDatabaseConnectTimeout,
						EnableTLS:      toPtr(true),
					},
				},
			},
		},
		{
			name: "new configuration - config file",
			input: struct {
				envs map[string]string
				args []string
			}{
				args: []string{"-config-path", "../testdata/config.yaml"},
			},
			want: &Configuration{
				Server: Server{
					Host: "localhost",
					Port: 3001,
					TLS: TLS{
						CertFile: "cert.pem",
						KeyFile:  "key.pem",
					},
				},
				Services: Services{
					Secret: Secret{
						EncryptionKey: "key",
						Timeout:       15 * time.Second,
					},
					Database: Database{
						URI:            "mongodb://localhost:27017",
						Address:        "localhost:27017",
						Database:       "test",
						Username:       "test",
						Password:       "test",
						Timeout:        15 * time.Second,
						ConnectTimeout: 15 * time.Second,
						EnableTLS:      toPtr(false),
					},
				},
			},
		},
		{
			name: "new configuration - config file, environment variables (override config file)",
			input: struct {
				envs map[string]string
				args []string
			}{
				args: []string{"-config-path", "../testdata/config.yaml"},
				envs: map[string]string{
					"BURNIT_LISTEN_HOST":              "localhost2",
					"BURNIT_LISTEN_PORT":              "3002",
					"BURNIT_TLS_CERT_FILE":            "cert2.pem",
					"BURNIT_TLS_KEY_FILE":             "key2.pem",
					"BURNIT_SECRETS_ENCRYPTION_KEY":   "key2",
					"BURNIT_SECRETS_TIMEOUT":          "20s",
					"BURNIT_DATABASE_URI":             "mongodb://localhost2:27018",
					"BURNIT_DATABASE_ADDRESS":         "localhost2:27018",
					"BURNIT_DATABASE":                 "test2",
					"BURNIT_DATABASE_USERNAME":        "test2",
					"BURNIT_DATABASE_PASSWORD":        "test2",
					"BURNIT_DATABASE_TIMEOUT":         "20s",
					"BURNIT_DATABASE_CONNECT_TIMEOUT": "20s",
					"BURNIT_DATABASE_ENABLE_TLS":      "true",
				},
			},
			want: &Configuration{
				Server: Server{
					Host: "localhost2",
					Port: 3002,
					TLS: TLS{
						CertFile: "cert2.pem",
						KeyFile:  "key2.pem",
					},
				},
				Services: Services{
					Secret: Secret{
						EncryptionKey: "key2",
						Timeout:       20 * time.Second,
					},
					Database: Database{
						URI:            "mongodb://localhost2:27018",
						Address:        "localhost2:27018",
						Database:       "test2",
						Username:       "test2",
						Password:       "test2",
						Timeout:        20 * time.Second,
						ConnectTimeout: 20 * time.Second,
						EnableTLS:      toPtr(true),
					},
				},
			},
		},
		{
			name: "new configuration - config file, environment variables, flags (override config file and environment variables)",
			input: struct {
				envs map[string]string
				args []string
			}{
				args: []string{
					"-config-path", "../testdata/config.yaml",
					"-host", "localhost3",
					"-port", "3003",
					"-tls-cert-file", "cert3.pem",
					"-tls-key-file", "key3.pem",
					"-encryption-key", "key3",
					"-timeout", "25s",
					"-database-uri", "mongodb://localhost3:27019",
					"-database-address", "localhost3:27019",
					"-database", "test3",
					"-database-user", "test3",
					"-database-password", "test3",
					"-database-timeout", "25s",
					"-database-connect-timeout", "25s",
					"-database-enable-tls", "false",
				},
				envs: map[string]string{
					"BURNIT_LISTEN_HOST":              "localhost2",
					"BURNIT_LISTEN_PORT":              "3002",
					"BURNIT_TLS_CERT_FILE":            "cert2.pem",
					"BURNIT_TLS_KEY_FILE":             "key2.pem",
					"BURNIT_SECRETS_ENCRYPTION_KEY":   "key2",
					"BURNIT_SECRETS_TIMEOUT":          "20s",
					"BURNIT_DATABASE_URI":             "mongodb://localhost2:27018",
					"BURNIT_DATABASE_ADDRESS":         "localhost2:27018",
					"BURNIT_DATABASE":                 "test2",
					"BURNIT_DATABASE_USERNAME":        "test2",
					"BURNIT_DATABASE_PASSWORD":        "test2",
					"BURNIT_DATABASE_TIMEOUT":         "20s",
					"BURNIT_DATABASE_CONNECT_TIMEOUT": "20s",
					"BURNIT_DATABASE_ENABLE_TLS":      "true",
				},
			},
			want: &Configuration{
				Server: Server{
					Host: "localhost3",
					Port: 3003,
					TLS: TLS{
						CertFile: "cert3.pem",
						KeyFile:  "key3.pem",
					},
				},
				Services: Services{
					Secret: Secret{
						EncryptionKey: "key3",
						Timeout:       25 * time.Second,
					},
					Database: Database{
						URI:            "mongodb://localhost3:27019",
						Address:        "localhost3:27019",
						Database:       "test3",
						Username:       "test3",
						Password:       "test3",
						Timeout:        25 * time.Second,
						ConnectTimeout: 25 * time.Second,
						EnableTLS:      toPtr(false),
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range test.input.envs {
				t.Setenv(k, v)
			}
			os.Args = append([]string{"cmd"}, test.input.args...)

			got, gotErr := New()

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("New() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr); diff != "" {
				t.Errorf("New() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestMergeConfigurations(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			dst  *Configuration
			srcs []*Configuration
		}
		want    *Configuration
		wantErr error
	}{
		{
			name: "merge configurations - empty dst and src",
			input: struct {
				dst  *Configuration
				srcs []*Configuration
			}{
				dst:  &Configuration{},
				srcs: []*Configuration{},
			},
			want: &Configuration{},
		},
		{
			name: "merge configurations - empty dst",
			input: struct {
				dst  *Configuration
				srcs []*Configuration
			}{
				dst: &Configuration{},
				srcs: []*Configuration{
					{
						Server: Server{
							Host: "localhost",
							Port: 8080,
						},
						Services: Services{
							Secret: Secret{
								EncryptionKey: "key",
								Timeout:       10 * time.Second,
							},
							Database: Database{
								URI:            "mongodb://localhost:27017",
								Address:        "localhost:27017",
								Database:       "test",
								Username:       "user",
								Password:       "password",
								Timeout:        10 * time.Second,
								ConnectTimeout: 10 * time.Second,
								EnableTLS:      toPtr(true),
							},
						},
					},
				},
			},
			want: &Configuration{
				Server: Server{
					Host: "localhost",
					Port: 8080,
				},
				Services: Services{
					Secret: Secret{
						EncryptionKey: "key",
						Timeout:       10 * time.Second,
					},
					Database: Database{
						URI:            "mongodb://localhost:27017",
						Address:        "localhost:27017",
						Database:       "test",
						Username:       "user",
						Password:       "password",
						Timeout:        10 * time.Second,
						ConnectTimeout: 10 * time.Second,
						EnableTLS:      toPtr(true),
					},
				},
			},
		},
		{
			name: "merge configurations - replace dst with src",
			input: struct {
				dst  *Configuration
				srcs []*Configuration
			}{
				dst: &Configuration{
					Server: Server{
						Host: "localhost",
						Port: 8080,
					},
					Services: Services{
						Secret: Secret{
							EncryptionKey: "key",
							Timeout:       10 * time.Second,
						},
						Database: Database{
							URI:            "mongodb://localhost:27017",
							Address:        "localhost:27017",
							Database:       "test",
							Username:       "user",
							Password:       "password",
							Timeout:        10 * time.Second,
							ConnectTimeout: 10 * time.Second,
							EnableTLS:      toPtr(true),
						},
					},
				},
				srcs: []*Configuration{
					{
						Server: Server{
							Host: "0.0.0.0",
							Port: 8081,
						},
						Services: Services{
							Secret: Secret{
								EncryptionKey: "new-key",
								Timeout:       20 * time.Second,
							},
							Database: Database{
								URI:            "mongodb://0.0.0.0:27017",
								Address:        "0.0.0.0:27017",
								Database:       "new-test",
								Username:       "new-user",
								Password:       "new-password",
								Timeout:        20 * time.Second,
								ConnectTimeout: 20 * time.Second,
								EnableTLS:      toPtr(true),
							},
						},
					},
				},
			},
			want: &Configuration{
				Server: Server{
					Host: "0.0.0.0",
					Port: 8081,
				},
				Services: Services{
					Secret: Secret{
						EncryptionKey: "new-key",
						Timeout:       20 * time.Second,
					},
					Database: Database{
						URI:            "mongodb://0.0.0.0:27017",
						Address:        "0.0.0.0:27017",
						Database:       "new-test",
						Username:       "new-user",
						Password:       "new-password",
						Timeout:        20 * time.Second,
						ConnectTimeout: 20 * time.Second,
						EnableTLS:      toPtr(true),
					},
				},
			},
		},
		{
			name: "merge configurations - replace dst with src in multiple srcs",
			input: struct {
				dst  *Configuration
				srcs []*Configuration
			}{
				dst: &Configuration{
					Server: Server{
						Host: "localhost",
						Port: 8080,
					},
					Services: Services{
						Secret: Secret{
							EncryptionKey: "key",
							Timeout:       10 * time.Second,
						},
						Database: Database{
							URI:            "mongodb://localhost:27017",
							Address:        "localhost:27017",
							Database:       "test",
							Username:       "user",
							Password:       "password",
							Timeout:        10 * time.Second,
							ConnectTimeout: 10 * time.Second,
							EnableTLS:      toPtr(true),
						},
					},
				},
				srcs: []*Configuration{
					{
						Server: Server{
							Host: "0.0.0.0",
							Port: 8081,
						},
						Services: Services{
							Secret: Secret{
								EncryptionKey: "new-key",
								Timeout:       20 * time.Second,
							},
							Database: Database{
								URI:            "mongodb://0.0.0.0:27017",
								Address:        "0.0.0.0:27017",
								Database:       "new-test",
								Username:       "new-user",
								Password:       "new-password",
								Timeout:        20 * time.Second,
								ConnectTimeout: 20 * time.Second,
								EnableTLS:      toPtr(true),
							},
						},
					},
					{
						Server: Server{
							Host: "0.0.0.0",
							Port: 8082,
						},
						Services: Services{
							Secret: Secret{
								EncryptionKey: "new-key-2",
								Timeout:       30 * time.Second,
							},
							Database: Database{
								URI:            "mongodb://0.0.0.0:27017",
								Address:        "0.0.0.0:27017",
								Database:       "new-test-2",
								Username:       "new-user-2",
								Password:       "new-password-2",
								Timeout:        20 * time.Second,
								ConnectTimeout: 20 * time.Second,
								EnableTLS:      toPtr(false),
							},
						},
					},
				},
			},
			want: &Configuration{
				Server: Server{
					Host: "0.0.0.0",
					Port: 8082,
				},
				Services: Services{
					Secret: Secret{
						EncryptionKey: "new-key-2",
						Timeout:       30 * time.Second,
					},
					Database: Database{
						URI:            "mongodb://0.0.0.0:27017",
						Address:        "0.0.0.0:27017",
						Database:       "new-test-2",
						Username:       "new-user-2",
						Password:       "new-password-2",
						Timeout:        20 * time.Second,
						ConnectTimeout: 20 * time.Second,
						EnableTLS:      toPtr(false),
					},
				},
			},
		},
		{
			name: "merge configurations - empty src",
			input: struct {
				dst  *Configuration
				srcs []*Configuration
			}{
				dst: &Configuration{
					Server: Server{
						Host: "localhost",
						Port: 8080,
					},
					Services: Services{
						Secret: Secret{
							EncryptionKey: "key",
							Timeout:       10 * time.Second,
						},
						Database: Database{
							URI:            "mongodb://localhost:27017",
							Address:        "localhost:27017",
							Database:       "test",
							Username:       "user",
							Password:       "password",
							Timeout:        10 * time.Second,
							ConnectTimeout: 10 * time.Second,
							EnableTLS:      toPtr(true),
						},
					},
				},
				srcs: []*Configuration{},
			},
			want: &Configuration{
				Server: Server{
					Host: "localhost",
					Port: 8080,
				},
				Services: Services{
					Secret: Secret{
						EncryptionKey: "key",
						Timeout:       10 * time.Second,
					},
					Database: Database{
						URI:            "mongodb://localhost:27017",
						Address:        "localhost:27017",
						Database:       "test",
						Username:       "user",
						Password:       "password",
						Timeout:        10 * time.Second,
						ConnectTimeout: 10 * time.Second,
						EnableTLS:      toPtr(true),
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got *Configuration

			gotErr := mergeConfigurations(test.input.dst, test.input.srcs...)

			if gotErr != nil && test.wantErr == nil {
				t.Errorf("mergeConfigurations() = unexpected error: %v\n", gotErr)
			}

			got = test.input.dst

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("mergeConfigurations() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}
