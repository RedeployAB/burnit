package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RedeployAB/redeploy-secrets/common/httperror"
	"github.com/RedeployAB/redeploy-secrets/secretdb/db"
	"github.com/RedeployAB/redeploy-secrets/secretdb/internal"
	"github.com/RedeployAB/redeploy-secrets/secretdb/models"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// NotFoundHandler handles all non used routes.
func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

// ReadSecretHandler reads a secret fron the database by ID.
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
		if len(res.Passphrase) > 0 && !internal.CompareHash(res.Passphrase, r.Header.Get("X-Passphrase")) {
			httperror.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		_, err = db.Delete(vars["id"], client)
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
		}

		rb := ResponseBody{
			Data: Data{
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

// CreateSecretHandler inserts a secret into the database.
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

		rb := &ResponseBody{
			Data: Data{ID: res.ID, CreatedAt: res.CreatedAt, ExpiresAt: res.ExpiresAt},
		}
		if err := json.NewEncoder(w).Encode(rb); err != nil {
			panic(err)
		}
	})
}

// UpdateSecretHandler handler updates a secret in the database.
func updateSecret(client *mongo.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		httperror.Error(w, "not implemented", http.StatusNotImplemented)
	})
}

// DeleteSecretHandler deletes a secret from the database.
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

// ResponseBody represents a secret type response body.
type ResponseBody struct {
	Data Data `json:"data"`
}

// Data represents the data part of the response body.
type Data struct {
	ID        primitive.ObjectID `json:"id,omitempty"`
	Secret    string             `json:"secret,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty"`
	ExpiresAt time.Time          `json:"expires_at,omitempty"`
}

// RequestBody represents a secret request body.
type RequestBody struct {
	Secret     string `json:"secret"`
	Passphrase string `json:"passphrase,omitempty"`
}
