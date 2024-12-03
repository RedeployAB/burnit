package ui

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/RedeployAB/burnit/internal/log"
	"github.com/RedeployAB/burnit/internal/middleware"
	"github.com/RedeployAB/burnit/internal/secret"
	"github.com/RedeployAB/burnit/internal/security"
	"github.com/RedeployAB/burnit/internal/session"
)

// Index handles requests to the index route.
func Index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/secrets", http.StatusMovedPermanently)
	})
}

// NotFound handles requests to not found routes.
func NotFound(ui UI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ui.Render(w, http.StatusNotFound, "not-found", nil)
	})
}

// CreateSecret handles requests to create a secret.
func CreateSecret(ui UI, secrets secret.Service, sessions session.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sessions are only implemented for CSRF tokens at the moment.
		// Use the CSRF token as the session ID when setting the session.
		sess := session.NewSession(session.WithCSRF(session.NewCSRF()))
		sessions.Set(sess.CSRF().Token(), sess)
		ui.Render(w, http.StatusOK, "secret-create", secretCreateResponse{CSRFToken: sess.CSRF().Token()})
	})
}

// GetSecret handles requests to get a secret.
func GetSecret(ui UI, secrets secret.Service, sessions session.Store, log log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, passphrase, err := extractIDAndPassphrase("/ui/secrets/", r.URL.Path)
		if err != nil || len(id) == 0 {
			http.Redirect(w, r, "/ui/secrets", http.StatusFound)
			return
		}

		if _, err = secrets.Get(id, passphrase, func(o *secret.GetOptions) {
			o.NoDecrypt = true
		}); err != nil {
			if errors.Is(err, secret.ErrSecretNotFound) {
				ui.Render(w, http.StatusNotFound, "secret-not-found", nil)
				return
			}

			requestID := requestIDFromContext(r.Context())
			log.Error("Failed to get secret.", uiLog(err, "GetSecret", requestID)...)
			ui.Render(w, http.StatusInternalServerError, "error", errorResponse{Title: "An error occured", Message: "Could not retrieve secret.", RequestID: requestID}, WithPartial())
			return
		}

		if len(passphrase) == 0 {
			// Sessions are only implemented for CSRF tokens at the moment.
			// Use the CSRF token as the session ID when setting the session.
			sess := session.NewSession(session.WithCSRF(session.NewCSRF()))
			sessions.Set(sess.CSRF().Token(), sess)
			ui.Render(w, http.StatusUnauthorized, "secret-get", secretGetResponse{ID: id, CSRFToken: sess.CSRF().Token()})
			return
		}

		decodedPassphrase, err := security.DecodeBase64(passphrase)
		if err != nil {
			ui.Render(w, http.StatusBadRequest, "error", errorResponse{Title: "Could not retrieve secret", Message: "Invalid passphrase."}, WithPartial())
			return
		}

		s, err := secrets.Get(id, string(decodedPassphrase), func(o *secret.GetOptions) {
			o.PassphraseHashed = true
		})
		if err != nil {
			if errors.Is(err, secret.ErrSecretNotFound) {
				ui.Render(w, http.StatusNotFound, "secret-not-found", nil)
				return
			}

			requestID := requestIDFromContext(r.Context())
			log.Error("Failed to get secret.", uiLog(err, "GetSecret", requestID)...)
			ui.Render(w, http.StatusInternalServerError, "error", errorResponse{Title: "An error occured", Message: "Could not retrieve secret.", RequestID: requestID}, WithPartial())
			return
		}

		response := secretGetResponse{
			ID:             s.ID,
			PassphraseHash: passphrase,
			Value:          s.Value,
		}

		ui.Render(w, http.StatusOK, "secret-get", response)
	})
}

