package server

import (
	"html/template"
	"net/http"

	"github.com/RedeployAB/burnit/internal/secret"
)

// uiCreateSecret handles requests to create a secret for the UI.
func (s server) uiCreateSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "secret-create", nil)
	})
}

// uiGetSecret handles requests to get a secret for the UI.
func (s server) uiGetSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, passhrase, err := extractIDAndPassphrase("/ui/secrets/", r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Add handling of the various scenarios of provided values.

		response := secretGetResponse{
			IDParameter:             id,
			PassphraseHashParameter: passhrase,
		}

		renderTemplate(w, "secret-get", response)
	})
}

// handlertCreateSecret handles requests to create a secret.
func (s server) handlerCreateSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		baseURL := r.FormValue("secret-create-base-url")
		if len(baseURL) == 0 {
			http.Error(w, "could not create secret: missing base URL", http.StatusBadRequest)
			return
		}

		secret, err := s.secrets.Create(secret.Secret{
			Value:      r.FormValue("value"),
			Passphrase: r.FormValue("passphrase"),
		})
		if err != nil {
			http.Error(w, "could not create secret: error in service", http.StatusInternalServerError)
			return
		}

		response := secretCreatedResponse{
			BaseURL:        r.FormValue("secret-create-base-url"),
			ID:             secret.ID,
			Passphrase:     secret.Passphrase,
			PassphraseHash: secret.PassphraseHash,
		}
		renderTemplate(w, "secret-created", response)
	})
}

// renderTemplate renders a template with the given data.
func renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	t, err := template.ParseFiles("templates/base.html", "templates/"+tmpl+".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// secretCreatedResponse is the response data for a created secret.
// Used for rendering the secret-created template.
type secretCreatedResponse struct {
	BaseURL        string
	ID             string
	Passphrase     string
	PassphraseHash string
}

// secretGetResponse is the response data for a fetched secret.
// Used for rendering the secret-get template.
type secretGetResponse struct {
	IDParameter             string
	PassphraseHashParameter string
}
