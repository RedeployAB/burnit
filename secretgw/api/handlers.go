package api

import (
	"encoding/json"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/common/httperror"
	"github.com/RedeployAB/redeploy-secrets/secretgw/client"
	"github.com/gorilla/mux"
)

// notFoundHandler handles all non used routes.
func (rt *Router) notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

// generateSecret makes calls to the secretgen service
// to generate a secret.
func (rt *Router) generateSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := rt.GeneratorServiceClient.Request(client.RequestOptions{
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
func (rt *Router) getSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)

		res, err := rt.DBServiceClient.Request(client.RequestOptions{
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
func (rt *Router) createSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		res, err := rt.DBServiceClient.Request(client.RequestOptions{
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
