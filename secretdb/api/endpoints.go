package api

import (
	"encoding/json"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/secretdb/internal"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

// NotFoundHandler handles all non used routes.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

// ReadSecretHandler reads a secret fron the database by ID.
func ReadSecretHandler(collection *mgo.Collection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)

		s, err := Find(vars["id"], collection)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Handle if passphrase is set on the secret.
		if len(s.Passphrase) > 0 && !internal.CompareHash(s.Passphrase, r.Header.Get("x-passphrase")) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		Delete(vars["id"], collection)
		sr := SecretResponseBody{
			Data: secretData{
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

// CreateSecretHandler inserts a secret into the database.
func CreateSecretHandler(collection *mgo.Collection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var sb SecretBody
		if err := json.NewDecoder(r.Body).Decode(&sb); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(&ErrorResponseBody{Error: "malformed request"}); err != nil {
				panic(err)
			}
			return
		}

		s, err := Insert(sb, collection)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		sr := SecretResponseBody{
			Data: secretData{
				ID:        s.ID.Hex(),
				CreatedAt: s.CreatedAt,
				ExpiresAt: s.ExpiresAt,
			},
		}
		if err := json.NewEncoder(w).Encode(&sr); err != nil {
			panic(err)
		}
	})
}

// UpdateSecretHandler handler updates a secret in the database.
func UpdateSecretHandler(collection *mgo.Collection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotImplemented)
		if err := json.NewEncoder(w).Encode(&ResponseBody{Data: "not implemented"}); err != nil {
			panic(err)
		}
	})
}

// DeleteSecretHandler deletes a secret from the database.
func DeleteSecretHandler(collection *mgo.Collection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		vars := mux.Vars(r)

		err := Delete(vars["id"], collection)
		if err != nil {
			if err.Error() == "not found" || err.Error() == "not valid ObjectId" {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
