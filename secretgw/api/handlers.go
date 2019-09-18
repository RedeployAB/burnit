package api

import (
	"encoding/json"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/common/httperror"
	"github.com/RedeployAB/redeploy-secrets/secretgw/config"
	"github.com/RedeployAB/redeploy-secrets/secretgw/internal"
	"github.com/gorilla/mux"
)

var (
	genAPIPath = "/api/v1/generate"
	dbAPIPath  = "/api/v1/secrets"
)

// Define clients. These should be initialized and
// ready at app start.
var genClient *internal.GeneratorClient
var dbClient *internal.DBClient

func init() {
	// Inititalize clients.
	genClient = &internal.GeneratorClient{BaseURL: config.Config.GeneratorBaseURL, Path: genAPIPath}
	dbClient = &internal.DBClient{BaseURL: config.Config.DBBaseURL, Path: dbAPIPath}
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
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		q := r.URL.RawQuery
		if q != "" {
			genClient.Path += "?" + q
		}

		s, err := genClient.Fetch()
		if err != nil {
			httperror.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&s); err != nil {
			panic(err)
		}
	})
}

// getSecretHandler makes calls to the secretdb service to
// get a secret.
func getSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)

		s, err := dbClient.Do("GET", vars["id"], r.Header, nil)
		if err != nil {
			status := internal.HandleHTTPError(err)
			w.WriteHeader(status)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&s); err != nil {
			panic(err)
		}
	})
}

// createSecret makes calls to secretdb to
// create a secret.
func createSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		s, err := dbClient.Do("POST", "", r.Header, r.Body)
		if err != nil {
			status := internal.HandleHTTPError(err)
			w.WriteHeader(status)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(&s); err != nil {
			panic(err)
		}
	})
}
