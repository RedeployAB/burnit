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
					TLSMode:  TLSModeRequire,
				},
			},
			want: "postgres://user:password@localhost/database?sslmode=require",
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
			want: "sqlserver://localhost?database=burnit",
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
					TLSMode:  TLSModeTrue,
				},
			},
			want: "sqlserver://user:password@localhost?database=database&encrypt=true",
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
