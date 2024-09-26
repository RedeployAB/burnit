package server

import (
	"net/http"
	"time"

	"github.com/RedeployAB/burnit/api"
	"github.com/RedeployAB/burnit/secret"
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var passphrase string
		if p := r.PathValue("passphrase"); len(p) > 0 {
			passphrase = p
		} else if p := r.Header.Get("Passphrase"); len(p) > 0 {
			passphrase = p
		}

		secret, err := s.secrets.Get(r.PathValue("id"), passphrase)
		if err != nil {
			if statusCode := errorCode(err); statusCode != 0 {
				writeError(w, statusCode, err)
				return
			}
			writeServerError(w)
			return
		}

		if err := encode(w, http.StatusOK, api.Secret{Value: secret.Value}); err != nil {
			writeServerError(w)
			return
		}
	})
}

// createSecret creates a new secret.
func (s server) createSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretRequest, err := decode[api.CreateSecretRequest](r)
		if err != nil {
			statusCode := errorCode(err)
			writeError(w, statusCode, err)
			return
		}

		secret, err := s.secrets.Create(toCreateSecret(&secretRequest))
		if err != nil {
			if statusCode := errorCode(err); statusCode != 0 {
				writeError(w, statusCode, err)
				return
			}
			writeServerError(w)
			return
		}

		w.Header().Set("Location", "/secrets/"+secret.ID)
		if err := encode(w, http.StatusCreated, toAPISecret(&secret)); err != nil {
			writeServerError(w)
			return
		}
	})
}

// toCreateSecret converts a CreateSecretRequest to a secret.
func toCreateSecret(s *api.CreateSecretRequest) secret.Secret {
	return secret.Secret{
		Value:      s.Value,
		Passphrase: s.Passphrase,
		TTL:        s.TTL,
	}
}

// toAPISecret converts a secret to an API secret.
func toAPISecret(s *secret.Secret) api.Secret {
	var expiresAt *time.Time
	if !s.ExpiresAt.IsZero() {
		expiresAt = &s.ExpiresAt
	}

	return api.Secret{
		ID:        s.ID,
		TTL:       s.TTL,
		ExpiresAt: expiresAt,
	}
}
