package internal

import (
	"fmt"
	"strconv"
	"strings"
)

// DBRequestError implements error interface and represents
// errors to secretdb
type DBRequestError struct {
	err  string
	code int
}

func (e *DBRequestError) Error() string {
	return fmt.Sprintf("code %d: %s", e.code, e.err)
}

// HandleHTTPError takes an error and parses it for HTTP errors.
func HandleHTTPError(err error) int {
	var status int
	errStr := err.Error()
	if strings.HasPrefix(errStr, "code ") {
		status, err = strconv.Atoi(errStr[5:8])
		if err != nil {
			status = 500
		}
	} else {
		status = 500
	}

	return status
}
