package api

// SecretBody represents incoming request body.
type SecretBody struct {
	Secret     string `json:"secret"`
	Passphrase string `json:"passphrase"`
}