// HandlerCreateSecret handles requests containing a form to create a secret.
func HandlerCreateSecret(ui UI, secrets secret.Service, sessions session.Store, log log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			sess := session.NewSession(session.WithCSRF(session.NewCSRF()))
			sessions.Set(sess.CSRF().Token(), sess)
			ui.Render(w, http.StatusOK, "secret-create", secretCreateResponse{CSRFToken: sess.CSRF().Token()}, WithPartial())
			return
		}
		if err := r.ParseForm(); err != nil {
			requestID := requestIDFromContext(r.Context())
			log.Error("Failed to parse form.", uiLog(err, "HandlerCreateSecret", requestID)...)
			ui.Render(w, http.StatusBadRequest, "error", errorResponse{Title: "An error occured", Message: "Could not parse form.", RequestID: requestID}, WithPartial())
			return
		}

		ok, statusCode, errResp, err := validateCSRFTToken(r.Context(), sessions, r.FormValue("csrf-token"))
		if err != nil {
			requestID := requestIDFromContext(r.Context())
			log.Error("Failed to validate CSRF token.", uiLog(err, "HandlerCreateSecret", requestID)...)
			ui.Render(w, statusCode, "error", errResp, WithPartial())
			return
		}
		if !ok {
			ui.Render(w, statusCode, "error", errResp, WithPartial())
			return
		}

		baseURL := r.FormValue("base-url")
		if len(baseURL) == 0 {
			ui.Render(w, http.StatusBadRequest, "error", errorResponse{Title: "An error occured", Message: "Missing base URL."}, WithPartial())
			return
		}

		ttl, err := time.ParseDuration(r.FormValue("ttl"))
		if err != nil {
			ui.Render(w, http.StatusBadRequest, "error", errorResponse{Title: "Could not create secret", Message: "Invalid expiration time."}, WithPartial())
			return
		}

		s, err := secrets.Create(secret.Secret{
			Value:      r.FormValue("value"),
			Passphrase: r.FormValue("custom-value"),
			TTL:        ttl,
		})
		if err != nil {
			var response errorResponse
			var statusCode int
			if !isSecretBadRequestError(err) {
				statusCode = http.StatusInternalServerError
				requestID := requestIDFromContext(r.Context())
				response = errorResponse{Title: "An error occured", Message: "Internal server error.", RequestID: requestID}
				log.Error("Failed to create secret.", uiLog(err, "HandlerCreateSecret", requestID)...)
			} else {
				statusCode = http.StatusBadRequest
				response = errorResponse{Title: "Could not create secret", Message: formatErrorMessage(err)}
			}

			ui.Render(w, statusCode, "error", response, WithPartial())
			return
		}

		if err := sessions.Delete(r.FormValue("csrf-token")); err != nil {
			log.Error("Failed to delete session.", uiLog(err, "HandlerCreateSecret", requestIDFromContext(r.Context()))...)
		}

		response := secretCreateResponse{
			BaseURL:        baseURL,
			ID:             s.ID,
			Passphrase:     s.Passphrase,
			PassphraseHash: base64.RawURLEncoding.EncodeToString(security.SHA256([]byte(s.Passphrase))),
		}

		ui.Render(w, http.StatusCreated, "secret-created", response, WithPartial())
	})
}

// HandlerGetSecret handles requests containing a form to get a secret.
// This form will be used when a passphrase is not provided in the URL.
func HandlerGetSecret(ui UI, secrets secret.Service, sessions session.Store, log log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			requestID := middleware.RequestIDFromContext(r.Context())
			log.Error("Failed to parse form.", uiLog(err, "HandlerGetSecret", requestID)...)
			ui.Render(w, http.StatusInternalServerError, "error", errorResponse{Title: "An error occured", Message: "Could not parse form.", RequestID: requestID}, WithPartial())
			return
		}

		ok, statusCode, errResp, err := validateCSRFTToken(r.Context(), sessions, r.FormValue("csrf-token"))
		if err != nil {
			log.Error("Failed to validate CSRF token.", uiLog(err, "HandlerGetSecret", requestIDFromContext(r.Context()))...)
			ui.Render(w, statusCode, "error", errResp, WithPartial())
			return
		}
		if !ok {
			ui.Render(w, statusCode, "error", errResp, WithPartial())
			return
		}

		id := r.FormValue("id")
		if len(id) == 0 {
			requestID := middleware.RequestIDFromContext(r.Context())
			log.Error("Missing ID in request.", uiLog(err, "HandlerGetSecret", requestID)...)
			ui.Render(w, http.StatusInternalServerError, "error", errorResponse{Title: "An error occured", Message: "Missing ID.", RequestID: requestID}, WithPartial())
			return
		}
		passphrase := r.FormValue("custom-value")
		if len(passphrase) == 0 {
			ui.Render(w, http.StatusOK, "secret-get", secretGetResponse{ID: id}, WithPartial())
			return
		}

		s, err := secrets.Get(id, passphrase)
		if err != nil {
			if errors.Is(err, secret.ErrSecretNotFound) {
				ui.Render(w, http.StatusNotFound, "secret-not-found", nil)
				return
			}
			if errors.Is(err, secret.ErrInvalidPassphrase) {
				ui.Render(w, http.StatusUnauthorized, "secret-get", secretGetResponse{ID: id, CSRFToken: r.FormValue("csrf-token")}, WithPartial())
				return
			}

			requestID := middleware.RequestIDFromContext(r.Context())
			log.Error("Failed to get secret.", uiLog(err, "HandlerGetSecret", requestID)...)
			ui.Render(w, http.StatusInternalServerError, "error", errorResponse{Title: "An error occured", Message: "Could not retrieve secret.", RequestID: requestID}, WithPartial())
			return
		}

		if err := sessions.Delete(r.FormValue("csrf-token")); err != nil {
			log.Error("Failed to delete session.", "handler", "HandlerGetSecret", "error", err, "requestId", requestIDFromContext(r.Context()))
		}

		response := secretGetResponse{
			ID:             s.ID,
			PassphraseHash: passphrase,
			Value:          s.Value,
		}

		ui.Render(w, http.StatusOK, "secret-get", response, WithPartial())
	})
}

