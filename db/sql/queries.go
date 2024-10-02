package sql

import (
	"fmt"
)

// querie contains queries used by the repository.
type queries struct {
	selectByID    string
	insert        string
	delete        string
	deleteExpired string
}

// createQueries creates the queries used by the repository.
func createQueries(driver Driver, table string) (queries, error) {
	var placeholders []string
	var now string
	switch driver {
	case DriverPostgres:
		placeholders = []string{"$1", "$2", "$3"}
		now = "NOW()"
	case DriverMSSQL:
		placeholders = []string{"@p1", "@p2", "@p3"}
		now = "GETUTCDATE()"
	case DriverMySQL, DriverMariaDB:
		placeholders = []string{"?", "?", "?"}
		now = "NOW()"
	case DriverSQLite:
		placeholders = []string{"?", "?", "?"}
		now = "DATETIME('now')"
	default:
		return queries{}, fmt.Errorf("unsupported driver: %s", driver)
	}

	return queries{
		selectByID:    fmt.Sprintf("SELECT id, value, expires_at FROM %s WHERE id = %s", table, placeholders[0]),
		insert:        fmt.Sprintf("INSERT INTO %s (id, value, expires_at) VALUES (%s, %s, %s)", table, placeholders[0], placeholders[1], placeholders[2]),
		delete:        fmt.Sprintf("DELETE FROM %s WHERE id = %s", table, placeholders[0]),
		deleteExpired: fmt.Sprintf("DELETE FROM %s WHERE expires_at < %s", table, now),
	}, nil
}
