package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/RedeployAB/burnit/internal/api"
	"github.com/RedeployAB/burnit/internal/secret"
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
		id, passphrase, err := extractIDAndPassphrase("/secrets/", r.URL.Path)
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		if len(passphrase) == 0 {
			p := r.Header.Get("Passphrase")
			passphrase, err = decodeBase64(p)
			if err != nil {
				writeError(w, http.StatusBadRequest, errors.New("invalid passphrase: should be base64 encoded"))
				return
			}
		}

		secret, err := s.secrets.Get(id, passphrase)
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
		id, passphrase, err := extractIDAndPassphrase("/secrets/", r.URL.Path)
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		if len(passphrase) == 0 {
			passphrase = r.Header.Get("Passphrase")
		}

		if err := s.secrets.Delete(id, func(o *secret.DeleteOptions) {
			o.VerifyPassphrase = true
			o.Passphrase = passphrase
		}); err != nil {
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
	var expiresAt time.Time
	if s.ExpiresAt != nil {
		expiresAt = s.ExpiresAt.Time
	}

	return secret.Secret{
		Value:      s.Value,
		Passphrase: s.Passphrase,
		TTL:        ttl,
		ExpiresAt:  expiresAt,
	}
}

// toAPISecret converts a secret to an API secret.
func toAPISecret(s *secret.Secret) api.Secret {
	var expiresAt *api.Time
	if !s.ExpiresAt.IsZero() {
		expiresAt = &api.Time{Time: s.ExpiresAt}
	}

	return api.Secret{
		ID:         s.ID,
		Passphrase: s.Passphrase,
		Path:       "/secrets/" + s.ID + "/" + s.PassphraseHash,
		TTL:        s.TTL.String(),
		ExpiresAt:  expiresAt,
	}
}
