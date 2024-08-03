package server

import (
	"net/http"

	"github.com/RedeployAB/burnit/api"
)

const (
	// contentType is the content type header.
	contentType = "Content-Type"
	// contentTypeJSON is the content type for JSON.
	contentTypeJSON = "application/json; charset=UTF-8"
	// contentTypeText is the content type for text.
	contentTypeText = "text/plain"
)

// generateSecret generates a new secret.
func (s server) generateSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := s.secrets.Generate(parseGenerateSecretQuery(r.URL.Query()))

		if header := r.Header.Get("Accept"); header == contentTypeText {
			writeValue(w, secret)
			return
		}

		if err := encode(w, http.StatusOK, api.Secret{Value: secret}); err != nil {
			writeServerError(w)
			return
		}
	})
}

// getSecret retrieves a secret.
func (s server) getSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

// createSecret creates a new secret.
func (s server) createSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
