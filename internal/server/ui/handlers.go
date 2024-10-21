package ui

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/RedeployAB/burnit/internal/secret"
)

var (
	//go:embed templates/*
	TemplateFS embed.FS

	//go:embed static/*
	StaticFS embed.FS
)

// secretService is an interface for the secret service.
type secretService interface {
	Get(id, passphrase string, options ...secret.GetOption) (secret.Secret, error)
	Create(secret secret.Secret) (secret.Secret, error)
	Delete(id string, options ...secret.DeleteOption) error
}

// CreateSecret handles requests to create a secret.
func CreateSecret(secrets secretService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, http.StatusOK, "secret-create", nil)
	})
}

// GetSecret handles requests to get a secret.
func GetSecret(secrets secretService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, passphrase, err := extractIDAndPassphrase("/ui/secrets/", r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		s, err := secrets.Get(id, passphrase)
		if err != nil {
			if errors.Is(err, secret.ErrSecretNotFound) {
				renderTemplate(w, http.StatusNotFound, "secret-not-found", nil)
				return
			}
			if errors.Is(err, secret.ErrInvalidPassphrase) {
				renderTemplate(w, http.StatusUnauthorized, "secret-get", secretGetResponse{ID: id})
				return
			}
			http.Error(w, "could not get secret: error in service", http.StatusInternalServerError)
			return
		}

		response := secretGetResponse{
			ID:             s.ID,
			PassphraseHash: passphrase,
			Value:          s.Value,
		}
		renderTemplate(w, http.StatusOK, "secret-get", response)
	})
}

// HandlerCreateSecret handles requests containing a form to create a secret.
func HandlerCreateSecret(secrets secretService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		baseURL := r.FormValue("secret-form-base-url")
		if len(baseURL) == 0 {
			http.Error(w, "could not create secret: missing base URL", http.StatusBadRequest)
			return
		}

		ttl, err := time.ParseDuration(r.FormValue("secret-form-ttl"))
		if err != nil {
			http.Error(w, "could not create secret: invalid expiration time", http.StatusBadRequest)
			return
		}

		secret, err := secrets.Create(secret.Secret{
			Value:      r.FormValue("secret-form-value"),
			Passphrase: r.FormValue("secret-form-passphrase"),
			TTL:        ttl,
		})
		if err != nil {
			http.Error(w, "could not create secret error in service", http.StatusInternalServerError)
			return
		}

		response := secretCreateResponse{
			BaseURL:        baseURL,
			ID:             secret.ID,
			Passphrase:     secret.Passphrase,
			PassphraseHash: secret.PassphraseHash,
		}

		renderTemplate(w, http.StatusCreated, "secret-created", response)
	})
}

// HandlerGetSecret handles requests containing a form to get a secret.
// This form will be used when a passphrase is not provided in the URL.
func HandlerGetSecret(secrets secretService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := r.FormValue("secret-result-form-id")
		if len(id) == 0 {
			http.Error(w, "could not get secret: missing ID", http.StatusBadRequest)
			return
		}
		passphrase := r.FormValue("secret-result-form-passphrase")
		if len(passphrase) == 0 {
			http.Error(w, "could not get secret: missing passphrase", http.StatusBadRequest)
			return
		}

		s, err := secrets.Get(id, passphrase)
		if err != nil {
			if errors.Is(err, secret.ErrSecretNotFound) {
				renderTemplate(w, http.StatusNotFound, "secret-not-found", nil)
				return
			}
			if errors.Is(err, secret.ErrInvalidPassphrase) {
				renderTemplate(w, http.StatusUnauthorized, "secret-get", secretGetResponse{ID: id})
				return
			}
			http.Error(w, "could not get secret: error in service", http.StatusInternalServerError)
			return
		}

		response := secretGetResponse{
			ID:             s.ID,
			PassphraseHash: passphrase,
			Value:          s.Value,
		}

		renderTemplate(w, http.StatusOK, "secret-get", response)
	})
}

// renderTemplate renders a template with the given data.
func renderTemplate(w http.ResponseWriter, statusCode int, tmpl string, data any) {
	t, err := template.ParseFS(TemplateFS, "templates/base.html", "templates/"+tmpl+".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	if err := t.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
}

// secretGetResponse is the response data for a get secret request.
type secretGetResponse struct {
	ID             string
	PassphraseHash string
	Value          string
}
