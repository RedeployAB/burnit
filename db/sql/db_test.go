package sql

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuildDSN(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			driver  Driver
			options *Options
		}
		want string
	}{
		{
			name: "postgres",
			input: struct {
				driver  Driver
				options *Options
			}{
				driver: DriverPostgres,
				options: &Options{
					Address: "localhost",
				},
			},
			want: "postgres://localhost/burnit",
		},
		{
			name: "postgres - database, username, password and TLS mode",
			input: struct {
				driver  Driver
				options *Options
			}{
				driver: DriverPostgres,
				options: &Options{
					Address:  "localhost",
					Database: "database",
					Username: "user",
					Password: "password",
					Postgres: PostgresOptions{
						SSLMode: PostgresSSLModeRequire,
					},
				},
			},
			want: "postgres://user:password@localhost/database?sslmode=require",
		},
		{
			name: "postgres with DSN/URI",
			input: struct {
				driver  Driver
				options *Options
			}{
				driver: DriverPostgres,
				options: &Options{
					Driver: DriverPostgres,
					DSN:    "postgres://user:password@localhost/database?sslmode=prefer",
				},
			},
			want: "postgres://user:password@localhost/database?sslmode=prefer",
		},
		{
			name: "mssql",
			input: struct {
				driver  Driver
				options *Options
			}{
				driver: DriverMSSQL,
				options: &Options{
					Address: "localhost",
				},
			},
			want: "sqlserver://localhost?database=Burnit",
		},
		{
			name: "mssql - database, username, password and TLS mode",
			input: struct {
				driver  Driver
				options *Options
			}{
				driver: DriverMSSQL,
				options: &Options{
					Address:  "localhost",
					Database: "database",
					Username: "user",
					Password: "password",
					MSSQL: MSSQLOptions{
						Encrypt: MSSQLEncryptTrue,
					},
				},
			},
			want: "sqlserver://user:password@localhost?database=database&encrypt=true",
		},
		{
			name: "mssql with DSN/URI",
			input: struct {
				driver  Driver
				options *Options
			}{
				driver: DriverMSSQL,
				options: &Options{
					Driver: DriverMSSQL,
					DSN:    "sqlserver://user:password@localhost?database=database&encrypt=true",
				},
			},
			want: "sqlserver://user:password@localhost?database=database&encrypt=true",
		},
		{
			name: "sqlite - file",
			input: struct {
				driver  Driver
				options *Options
			}{
				driver: DriverSQLite,
				options: &Options{
					SQLite: SQLiteOptions{
						File: "file.db",
					},
				},
			},
			want: "file:file.db",
		},
		{
			name: "sqlite - default file",
			input: struct {
				driver  Driver
				options *Options
			}{
				driver:  DriverSQLite,
				options: &Options{},
			},
			want: "file:burnit.db",
		},
		{
			name: "sqlite - in-memory",
			input: struct {
				driver  Driver
				options *Options
			}{
				driver: DriverSQLite,
				options: &Options{
					SQLite: SQLiteOptions{
						InMemory: true,
					},
				},
			},
			want: ":memory:",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := buildDSN(test.input.driver, test.input.options)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("buildDSN() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}
