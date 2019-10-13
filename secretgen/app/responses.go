package app

// secretResponse represents a secret response.
type secretResponse struct {
	Data secretData `json:"data"`
}

// secret represents a secret.
type secretData struct {
	Secret string `json:"secret"`
}
