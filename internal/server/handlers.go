package server

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/RedeployAB/burnit/internal/api"
	"github.com/RedeployAB/burnit/internal/log"
	"github.com/RedeployAB/burnit/internal/middleware"
	"github.com/RedeployAB/burnit/internal/secret"
	"github.com/RedeployAB/burnit/internal/security"
	"github.com/RedeployAB/burnit/internal/ui"
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
func index(ui ui.UI, log log.Logger) http.Handler {
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
			requestID := requestIDFromContext(r.Context())
			log.Error("Failed to encode response.", serviceLog(err, "index", requestID)...)
			writeServerError(w, requestID)
			return
		}
	})
}

// notFound returns a handler for handling not found routes.
func notFound(ui ui.UI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ui != nil && strings.Contains(r.Header.Get("Accept"), contentTypeHTML) {
			ui.Render(w, http.StatusNotFound, "not-found", nil)
			return
		}
		writeError(w, errors.New("not found"), http.StatusNotFound, "NotFound")
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
			requestID := requestIDFromContext(r.Context())
			log.Error("Failed to encode response.", serviceLog(err, "generateSecret", requestID)...)
			writeServerError(w, requestID)
			return
		}
	})
}

// getSecret retrieves a secret.
func getSecret(secrets secret.Service, log log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if len(id) == 0 {
			writeError(w, errors.New("secret ID is required"), http.StatusBadRequest, "SecretIDRequired")
			return
		}

		passphrase, err := getPassphrase(r.Header)
		if err != nil {
			statusCode, code := errorCode(err)
			writeError(w, err, statusCode, code)
			return
		}

		secret, err := secrets.Get(id, passphrase)
		if err != nil {
			if statusCode, code := errorCode(err); statusCode != 0 {
				writeError(w, err, statusCode, code)
				return
			}
			requestID := requestIDFromContext(r.Context())
			log.Error("Failed to get secret.", serviceLog(err, "getSecret", requestID)...)
			writeServerError(w, requestID)
			return
		}

		if err := encode(w, http.StatusOK, api.Secret{Value: secret.Value}); err != nil {
			requestID := requestIDFromContext(r.Context())
			log.Error("Failed to encode response.", serviceLog(err, "getSecret", requestID)...)
			writeServerError(w, requestID)
			return
		}
	})
}

// createSecret creates a new secret.
func createSecret(secrets secret.Service, log log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretRequest, err := decode[api.CreateSecretRequest](r)
		if err != nil {
			statusCode, code := errorCode(err)
			writeError(w, err, statusCode, code)
			return
		}

		secret, err := secrets.Create(toCreateSecret(&secretRequest))
		if err != nil {
			if statusCode, code := errorCode(err); statusCode != 0 {
				writeError(w, err, statusCode, code)
				return
			}
			requestID := requestIDFromContext(r.Context())
			log.Error("Failed to create secret.", serviceLog(err, "createSecret", requestID)...)
			writeServerError(w, requestID)
			return
		}

		w.Header().Set("Location", "/secrets/"+secret.ID)
		if err := encode(w, http.StatusCreated, toAPISecret(&secret)); err != nil {
			requestID := requestIDFromContext(r.Context())
			log.Error("Failed to encode response.", serviceLog(err, "createSecret", requestID)...)
			writeServerError(w, requestID)
			return
		}
	})
}

// deleteSecret deletes a secret.
func deleteSecret(secrets secret.Service, log log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if len(id) == 0 {
			writeError(w, errors.New("secret ID is required"), http.StatusBadRequest, "SecretIDRequired")
			return
		}

		passphrase, err := getPassphrase(r.Header)
		if err != nil {
			statusCode, code := errorCode(err)
			writeError(w, err, statusCode, code)
			return
		}

		if err := secrets.Delete(id, func(o *secret.DeleteOptions) {
			o.VerifyPassphrase = true
			o.Passphrase = passphrase
		}); err != nil {
			if statusCode, code := errorCode(err); statusCode != 0 {
				writeError(w, err, statusCode, code)
				return
			}
			requestID := requestIDFromContext(r.Context())
			log.Error("Failed to delete secret.", serviceLog(err, "deleteSecret", requestID)...)
			writeServerError(w, requestID)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

// parseGenerateSecretQuery parses the query parameters for length
// and special characters.
func parseGenerateSecretQuery(v url.Values) secret.GenerateOption {
	var length string
	if l, ok := v["length"]; ok {
		length = l[0]
	}
	if l, ok := v["l"]; ok {
		length = l[0]
	}

	var specialCharacters string
	if sc, ok := v["specialCharacters"]; ok {
		specialCharacters = sc[0]
	}
	if sc, ok := v["sc"]; ok {
		specialCharacters = sc[0]
	}

	l, err := strconv.Atoi(length)
	if err != nil {
		l = defaultLength
	}
	sc, err := strconv.ParseBool(specialCharacters)
	if err != nil {
		sc = false
	}

	return func(o *secret.GenerateOptions) {
		o.Length = l
		o.SpecialCharacters = sc
	}
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

// serviceLog formats the log message for a service.
func serviceLog(err error, handler, requestID string) []any {
	return []any{"type", "service", "handler", handler, "error", err, "requestId", requestID}
}

// requestIDFromContext returns the request ID from the context.
func requestIDFromContext(ctx context.Context) string {
	return middleware.RequestIDFromContext(ctx)
}
