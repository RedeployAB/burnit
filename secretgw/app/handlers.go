package app

import (
	"encoding/json"
	"net/http"

	"github.com/RedeployAB/burnit/common/httperror"
	"github.com/RedeployAB/burnit/secretgw/internal/client"
	"github.com/gorilla/mux"
)

// notFoundHandler handles all non used routes.
func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

// generateSecret makes calls to the secretgen service
// to generate a secret.
func (s *Server) generateSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := s.generatorService.Request(client.RequestOptions{
			Method: client.GET,
			Query:  r.URL.RawQuery,
		})
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&res); err != nil {
			panic(err)
		}
	})
}

// getSecret makes calls to the secretdb service to
// get a secret.
func (s *Server) getSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)

		res, err := s.dbService.Request(client.RequestOptions{
			Method: client.GET,
			Header: r.Header,
			Params: vars,
		})
		if err != nil {
			status := client.HandleHTTPError(err)
			w.WriteHeader(status)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&res); err != nil {
			panic(err)
		}
	})
}

// createSecret makes calls to secretdb to
// create a secret.
func (s *Server) createSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		res, err := s.dbService.Request(client.RequestOptions{
			Method: client.POST,
			Header: r.Header,
			Body:   r.Body,
		})
		if err != nil {
			status := client.HandleHTTPError(err)
			w.WriteHeader(status)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(&res); err != nil {
			panic(err)
		}
	})
}
