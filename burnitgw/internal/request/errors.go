package request

import (
	"strconv"
	"strings"
)

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
