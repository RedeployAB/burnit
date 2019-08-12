package api

// SecretResponseBody represents a secret response.
type SecretResponseBody struct {
	Data secretData `json:"data"`
}

type secretData struct {
	Secret string `json:"secret"`
}
