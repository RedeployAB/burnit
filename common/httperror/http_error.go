package httperror

import (
	"encoding/json"
	"net/http"
)

// Error works like the http packages Error, but it responds
// with JSON {"error":"error message", "code":"httpErrorCode"}.
func Error(w http.ResponseWriter, error string, code int) {
	w.WriteHeader(code)
	er := errorResponse{Error: error, Code: code}
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
