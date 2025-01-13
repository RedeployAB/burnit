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
				"-tls-cert-file", "cert",
				"-tls-key-file", "key",
				"-rate-limiter-rate", "10",
				"-rate-limiter-burst", "10",
				"-rate-limiter-cleanup-interval", "15s",
				"-cors-origin", "origin",
				"-secret-service-timeout", "15s",
				"-database-driver", "postgres",
				"-database-uri", "uri",
				"-database-address", "address",
				"-database", "database",
				"-database-user", "user",
				"-database-password", "password",
				"-database-timeout", "15s",
				"-database-connect-timeout", "15s",
				"-database-mongo-enable-tls", "true",
				"-database-postgres-ssl-mode", "enable",
				"-database-mssql-encrypt", "true",
				"-database-sqlite-file", "file.db",
				"-database-sqlite-in-memory", "true",
				"-database-redis-dial-timeout", "15s",
				"-database-redis-max-retries", "10",
				"-database-redis-min-retry-backoff", "15s",
				"-database-redis-max-retry-backoff", "15s",
				"-database-redis-enable-tls", "true",
				"-session-service-timeout", "15s",
				"-runtime-render", "true",
				"-session-database-driver", "postgres",
				"-session-database-uri", "uri",
				"-session-database-address", "address",
				"-session-database", "database",
				"-session-database-user", "user",
				"-session-database-password", "password",
				"-session-database-timeout", "15s",
				"-session-database-connect-timeout", "15s",
				"-session-database-mongo-enable-tls", "true",
				"-session-database-postgres-ssl-mode", "enable",
				"-session-database-mssql-encrypt", "true",
				"-session-database-sqlite-file", "file.db",
				"-session-database-sqlite-in-memory", "true",
				"-session-database-redis-dial-timeout", "15s",
				"-session-database-redis-max-retries", "10",
				"-session-database-redis-min-retry-backoff", "15s",
				"-session-database-redis-max-retry-backoff", "15s",
				"-session-database-redis-enable-tls", "true",
				"-local-development", "true",
			},
			want: flags{
				configPath:                          "path",
				host:                                "host",
				port:                                3001,
				tlsCertFile:                         "cert",
				tlsKeyFile:                          "key",
				corsOrigin:                          "origin",
				rateLimiterRate:                     10,
				rateLimiterBurst:                    10,
				rateLimiterCleanupInterval:          time.Second * 15,
				secretServiceTimeout:                time.Second * 15,
				databaseDriver:                      "postgres",
				databaseURI:                         "uri",
				databaseAddr:                        "address",
				database:                            "database",
				databaseUser:                        "user",
				databasePass:                        "password",
				databaseTimeout:                     time.Second * 15,
				databaseConnectTimeout:              time.Second * 15,
				databaseMongoEnableTLS:              toPtr(true),
				databasePostgresSSLMode:             "enable",
				databaseMSSQLEncrypt:                "true",
				databaseSQLiteFile:                  "file.db",
				databaseSQLiteInMemory:              toPtr(true),
				databaseRedisDialTimeout:            time.Second * 15,
				databaseRedisMaxRetries:             10,
				databaseRedisMinRetryBackoff:        time.Second * 15,
				databaseRedisMaxRetryBackoff:        time.Second * 15,
				databaseRedisEnableTLS:              toPtr(true),
				sessionServiceTimeout:               time.Second * 15,
				runtimeRender:                       toPtr(true),
				sessionDatabaseDriver:               "postgres",
				sessionDatabaseURI:                  "uri",
				sessionDatabaseAddr:                 "address",
				sessionDatabase:                     "database",
				sessionDatabaseUser:                 "user",
				sessionDatabasePass:                 "password",
				sessionDatabaseTimeout:              time.Second * 15,
				sessionDatabaseConnectTimeout:       time.Second * 15,
				sessionDatabaseMongoEnableTLS:       toPtr(true),
				sessionDatabasePostgresSSLMode:      "enable",
				sessionDatabaseMSSQLEncrypt:         "true",
				sessionDatabaseSQLiteFile:           "file.db",
				sessionDatabaseSQLiteInMemory:       toPtr(true),
				sessionDatabaseRedisDialTimeout:     time.Second * 15,
				sessionDatabaseRedisMaxRetries:      10,
				sessionDatabaseRedisMinRetryBackoff: time.Second * 15,
				sessionDatabaseRedisMaxRetryBackoff: time.Second * 15,
				sessionDatabaseRedisEnableTLS:       toPtr(true),
				localDevelopment:                    toPtr(true),
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
