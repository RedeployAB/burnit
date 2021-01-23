package server

import (
	"encoding/json"
	"net/http"

	"github.com/RedeployAB/burnit/burnitgw/service/request"
	"github.com/RedeployAB/burnit/common/httperror"
	"github.com/gorilla/mux"
)

// notFoundHandler handles all non used routes.
func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	httperror.Write(w, http.StatusNotFound, "", "")
	return
}

// generateSecret makes calls to the burnitgen service
// to generate a secret.
func (s *Server) generateSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret, err := s.generatorService.Generate(r)
		if err != nil {
			writeError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&secret); err != nil {
			panic(err)
		}
	})
}

// getSecret makes calls to the burnitdb service to
// get a secret.
func (s *Server) getSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		secret, err := s.dbService.Get(r, vars)
		if err != nil {
			writeError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&secret); err != nil {
			panic(err)
		}
	})
}

// createSecret makes calls to burnitdb to
// create a secret.
func (s *Server) createSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret, err := s.dbService.Create(r)
		if err != nil {
			writeError(w, err)
			return
		}
		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Location", "/secrets/"+secret.ID)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(&secret); err != nil {
			panic(err)
		}
	})
}

func writeError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *request.Error:
		httperror.Write(w, e.StatusCode, e.Code, e.Message)
	default:
		httperror.Write(w, http.StatusInternalServerError, "", "")
	}
}
