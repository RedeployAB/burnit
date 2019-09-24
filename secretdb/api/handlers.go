package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RedeployAB/redeploy-secrets/common/httperror"
	"github.com/RedeployAB/redeploy-secrets/common/security"
	"github.com/RedeployAB/redeploy-secrets/secretdb/db"
	"github.com/RedeployAB/redeploy-secrets/secretdb/models"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

		res, err := db.Find(vars["id"], client)
		if err != nil {
			httperror.Error(w, "not found", http.StatusNotFound)
			return
		}
		// Handle if passphrase is set on the secret.
		if len(res.Passphrase) > 0 && !security.CompareHash(res.Passphrase, r.Header.Get("X-Passphrase")) {
			httperror.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		_, err = db.Delete(vars["id"], client)
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
		}

		rb := responseBody{
			Data: data{
				ID:        res.ID,
				Secret:    res.Secret,
				CreatedAt: res.CreatedAt,
				ExpiresAt: res.ExpiresAt,
			},
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&rb); err != nil {
			panic(err)
		}

	})
}

// createSecretHandler inserts a secret into the database.
func createSecret(client *mongo.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var s models.Secret
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			httperror.Error(w, "malformed JSON", http.StatusBadRequest)
			return
		}

		res, err := db.Insert(s, client)
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusCreated)

		rb := &responseBody{
			Data: data{ID: res.ID, CreatedAt: res.CreatedAt, ExpiresAt: res.ExpiresAt},
		}
		if err := json.NewEncoder(w).Encode(rb); err != nil {
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
		res, err := db.Delete(vars["id"], client)
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if res == 0 {
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

		err := db.DeleteExpired(client)
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	})
}

// responseBody represents a secret type response body.
type responseBody struct {
	Data data `json:"data"`
}

// data represents the data part of the response body.
type data struct {
	ID        primitive.ObjectID `json:"id,omitempty"`
	Secret    string             `json:"secret,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty"`
	ExpiresAt time.Time          `json:"expiresAt,omitempty"`
}
