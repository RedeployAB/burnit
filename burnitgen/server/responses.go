package server

// secretResponse defines the structure of an outgoing
// secret response.
type secretResponse struct {
	Secret secretBody `json:"secret"`
}

// secretBody defines the structure of an outgoing
// secret response body.
type secretBody struct {
	Value string `json:"value"`
}

// newSecretResponse creates a response from a secret (string).
func newSecretResponse(s string) *secretResponse {
	return &secretResponse{Secret: secretBody{Value: s}}
}
