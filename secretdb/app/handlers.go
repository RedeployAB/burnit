package app

import (
	"encoding/json"
	"net/http"

	"github.com/RedeployAB/burnit/common/httperror"
	"github.com/gorilla/mux"
)

// notFound handles all non used routes.
func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

// getSecret reads a secret fron the database by ID.
func (s *Server) getSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)

		secret, err := s.repository.Find(vars["id"])
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if secret == nil {
			httperror.Error(w, "not found", http.StatusNotFound)
			return
		}
		if !secret.VerifyPassphrase(r.Header.Get("X-Passphrase")) {
			httperror.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		_, err = s.repository.Delete(vars["id"])
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
		}

		sr := secretResponse{
			Data: secretData{
				ID:        secret.ID,
				Secret:    secret.Secret,
				CreatedAt: secret.CreatedAt,
				ExpiresAt: secret.ExpiresAt,
			},
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&sr); err != nil {
			panic(err)
		}
	})
}

// createSecret inserts a secret into the database.
/* func (s *Server) createSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret, err := dto.NewSecret(r.Body)
		if err != nil {
			httperror.Error(w, "malformed JSON", http.StatusBadRequest)
			return
		}

		secret, err = s.repository.Insert(secret)
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
		}

		sr := secretResponse{
			Data: secretData{
				ID:        secret.ID,
				CreatedAt: secret.CreatedAt,
				ExpiresAt: secret.ExpiresAt,
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
func (s *Server) updateSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		httperror.Error(w, "not implemented", http.StatusNotImplemented)
	})
}

// deleteSecretHandler deletes a secret from the database.
func (s *Server) deleteSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		res, err := s.repository.Delete(vars["id"])
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
} */

// deleteExpiredSecrets will delete all secrets where ExpiresAt has
// passed.
/* func (s *Server) deleteExpiredSecrets() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Contente-Type", "application/json; charset=UTF-8")

		_, err := s.repository.DeleteExpired()
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
} */
