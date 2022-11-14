package server

import (
	"encoding/json"
	"net/http"
)

// secretError represents a sercret error response.
type secretError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// JSON marshals secretError to JSON.
func (s *secretError) JSON() []byte {
	b, _ := json.Marshal(&s)
	return b
}

// newSecretError creates and returns a secretError.
func newSecretError(statusCode int, message string) *secretError {
	return &secretError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// writeError writes error to http.ResponseWriter.
func writeError(w http.ResponseWriter, statusCode int, message string) {
	if len(message) == 0 {
		switch statusCode {
		case http.StatusBadRequest:
			message = "malformed JSON"
		case http.StatusUnauthorized:
			message = "unauthorized"
		case http.StatusNotFound:
			message = "not found"
		case http.StatusInternalServerError:
			message = "internal server error"
		default:
			message = "unhandled error"
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	w.Write(newSecretError(statusCode, message).JSON())
}
