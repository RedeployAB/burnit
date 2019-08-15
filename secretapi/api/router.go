package api

import (
	mw "github.com/RedeployAB/redeploy-secrets/common/middleware"
	"github.com/RedeployAB/redeploy-secrets/secretapi/config"
	"github.com/gorilla/mux"
)

var genAPIVer = "v1"
var dbAPIVer = "v1"

// NewRouter returns a mux Router.
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Generator routes.
	g := r.PathPrefix("/api/" + genAPIVer).Subrouter()
	g.Handle("/generate", generateSecretHandler()).Methods("GET")

	// DB routes.
	d := r.PathPrefix("/api/" + dbAPIVer).Subrouter()
	d.Handle("/secrets/{id}", readSecretHandler()).Methods("GET")
	d.Handle("/secrets", createSecretHandler()).Methods("POST")
	// Init middleware for all db routes.
	amw := mw.AuthHeader{Token: config.Config.DBAPIKey}
	hmw := mw.HeaderStrip{Exceptions: []string{"X-Passphrase"}}
	d.Use(mw.Logger, hmw.Strip, amw.AddAuthHeader)
	// All other routes.
	r.PathPrefix("/").HandlerFunc(notFoundHandler)
	// Attach logger.
	r.Use(mw.Logger)

	return r
}
