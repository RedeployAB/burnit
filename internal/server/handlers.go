package server

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/RedeployAB/burnit/internal/api"
	"github.com/RedeployAB/burnit/internal/frontend"
	"github.com/RedeployAB/burnit/internal/log"
	"github.com/RedeployAB/burnit/internal/secret"
	"github.com/RedeployAB/burnit/internal/security"
	"github.com/RedeployAB/burnit/internal/version"
)

const (
	// contentType is the content type header.
	contentType = "Content-Type"
	// contentTypeJSON is the content type for JSON.
	contentTypeJSON = "application/json; charset=UTF-8"
	// contentTypeHTML is the content type for HTML.
	contentTypeHTML = "text/html"
	// contentTypeText is the content type for text.
	contentTypeText = "text/plain"
)

// index returns a handler for handling the index route.
func index(ui frontend.UI, log log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ui != nil && strings.Contains(r.Header.Get("Accept"), contentTypeHTML) {
			http.Redirect(w, r, "/ui/secrets", http.StatusMovedPermanently)
			return
		}

		if err := encode(w, http.StatusOK, api.Index{
			Name:    "burnit",
			Version: version.Version(),
			Endpoints: []string{
				"/secret",
				"/secrets",
			},
		}); err != nil {
			log.Error("Failed to encode response.", "error", err)
			writeServerError(w)
			return
		}
	})
}

// notFound returns a handler for handling not found routes.
func notFound(ui frontend.UI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ui != nil && strings.Contains(r.Header.Get("Accept"), contentTypeHTML) {
			ui.Render(w, http.StatusNotFound, "not-found", nil)
			return
		}
		writeError(w, http.StatusNotFound, errors.New("not found"))
	})
}

// generateSecret generates a new secret.
func generateSecret(secrets secret.Service, log log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := secrets.Generate(parseGenerateSecretQuery(r.URL.Query()))

		if header := r.Header.Get("Accept"); header == contentTypeText {
			writeValue(w, secret)
			return
		}

		if err := encode(w, http.StatusOK, api.Secret{Value: secret}); err != nil {
			log.Error("Failed to encode response.", "error", err)
			writeServerError(w)
			return
		}
	})
}

// getSecret retrieves a secret.
func getSecret(secrets secret.Service, log log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if len(id) == 0 {
			writeError(w, http.StatusBadRequest, errors.New("secret ID is required"))
			return
		}

		passphrase, err := getPassphrase(r.Header)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err)
			return
		}

		secret, err := secrets.Get(id, passphrase)
		if err != nil {
			if statusCode := errorCode(err); statusCode != 0 {
				writeError(w, statusCode, err)
				return
			}
			log.Error("Failed to get secret.", "error", err)
			writeServerError(w)
			return
		}

		if err := encode(w, http.StatusOK, api.Secret{Value: secret.Value}); err != nil {
			log.Error("Failed to encode response.", "error", err)
			writeServerError(w)
			return
		}
	})
}

// createSecret creates a new secret.
func createSecret(secrets secret.Service, log log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretRequest, err := decode[api.CreateSecretRequest](r)
		if err != nil {
			writeError(w, errorCode(err), err)
			return
		}

		secret, err := secrets.Create(toCreateSecret(&secretRequest))
		if err != nil {
			if statusCode := errorCode(err); statusCode != 0 {
				writeError(w, statusCode, err)
				return
			}
			log.Error("Failed to create secret.", "error", err)
			writeServerError(w)
			return
		}

		w.Header().Set("Location", "/secrets/"+secret.ID)
		if err := encode(w, http.StatusCreated, toAPISecret(&secret)); err != nil {
			log.Error("Failed to encode response.", "error", err)
			writeServerError(w)
			return
		}
	})
}

// deleteSecret deletes a secret.
func deleteSecret(secrets secret.Service, log log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if len(id) == 0 {
			writeError(w, http.StatusBadRequest, errors.New("secret ID is required"))
			return
		}

		passphrase, err := getPassphrase(r.Header)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err)
			return
		}

		if err := secrets.Delete(id, func(o *secret.DeleteOptions) {
			o.VerifyPassphrase = true
			o.Passphrase = passphrase
		}); err != nil {
			if statusCode := errorCode(err); statusCode != 0 {
				writeError(w, statusCode, err)
				return
			}
			log.Error("Failed to delete secret.", "error", err)
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

// getPassphrase retrieves the passphrase from the headers and
// decodes it.
func getPassphrase(header http.Header) (string, error) {
	passphrase := header.Get("Passphrase")
	if len(passphrase) == 0 {
		return "", ErrPassphraseRequired
	}

	decodedPassphrase, err := security.DecodeBase64(passphrase)
	if err != nil {
		return "", ErrPassphraseNotBase64
	}

	return string(decodedPassphrase), nil
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
		TTL:        s.TTL.String(),
		ExpiresAt:  expiresAt,
	}
}
