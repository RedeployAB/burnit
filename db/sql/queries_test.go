package sql

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCreateQueries(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			driver Driver
			table  string
		}
		want    queries
		wantErr error
	}{
		{
			name: "postgres",
			input: struct {
				driver Driver
				table  string
			}{
				driver: DriverPostgres,
				table:  "secrets",
			},
			want: queries{
				selectByID:    "SELECT id, value, expires_at FROM secrets WHERE id = $1",
				insert:        "INSERT INTO secrets (id, value, expires_at) VALUES ($1, $2, $3)",
				delete:        "DELETE FROM secrets WHERE id = $1",
				deleteExpired: "DELETE FROM secrets WHERE expires_at < NOW()",
			},
		},
		{
			name: "mssql",
			input: struct {
				driver Driver
				table  string
			}{
				driver: DriverMSSQL,
				table:  "secrets",
			},
			want: queries{
				selectByID:    "SELECT id, value, expires_at FROM secrets WHERE id = @p1",
				insert:        "INSERT INTO secrets (id, value, expires_at) VALUES (@p1, @p2, @p3)",
				delete:        "DELETE FROM secrets WHERE id = @p1",
				deleteExpired: "DELETE FROM secrets WHERE expires_at < GETUTCDATE()",
			},
		},
		{
			name: "sqlite",
			input: struct {
				driver Driver
				table  string
			}{
				driver: DriverSQLite,
				table:  "secrets",
			},
			want: queries{
				selectByID:    "SELECT id, value, expires_at FROM secrets WHERE id = ?",
				insert:        "INSERT INTO secrets (id, value, expires_at) VALUES (?, ?, ?)",
				delete:        "DELETE FROM secrets WHERE id = ?",
				deleteExpired: "DELETE FROM secrets WHERE expires_at < DATETIME('now')",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := createQueries(test.input.driver, test.input.table)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(queries{})); diff != "" {
				t.Errorf("createQueries() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("createQueries() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}
