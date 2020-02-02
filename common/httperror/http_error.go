package httperror

import (
	"encoding/json"
	"net/http"
)

// Error works like the http packages Error, but it responds
// with JSON {"error":"error message", "code":"httpErrorCode"}.
func Error(w http.ResponseWriter, code int) {
	var err string
	switch code {
	case http.StatusNotFound:
		err = "not found"
	case http.StatusBadRequest:
		err = "malformed JSON"
	case http.StatusUnauthorized:
		err = "unauthorized"
	case http.StatusForbidden:
		err = "forbidden"
	case http.StatusInternalServerError:
		err = "internal server error"
	case http.StatusNotImplemented:
		err = "not implemented"
	}
	er := errorResponse{Error: err, Code: code}
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(&er); err != nil {
		panic(err)
	}
}

// authenticationErrorResponse represents an HTTP error response
// message to be encoded as JSON.
type errorResponse struct {
	Error string `json:"error,omitempty"`
	Code  int    `json:"code,omitempty"`
}
