package config

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestMergeConfigurations(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			dst *Configuration
			src *Configuration
		}
		want    *Configuration
		wantErr error
	}{
		{
			name: "merge configurations - empty dst and src",
			input: struct {
				dst *Configuration
				src *Configuration
			}{
				dst: &Configuration{},
				src: &Configuration{},
			},
			want: &Configuration{},
		},
		{
			name: "merge configurations - empty dst",
			input: struct {
				dst *Configuration
				src *Configuration
			}{
				dst: &Configuration{},
				src: &Configuration{
					Server: Server{
						Host: "localhost",
						Port: "8080",
					},
					Services: Services{
						Secrets: Secrets{
							EncryptionKey: "key",
							Timeout:       10 * time.Second,
						},
						Database: Database{
							URI:            "mongodb://localhost:27017",
							Address:        "localhost:27017",
							Database:       "test",
							Username:       "user",
							Password:       "password",
							Collection:     "collection",
							Timeout:        10 * time.Second,
							ConnectTimeout: 10 * time.Second,
							EnableTLS:      toPtr(true),
						},
					},
				},
			},
			want: &Configuration{
				Server: Server{
					Host: "localhost",
					Port: "8080",
				},
				Services: Services{
					Secrets: Secrets{
						EncryptionKey: "key",
						Timeout:       10 * time.Second,
					},
					Database: Database{
						URI:            "mongodb://localhost:27017",
						Address:        "localhost:27017",
						Database:       "test",
						Username:       "user",
						Password:       "password",
						Collection:     "collection",
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
				dst *Configuration
				src *Configuration
			}{
				dst: &Configuration{
					Server: Server{
						Host: "localhost",
						Port: "8080",
					},
					Services: Services{
						Secrets: Secrets{
							EncryptionKey: "key",
							Timeout:       10 * time.Second,
						},
						Database: Database{
							URI:            "mongodb://localhost:27017",
							Address:        "localhost:27017",
							Database:       "test",
							Username:       "user",
							Password:       "password",
							Collection:     "collection",
							Timeout:        10 * time.Second,
							ConnectTimeout: 10 * time.Second,
							EnableTLS:      toPtr(true),
						},
					},
				},
				src: &Configuration{
					Server: Server{
						Host: "0.0.0.0",
						Port: "8081",
					},
					Services: Services{
						Secrets: Secrets{
							EncryptionKey: "new-key",
							Timeout:       20 * time.Second,
						},
						Database: Database{
							URI:            "mongodb://0.0.0.0:27017",
							Address:        "0.0.0.0:27017",
							Database:       "new-test",
							Username:       "new-user",
							Password:       "new-password",
							Collection:     "new-collection",
							Timeout:        20 * time.Second,
							ConnectTimeout: 20 * time.Second,
							EnableTLS:      toPtr(true),
						},
					},
				},
			},
			want: &Configuration{
				Server: Server{
					Host: "0.0.0.0",
					Port: "8081",
				},
				Services: Services{
					Secrets: Secrets{
						EncryptionKey: "new-key",
						Timeout:       20 * time.Second,
					},
					Database: Database{
						URI:            "mongodb://0.0.0.0:27017",
						Address:        "0.0.0.0:27017",
						Database:       "new-test",
						Username:       "new-user",
						Password:       "new-password",
						Collection:     "new-collection",
						Timeout:        20 * time.Second,
						ConnectTimeout: 20 * time.Second,
						EnableTLS:      toPtr(true),
					},
				},
			},
		},
		{
			name: "merge configurations - empty src",
			input: struct {
				dst *Configuration
				src *Configuration
			}{
				dst: &Configuration{
					Server: Server{
						Host: "localhost",
						Port: "8080",
					},
					Services: Services{
						Secrets: Secrets{
							EncryptionKey: "key",
							Timeout:       10 * time.Second,
						},
						Database: Database{
							URI:            "mongodb://localhost:27017",
							Address:        "localhost:27017",
							Database:       "test",
							Username:       "user",
							Password:       "password",
							Collection:     "collection",
							Timeout:        10 * time.Second,
							ConnectTimeout: 10 * time.Second,
							EnableTLS:      toPtr(true),
						},
					},
				},
				src: &Configuration{},
			},
			want: &Configuration{
				Server: Server{
					Host: "localhost",
					Port: "8080",
				},
				Services: Services{
					Secrets: Secrets{
						EncryptionKey: "key",
						Timeout:       10 * time.Second,
					},
					Database: Database{
						URI:            "mongodb://localhost:27017",
						Address:        "localhost:27017",
						Database:       "test",
						Username:       "user",
						Password:       "password",
						Collection:     "collection",
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

			gotErr := mergeConfigurations(test.input.dst, test.input.src)

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

// toPtr returns a pointer to the value v.
func toPtr[T any](v T) *T {
	return &v
}
