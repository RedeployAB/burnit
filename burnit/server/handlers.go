package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RedeployAB/burnit/burnit/secrets"
	"github.com/RedeployAB/burnit/common/httperror"
	"github.com/gorilla/mux"
)

// notFound handles all non used routes.
func (s *server) notFound(w http.ResponseWriter, r *http.Request) {
	httperror.Write(w, http.StatusNotFound, "", "")
}

// getSecret reads a secret fron the database by ID.
func (s *server) getSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sec, err := s.secrets.Get(vars["id"], r.Header.Get("passphrase"))
		if err != nil {
			httperror.Write(w, http.StatusInternalServerError, "", "")
			return
		}

		if sec == nil {
			httperror.Write(w, http.StatusNotFound, "", "")
			return
		}

		if len(sec.Value) == 0 {
			httperror.Write(w, http.StatusUnauthorized, "", "")
			return
		}

		_, err = s.secrets.Delete(sec.ID)
		if err != nil {
			httperror.Write(w, http.StatusInternalServerError, "", "")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(newSecretResponse(sec)); err != nil {
			panic(err)
		}
	})
}

// createSecret inserts a secret into the database.
func (s *server) createSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sec, err := secrets.NewFromJSON(r.Body)
		if err != nil {
			httperror.Write(w, http.StatusBadRequest, "", err.Error())
			return
		}
		defer r.Body.Close()

		sec, err = s.secrets.Create(sec)
		if err != nil {
			httperror.Write(w, http.StatusInternalServerError, "", "")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Location", "/secrets/"+sec.ID)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(newSecretResponse(sec)); err != nil {
			panic(err)
		}
	})
}

// deleteSecret deletes a secret from the database.
func (s *server) deleteSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		deleted, err := s.secrets.Delete(vars["id"])
		if err != nil {
			httperror.Write(w, http.StatusInternalServerError, "", "")
			return
		}
		if deleted == 0 {
			httperror.Write(w, http.StatusNotFound, "", "")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	})
}

// generateSecret generates a secret.
func (s *server) generateSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		length := query.Get("length")
		specialchars := query.Get("specialchars")

		l, err := strconv.Atoi(length)
		if err != nil {
			l = 16
		}
		sc, err := strconv.ParseBool(specialchars)
		if err != nil {
			sc = false
		}

		secret := s.secrets.Generate(l, sc)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(newSecretResponse(secret)); err != nil {
			panic(err)
		}
	})
}
