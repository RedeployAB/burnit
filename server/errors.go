package server

import (
	"errors"
	"net/http"

	"github.com/RedeployAB/burnit/secret"
)

var (
	// ErrEmptyRequest is returned when the request body is empty.
	ErrEmptyRequest = errors.New("empty request")
	// ErrMalformedRequest is returned when the request body is invalid.
	ErrMalformedRequest = errors.New("malformed request")
	// ErrInvalidRequest is returned when the request is invalid.
	ErrInvalidRequest = errors.New("invalid request")
	// ErrInvalidPath is returned when the path is invalid.
	ErrInvalidPath = errors.New("invalid path")
)

// responseError represents an error response.
type responseError struct {
	StatusCode int    `json:"statusCode"`
	Err        string `json:"error"`
}

// Error returns the error message.
func (e *responseError) Error() string {
	return e.Err
}

// writeError writes an error response to the caller.
func writeError(w http.ResponseWriter, statusCode int, err error) {
	if err == nil {
		err = errors.New("internal server error")
	}

	respErr := &responseError{
		StatusCode: statusCode,
		Err:        err.Error(),
	}

	if err := encode(w, statusCode, respErr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// writeServerError writes a server error response to the caller.
func writeServerError(w http.ResponseWriter) {
	writeError(w, http.StatusInternalServerError, nil)
}

// errorCode returns the status code for the given error.
func errorCode(err error) int {
	for statusCode, errs := range errorCodeMaps {
		for _, e := range errs {
			if errors.Is(err, e) {
				return statusCode
			}
		}
	}
	return 0
}

// errorCodeMaps maps errors to status codes.
var errorCodeMaps = map[int][]error{
	http.StatusBadRequest: {
		ErrEmptyRequest,
		ErrInvalidRequest,
		ErrMalformedRequest,
	},
	http.StatusUnauthorized: {
		secret.ErrInvalidPassphrase,
	},
	http.StatusNotFound: {
		secret.ErrSecretNotFound,
	},
}
