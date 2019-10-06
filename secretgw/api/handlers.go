package api

import (
	"encoding/json"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/common/httperror"
	"github.com/RedeployAB/redeploy-secrets/secretgw/client"
	"github.com/RedeployAB/redeploy-secrets/secretgw/config"
	"github.com/gorilla/mux"
)

var (
	genAPIPath = "/api/v1/generate"
	dbAPIPath  = "/api/v1/secrets"
)

// Define clients. These should be initialized and
// ready at app start.
var genSvc *client.APIClient
var dbSvc *client.APIClient

func init() {
	// Inititalize clients.
	genSvc = &client.APIClient{BaseURL: config.Config.GeneratorBaseURL, Path: config.Config.GeneratorServicePath}
	dbSvc = &client.APIClient{BaseURL: config.Config.DBBaseURL, Path: config.Config.DBServicePath}
}

// notFoundHandler handles all non used routes.
func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

// generateSecret makes calls to the secretgen service
// to generate a secret.
func generateSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := genSvc.Request(client.RequestOptions{
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
func getSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)

		res, err := dbSvc.Request(client.RequestOptions{
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
func createSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		res, err := dbSvc.Request(client.RequestOptions{
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
