package httperror

import (
	"encoding/json"
	"net/http"
)

// Error works like the http packages Error, but it responds
// with JSON {"message":"error message", "statusCode":"httpErrorCode"}.
func Error(w http.ResponseWriter, statusCode int) {
	var err string
	switch statusCode {
	case http.StatusBadRequest:
		err = "malformed JSON"
	case http.StatusUnauthorized:
		err = "unauthorized"
	case http.StatusForbidden:
		err = "forbidden"
	case http.StatusNotFound:
		err = "not found"
	case http.StatusInternalServerError:
		err = "internal server error"
	case http.StatusNotImplemented:
		err = "not implemented"
	}
	er := errorResponse{Message: err, StatusCode: statusCode}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(&er); err != nil {
		panic(err)
	}
}

// authenticationErrorResponse represents an HTTP error response
// message to be encoded as JSON.
type errorResponse struct {
	Message    string `json:"message,omitempty"`
	Code       string `json:"code,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
}
