package server

// responseBody represents a secret response.
type responseBody struct {
	Data data `json:"data"`
}

// data represents the data part of the response body.
type data struct {
	Secret string `json:"secret"`
}

// response creates a response from a Secret (string).
func response(s string) *responseBody {
	return &responseBody{Data: data{Secret: s}}
}
