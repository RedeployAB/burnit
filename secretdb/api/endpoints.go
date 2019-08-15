package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RedeployAB/redeploy-secrets/common/httperror"
	"github.com/RedeployAB/redeploy-secrets/secretdb/db"
	"github.com/RedeployAB/redeploy-secrets/secretdb/internal"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

// NotFoundHandler handles all non used routes.
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

// ReadSecretHandler reads a secret fron the database by ID.
func readSecretHandler(collection *mgo.Collection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)

		sm, err := db.Find(vars["id"], collection)
		if err != nil {
			httperror.Error(w, "not found", http.StatusNotFound)
			return
		}
		// Handle if passphrase is set on the secret.
		if len(sm.Passphrase) > 0 && !internal.CompareHash(sm.Passphrase, r.Header.Get("x-passphrase")) {
			httperror.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		db.Delete(vars["id"], collection)
		sr := secretResponseBody{
			Data: secretData{
				Secret:    sm.Secret,
				CreatedAt: sm.CreatedAt,
				ExpiresAt: sm.ExpiresAt,
			},
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&sr); err != nil {
			panic(err)
		}
	})
}

// CreateSecretHandler inserts a secret into the database.
func createSecretHandler(collection *mgo.Collection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var s db.Secret
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			httperror.Error(w, "malformed JSON", http.StatusBadRequest)
			return
		}

		sm, err := db.Insert(s, collection)
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		sr := secretResponseBody{
			Data: secretData{ID: sm.ID.Hex(), CreatedAt: sm.CreatedAt, ExpiresAt: sm.ExpiresAt},
		}
		if err := json.NewEncoder(w).Encode(&sr); err != nil {
			panic(err)
		}
	})
}

// UpdateSecretHandler handler updates a secret in the database.
func updateSecretHandler(collection *mgo.Collection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		httperror.Error(w, "not implemented", http.StatusNotImplemented)
	})
}

// DeleteSecretHandler deletes a secret from the database.
func deleteSecretHandler(collection *mgo.Collection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		vars := mux.Vars(r)

		err := db.Delete(vars["id"], collection)
		if err != nil {
			if err.Error() == "not found" || err.Error() == "not valid ObjectId" {
				httperror.Error(w, "not found", http.StatusNotFound)
			} else {
				httperror.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// SecretResponseBody represents a secret type response.
type secretResponseBody struct {
	Data secretData `json:"data"`
}

type secretData struct {
	ID        string    `json:"id,omitempty"`
	Secret    string    `json:"secret,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
