package api

import (
	"encoding/json"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/common/httperror"
	"github.com/RedeployAB/redeploy-secrets/secretdb/db"
	"github.com/RedeployAB/redeploy-secrets/secretdb/secret"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"
)

// notFoundHandler handles all non used routes.
func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

// getSecretHandler reads a secret fron the database by ID.
func getSecret(client *mongo.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)

		repo := db.SecretRepository{Client: client}
		res, err := repo.Find(vars["id"])
		if err != nil {
			if err.Error() == "mongo: no documents in result" {
				httperror.Error(w, "not found", http.StatusNotFound)
			} else {
				httperror.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}
		// Handle if passphrase is set on the secret.
		h, p := res.Passphrase, r.Header.Get("X-Passphrase")
		if !secret.VerifyPassphrase(h, p) {
			httperror.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		_, err = repo.Delete(vars["id"])
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
		}

		sr := secretResponseBody{
			Data: secretResponse{
				ID:        res.ID,
				Secret:    res.Secret,
				CreatedAt: res.CreatedAt,
				ExpiresAt: res.ExpiresAt,
			},
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&sr); err != nil {
			panic(err)
		}

	})
}

// createSecretHandler inserts a secret into the database.
func createSecret(client *mongo.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		//var s secret.Secret
		s, err := secret.NewSecret(r.Body)
		if err != nil {
			httperror.Error(w, "malformed JSON", http.StatusBadRequest)
			return
		}
		/* 		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			httperror.Error(w, "malformed JSON", http.StatusBadRequest)
			return
		} */
		repo := db.SecretRepository{Client: client}
		res, err := repo.Insert(s)
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusCreated)

		sr := secretResponseBody{
			Data: secretResponse{ID: res.ID, CreatedAt: res.CreatedAt, ExpiresAt: res.ExpiresAt},
		}
		if err := json.NewEncoder(w).Encode(&sr); err != nil {
			panic(err)
		}
	})
}

// updateSecretHandler handler updates a secret in the database.
func updateSecret(client *mongo.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		httperror.Error(w, "not implemented", http.StatusNotImplemented)
	})
}

// deleteSecretHandler deletes a secret from the database.
func deleteSecret(client *mongo.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		vars := mux.Vars(r)

		repo := db.SecretRepository{Client: client}
		res, err := repo.Delete(vars["id"])
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if res == 0 || res == -1 {
			httperror.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

// deleteExpiredSecrets will delete all secrets where ExpiresAt has
// passed.
func deleteExpiredSecrets(client *mongo.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Contente-Type", "application/json; charset=UTF-8")

		_, err := db.DeleteExpired(client)
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
