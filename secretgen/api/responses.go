package api

// secretResponse represents a secret response.
type secretResponseBody struct {
	Data secretResponse `json:"data"`
}

// secret represents a secret.
type secretResponse struct {
	Secret string `json:"secret"`
}