// extractIDAndPassphrase extracts the ID and passphrase from the path.
func extractIDAndPassphrase(route, path string) (string, string, error) {
	path = strings.TrimPrefix(path, route)
	parts := strings.Split(path, "/")

	if len(parts) > 2 {
		return "", "", fmt.Errorf("invalid path: %s", path)
	}
	if len(parts) == 1 {
		return parts[0], "", nil
	}
	return parts[0], parts[1], nil
}

// secretCreateResponse is the response data for a create secret request.
type secretCreateResponse struct {
	BaseURL        string
	ID             string
	Passphrase     string
	PassphraseHash string
	CSRFToken      string
}

// secretGetResponse is the response data for a get secret request.
type secretGetResponse struct {
	ID             string
	PassphraseHash string
	Value          string
	CSRFToken      string
}

// errorResponse is the response data for an error.
type errorResponse struct {
	RequestID string
	Title     string
	Message   string
}

// formatErrorMessage formats an error message.
func formatErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	re := regexp.MustCompile(`^.*:\s+`)
	msg = re.ReplaceAllString(msg, "")
	if !strings.HasSuffix(msg, ".") {
		msg += "."
	}
	return strings.ToUpper(msg[:1]) + msg[1:]
}

// isSecretBadRequestError returns true if the error is a bad request error.
func isSecretBadRequestError(err error) bool {
	errs := []error{
		secret.ErrValueInvalid,
		secret.ErrInvalidPassphrase,
		secret.ErrValueTooManyCharacters,
		secret.ErrInvalidExpirationTime,
		secret.ErrPassphraseInvalid,
		secret.ErrPassphraseTooFewCharacters,
		secret.ErrPassphraseTooManyCharacters,
	}

	for _, e := range errs {
		if errors.Is(err, e) {
			return true
		}
	}
	return false
}

// validateCSRFTToken validates the CSRF token in the session.
// When the token has been validated, and the session is no longer needed,
// and the session is deleted. For future implementations, the session
// should be deleted when the session is no longer needed.
//
// In this implementation the CSRF token is the session ID, since
// sessions have only been implemented for CSRF tokens.
func validateCSRFTToken(ctx context.Context, sessions session.Store, id string) (bool, int, errorResponse, error) {
	sess, err := sessions.Get(id)
	if err != nil {
		title := "Could not retrieve session"
		if errors.Is(err, session.ErrSessionNotFound) {
			return false, http.StatusBadRequest, errorResponse{Title: title, Message: "Session not found."}, errors.New("session not found")
		}
		if errors.Is(err, session.ErrSessionExpired) {
			return false, http.StatusBadRequest, errorResponse{Title: title, Message: "Session expired. Please refresh the page."}, nil
		}
		return false, http.StatusInternalServerError, errorResponse{Title: "An error occured", Message: "Could not validate CSRF token.", RequestID: middleware.RequestIDFromContext(ctx)}, err
	}
	title := "Invalid CSRF token"
	csrf := sess.CSRF()
	t := csrf.Token()
	if len(t) == 0 {
		return false, http.StatusBadRequest, errorResponse{Title: title, Message: "CSRF token not found."}, errors.New("CSRF token not found")
	}
	if t != id {
		return false, http.StatusBadRequest, errorResponse{Title: title, Message: "CSRF token does not match with session."}, errors.New("CSRF token does not match with session")
	}
	if csrf.Expired() {
		return false, http.StatusBadRequest, errorResponse{Title: title, Message: "CSRF token expired. Please refresh the page."}, nil
	}
	return true, 0, errorResponse{}, nil
}

// uiLog formats the log message for the UI.
func uiLog(err error, handler, requestID string) []any {
	return []any{"type", "ui", "handler", handler, "error", err, "requestId", requestID}
}

// requestIDFromContext returns the request ID from the context.
func requestIDFromContext(ctx context.Context) string {
	return middleware.RequestIDFromContext(ctx)
}
