package server

import (
	"net/http"
	"strconv"

	"github.com/RedeployAB/burnit/burnit/secrets"
	"github.com/gorilla/mux"
)

// notFound handles all non used routes.
func (s *server) notFound(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotFound, "")
}

// getSecret reads a secret fron the database by ID.
func (s *server) getSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sec, err := s.secrets.Get(vars["id"], r.Header.Get("passphrase"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, "")
			return
		}
		if sec == nil {
			writeError(w, http.StatusNotFound, "")
			return
		}

		if len(sec.Value) == 0 {
			writeError(w, http.StatusUnauthorized, "")
			return
		}

		_, err = s.secrets.Delete(sec.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(newSecret(sec).JSON())
	})
}

// createSecret inserts a secret into the database.
func (s *server) createSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sec, err := secrets.NewFromJSON(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		sec, err = s.secrets.Create(sec)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Location", "/secrets/"+sec.ID)
		w.WriteHeader(http.StatusCreated)
		w.Write(newSecret(sec).JSON())
	})
}

// deleteSecret deletes a secret from the database.
func (s *server) deleteSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		deleted, err := s.secrets.Delete(vars["id"])
		if err != nil {
			writeError(w, http.StatusInternalServerError, "")
			return
		}
		if deleted == 0 {
			writeError(w, http.StatusNotFound, "")
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
		w.Write(newSecret(secret).JSON())
	})
}
