package server

// secretResponse defines the structure of an outgoing
// secret response.
type secretResponse struct {
	Value string `json:"value"`
}

// newSecretResponse creates a response from a secret (string).
func newSecretResponse(s string) *secretResponse {
	return &secretResponse{Value: s}
}
