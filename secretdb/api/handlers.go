package api

import (
	"encoding/json"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/common/httperror"
	"github.com/RedeployAB/redeploy-secrets/secretdb/db"
	"github.com/RedeployAB/redeploy-secrets/secretdb/secret"
	"github.com/gorilla/mux"
)

// notFoundHandler handles all non used routes.
func (rt *Router) notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

// getSecretHandler reads a secret fron the database by ID.
func (rt *Router) getSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)

		repo := db.SecretRepository{DB: rt.DB, Config: rt.Config.Server}
		s, err := repo.Find(vars["id"])
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if s == nil {
			httperror.Error(w, "not found", http.StatusNotFound)
			return
		}

		if !s.VerifyPassphrase(r.Header.Get("X-Passphrase")) {
			httperror.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		_, err = repo.Delete(vars["id"])
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
		}

		sr := secretResponseBody{
			Data: secretResponse{
				ID:        s.ID,
				Secret:    s.Secret,
				CreatedAt: s.CreatedAt,
				ExpiresAt: s.ExpiresAt,
			},
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&sr); err != nil {
			panic(err)
		}

	})
}

// createSecretHandler inserts a secret into the database.
func (rt *Router) createSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := secret.NewSecret(r.Body)
		if err != nil {
			httperror.Error(w, "malformed JSON", http.StatusBadRequest)
			return
		}

		repo := db.SecretRepository{DB: rt.DB, Config: rt.Config.Server}
		s, err = repo.Insert(s)
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
		}

		sr := secretResponseBody{
			Data: secretResponse{
				ID:        s.ID,
				CreatedAt: s.CreatedAt,
				ExpiresAt: s.ExpiresAt,
			},
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(&sr); err != nil {
			panic(err)
		}
	})
}

// updateSecretHandler handler updates a secret in the database.
func (rt *Router) updateSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		httperror.Error(w, "not implemented", http.StatusNotImplemented)
	})
}

// deleteSecretHandler deletes a secret from the database.
func (rt *Router) deleteSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		repo := db.SecretRepository{DB: rt.DB, Config: rt.Config.Server}
		res, err := repo.Delete(vars["id"])
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if res == 0 || res == -1 {
			httperror.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	})
}

// deleteExpiredSecrets will delete all secrets where ExpiresAt has
// passed.
func (rt *Router) deleteExpiredSecrets() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Contente-Type", "application/json; charset=UTF-8")

		repo := db.SecretRepository{DB: rt.DB, Config: rt.Config.Server}
		_, err := repo.DeleteExpired()
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
