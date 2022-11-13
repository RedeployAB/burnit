package httperror

import (
	"encoding/json"
	"net/http"
)

// Write writes status code and error body on an http.ResponseWriter with
// JSON {"message":"error message","code":"ErrorCode","statusCode":"httpErrorCode"}.
// If empty string is provided as message, an appropriate message
// will be evauluated from provided status code.
func Write(w http.ResponseWriter, statusCode int, code, message string) {
	if len(message) == 0 {
		switch statusCode {
		case http.StatusBadRequest:
			message = "malformed JSON"
		case http.StatusUnauthorized:
			message = "unauthorized"
		case http.StatusForbidden:
			message = "forbidden"
		case http.StatusNotFound:
			message = "not found"
		case http.StatusInternalServerError:
			message = "internal server error"
		case http.StatusNotImplemented:
			message = "not implemented"
		}
	}

	er := errorResponse{StatusCode: statusCode, Code: code, Message: message}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(&er); err != nil {
		panic(err)
	}
}

// errorResponse represents an HTTP error response
// message to be encoded as JSON.
type errorResponse struct {
	Message    string `json:"message,omitempty"`
	Code       string `json:"code,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
}
