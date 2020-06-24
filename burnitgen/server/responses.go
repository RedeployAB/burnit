package server

// responseBody represents a secret response.
type responseBody struct {
	Secret secret `json:"secret"`
}

// secret represents the data part of the response body.
type secret struct {
	Value string `json:"value"`
}

// response creates a response from a Secret (string).
func response(s string) *responseBody {
	return &responseBody{Secret: secret{Value: s}}
}
