package config

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseFlags(t *testing.T) {
	var tests = []struct {
		name    string
		input   []string
		want    flags
		wantErr error
	}{
		{
			name:  "parse flags - empty",
			input: []string{},
			want:  flags{},
		},
		{
			name: "parse flags - all flags",
			input: []string{
				"-config-path", "path",
				"-host", "host",
				"-port", "3001",
				"-encryption-key", "key",
				"-timeout", "15s",
				"-database-uri", "uri",
				"-database-address", "address",
				"-database", "database",
				"-database-user", "user",
				"-database-password", "password",
				"-database-timeout", "15s",
				"-database-connect-timeout", "15s",
				"-database-enable-tls", "false",
			},
			want: flags{
				configPath:             "path",
				host:                   "host",
				port:                   "3001",
				encryptionKey:          "key",
				timeout:                time.Second * 15,
				databaseURI:            "uri",
				databaseAddr:           "address",
				database:               "database",
				databaseUser:           "user",
				databasePass:           "password",
				databaseTimeout:        time.Second * 15,
				databaseConnectTimeout: time.Second * 15,
				databaseEnableTLS:      toPtr(false),
			},
		},
		{
			name: "parse flags - set enable-tls to true",
			input: []string{
				"-database-enable-tls", "true",
			},
			want: flags{
				databaseEnableTLS: toPtr(true),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _, gotErr := parseFlags(test.input)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(flags{})); diff != "" {
				t.Errorf("parseFlags() = unexpected result (-want +got)\n%s\n", diff)
			}

			if test.wantErr == nil && gotErr != nil {
				t.Errorf("parseFlags() = unexpected error: %v\n", gotErr)
			}
		})
	}
}
