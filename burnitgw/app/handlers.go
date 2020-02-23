package app

import (
	"encoding/json"
	"net/http"

	"github.com/RedeployAB/burnit/burnitgw/internal/request"
	"github.com/RedeployAB/burnit/common/httperror"
	"github.com/gorilla/mux"
)

// notFoundHandler handles all non used routes.
func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	httperror.Error(w, http.StatusNotFound)
}

// generateSecret makes calls to the burnitgen service
// to generate a secret.
func (s *Server) generateSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := s.generatorService.Request(request.Options{
			Method: request.GET,
			Query:  r.URL.RawQuery,
		})
		if err != nil {
			httperror.Error(w, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&res); err != nil {
			panic(err)
		}
	})
}

// getSecret makes calls to the burnitdb service to
// get a secret.
func (s *Server) getSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		res, err := s.dbService.Request(request.Options{
			Method: request.GET,
			Header: r.Header,
			Params: vars,
		})
		if err != nil {
			status := request.HandleHTTPError(err)
			w.WriteHeader(status)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&res); err != nil {
			panic(err)
		}
	})
}

// createSecret makes calls to burnitdb to
// create a secret.
func (s *Server) createSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := s.dbService.Request(request.Options{
			Method: request.POST,
			Header: r.Header,
			Body:   r.Body,
		})
		if err != nil {
			status := request.HandleHTTPError(err)
			w.WriteHeader(status)
			return
		}
		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(&res); err != nil {
			panic(err)
		}
	})
}
