package server

import (
	"encoding/json"
	"net/http"

	"github.com/RedeployAB/burnit/burnit/secret"
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
		sec, err := s.secretService.Get(vars["id"], r.Header.Get("passphrase"))
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

		_, err = s.secretService.Delete(sec.ID)
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
		sec, err := secret.NewFromJSON(r.Body)
		if err != nil {
			httperror.Write(w, http.StatusBadRequest, "", err.Error())
			return
		}
		defer r.Body.Close()

		sec, err = s.secretService.Create(sec)
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

// updateSecret handler updates a secret in the database.
func (s *server) updateSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httperror.Write(w, http.StatusNotImplemented, "", "")
	})
}

// deleteSecret deletes a secret from the database.
func (s *server) deleteSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		deleted, err := s.secretService.Delete(vars["id"])
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
