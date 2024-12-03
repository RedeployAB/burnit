package api

import "encoding/json"

// Error represents an error response.
type Error struct {
	StatusCode int    `json:"statusCode"`
	Code       string `json:"code,omitempty"`
	Err        string `json:"error"`
	RequestID  string `json:"requestId,omitempty"`
}

// Error returns the error message.
func (e Error) Error() string {
	return e.Err
}

// JSON returns the JSON encoded error.
func (e Error) JSON() []byte {
	b, _ := json.Marshal(e)
	return b
}
