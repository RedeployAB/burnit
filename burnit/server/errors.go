package server

import (
	"encoding/json"
	"net/http"
)

// errorResponse represents an error response.
type errorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// JSON marshals errorResponse to JSON.
func (s *errorResponse) JSON() []byte {
	b, _ := json.Marshal(&s)
	return b
}

// newErrorResponse creates and returns a secretError.
func newErrorResponse(statusCode int, message string) *errorResponse {
	return &errorResponse{
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
	w.Write(newErrorResponse(statusCode, message).JSON())
}
