package server

import (
	"errors"
	"net/http"
)

var (
	// ErrEmptyRequest is returned when the request body is empty.
	ErrEmptyRequest = errors.New("empty request")
	// ErrMalformedRequest is returned when the request body is invalid.
	ErrMalformedRequest = errors.New("malformed request")
	// ErrInvalidRequest is returned when the request is invalid.
	ErrInvalidRequest = errors.New("invalid request")
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
