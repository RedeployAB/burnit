package sql

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCreateSessionQueries(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			driver Driver
			table  string
		}
		want    sessionQueries
		wantErr error
	}{
		{
			name: "postgres",
			input: struct {
				driver Driver
				table  string
			}{
				driver: DriverPostgres,
				table:  "sessions",
			},
			want: sessionQueries{
				selectByID:        "SELECT id, expires_at, csrf_token, csrf_expires_at FROM sessions WHERE id = $1",
				selectByCSRFToken: "SELECT id, expires_at, csrf_token, csrf_expires_at FROM sessions WHERE csrf_token = $1",
				upsert:            "INSERT INTO sessions (id, expires_at, csrf_token, csrf_expires_at) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET expires_at = EXCLUDED.expires_at, csrf_token = EXCLUDED.csrf_token, csrf_expires_at = EXCLUDED.csrf_expires_at",
				delete:            "DELETE FROM sessions WHERE id = $1",
				deleteExpired:     "DELETE FROM sessions WHERE expires_at < NOW() AT TIME ZONE 'UTC'",
			},
		},
		{
			name: "mssql",
			input: struct {
				driver Driver
				table  string
			}{
				driver: DriverMSSQL,
				table:  "sessions",
			},
			want: sessionQueries{
				selectByID:        "SELECT ID, ExpiresAt, CSRFToken, CSRFExpiresAt FROM Sessions WHERE ID = @p1",
				selectByCSRFToken: "SELECT ID, ExpiresAt, CSRFToken, CSRFExpiresAt FROM Sessions WHERE CSRFToken = @p1",
				upsert:            "MERGE INTO Sessions AS target USING (VALUES (@p1, @p2, @p3, @p4)) AS source (ID, ExpiresAt, CSRFToken, CSRFExpiresAt) ON target.ID = source.ID WHEN MATCHED THEN UPDATE SET target.ExpiresAt = source.ExpiresAt, target.CSRFToken = source.CSRFToken, target.CSRFExpiresAt = source.CSRFExpiresAt WHEN NOT MATCHED THEN INSERT (ID, ExpiresAt, CSRFToken, CSRFExpiresAt) VALUES (source.ID, source.ExpiresAt, source.CSRFToken, source.CSRFExpiresAt);",
				delete:            "DELETE FROM Sessions WHERE ID = @p1",
				deleteExpired:     "DELETE FROM Sessions WHERE ExpiresAt < GETUTCDATE()",
			},
		},
		{
			name: "sqlite",
			input: struct {
				driver Driver
				table  string
			}{
				driver: DriverSQLite,
				table:  "sessions",
			},
			want: sessionQueries{
				selectByID:        "SELECT id, expires_at, csrf_token, csrf_expires_at FROM sessions WHERE id = ?1",
				selectByCSRFToken: "SELECT id, expires_at, csrf_token, csrf_expires_at FROM sessions WHERE csrf_token = ?1",
				upsert:            "INSERT INTO sessions (id, expires_at, csrf_token, csrf_expires_at) VALUES (?1, ?2, ?3, ?4) ON CONFLICT(id) DO UPDATE SET expires_at = excluded.expires_at, csrf_token = excluded.csrf_token, csrf_expires_at = excluded.csrf_expires_at",
				delete:            "DELETE FROM sessions WHERE id = ?1",
				deleteExpired:     "DELETE FROM sessions WHERE expires_at < DATETIME('now')",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := createSessionQueries(test.input.driver, test.input.table)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(sessionQueries{})); diff != "" {
				t.Errorf("createSessionQueries() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("createSessionQueries() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}
