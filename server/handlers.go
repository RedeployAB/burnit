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
			s.log.Error("Failed to encode response.", "error", err)
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
			s.log.Error("Failed to get secret.", "error", err)
			writeServerError(w)
			return
		}

		if err := encode(w, http.StatusOK, api.Secret{Value: secret.Value}); err != nil {
			s.log.Error("Failed to encode response.", "error", err)
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
			writeError(w, errorCode(err), err)
			return
		}

		secret, err := s.secrets.Create(toCreateSecret(&secretRequest))
		if err != nil {
			if statusCode := errorCode(err); statusCode != 0 {
				writeError(w, statusCode, err)
				return
			}
			s.log.Error("Failed to create secret.", "error", err)
			writeServerError(w)
			return
		}

		w.Header().Set("Location", "/secrets/"+secret.ID)
		if err := encode(w, http.StatusCreated, toAPISecret(&secret)); err != nil {
			s.log.Error("Failed to encode response.", "error", err)
			writeServerError(w)
			return
		}
	})
}

// deleteSecret deletes a secret.
func (s server) deleteSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := s.secrets.Delete(r.PathValue("id")); err != nil {
			if statusCode := errorCode(err); statusCode != 0 {
				writeError(w, statusCode, err)
				return
			}
			s.log.Error("Failed to delete secret.", "error", err)
			writeServerError(w)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

// toCreateSecret converts a CreateSecretRequest to a secret.
func toCreateSecret(s *api.CreateSecretRequest) secret.Secret {
	var ttl time.Duration
	if len(s.TTL) > 0 {
		var err error
		ttl, err = time.ParseDuration(s.TTL)
		if err != nil {
			// This should never happen since the request is validated.
			ttl = 5 * time.Minute
		}
	}

	return secret.Secret{
		Value:      s.Value,
		Passphrase: s.Passphrase,
		TTL:        ttl,
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
		TTL:       s.TTL.String(),
		ExpiresAt: expiresAt,
	}
}
