package server

import (
	"encoding/json"
	"net/http"

	"github.com/RedeployAB/burnit/burnitdb/internal/dto"
	"github.com/RedeployAB/burnit/burnitdb/internal/mappers"
	"github.com/RedeployAB/burnit/common/httperror"
	"github.com/gorilla/mux"
)

// notFound handles all non used routes.
func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	httperror.Error(w, http.StatusNotFound)
}

// getSecret reads a secret fron the database by ID.
func (s *Server) getSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		secretModel, err := s.repository.Find(vars["id"])
		if err != nil {
			httperror.Error(w, http.StatusInternalServerError)
			return
		}

		if secretModel == nil {
			httperror.Error(w, http.StatusNotFound)
			return
		}

		secretDTO := mappers.Secret{}.ToDTO(secretModel)
		passphrase := r.Header.Get("X-Passphrase")
		if len(passphrase) > 0 && !s.compareHash(secretDTO.Passphrase, passphrase) {
			httperror.Error(w, http.StatusUnauthorized)
			return
		}

		_, err = s.repository.Delete(secretDTO.ID)
		if err != nil {
			httperror.Error(w, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response(secretDTO)); err != nil {
			panic(err)
		}
	})
}

// createSecret inserts a secret into the database.
func (s *Server) createSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretDTO, err := dto.NewSecret(r.Body)
		if err != nil {
			httperror.Error(w, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		secretModel, err := s.repository.Insert(mappers.Secret{}.ToPersistance(secretDTO))
		if err != nil {
			httperror.Error(w, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response(mappers.Secret{}.ToDTO(secretModel))); err != nil {
			panic(err)
		}
	})
}

// updateSecret handler updates a secret in the database.
func (s *Server) updateSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		httperror.Error(w, http.StatusNotImplemented)
	})
}

// deleteSecret deletes a secret from the database.
func (s *Server) deleteSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		res, err := s.repository.Delete(vars["id"])
		if err != nil {
			httperror.Error(w, http.StatusInternalServerError)
			return
		}
		if res == 0 {
			httperror.Error(w, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	})
}
