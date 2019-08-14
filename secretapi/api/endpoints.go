package api

import (
	"encoding/json"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/secretapi/config"
	"github.com/RedeployAB/redeploy-secrets/secretapi/internal"
	"github.com/gorilla/mux"
)

var (
	genAPIPath = "/api/v1/secret"
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

// NotFoundHandler handles all non used routes.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

// GenerateSecretHandler makes calls to the secretgen service
// to generate a secret.
func GenerateSecretHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		q := r.URL.RawQuery
		if q != "" {
			genClient.Path += "?" + q
		}

		s, err := genClient.Fetch()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&s); err != nil {
			panic(err)
		}

	})
}

// ReadSecretHandler makes calls to the secretdb service to
// get a secret.
func ReadSecretHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)
		header := handleHeaders(r.Header)

		s, err := dbClient.Do("GET", vars["id"], header, nil)
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

// CreateSecretHandler makes calls to secretdb to
// create a secret.
func CreateSecretHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		s, err := dbClient.Do("POST", "", nil, r.Body)
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
