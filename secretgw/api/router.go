package api

import (
	mw "github.com/RedeployAB/burnit/common/middleware"
	"github.com/RedeployAB/burnit/secretgw/client"
	"github.com/RedeployAB/burnit/secretgw/config"
	"github.com/gorilla/mux"
)

var genAPIVer = "v0"
var dbAPIVer = "v0"

// Router represents a mux,Router.
type Router struct {
	*mux.Router
	Config                 config.Configuration
	GeneratorServiceClient *client.APIClient
	DBServiceClient        *client.APIClient
}

// NewRouter returns a mux Router.
func NewRouter(config config.Configuration) *Router {
	r := &Router{Router: mux.NewRouter()}

	r.GeneratorServiceClient = &client.APIClient{BaseURL: config.GeneratorBaseURL, Path: config.GeneratorServicePath}
	r.DBServiceClient = &client.APIClient{BaseURL: config.DBBaseURL, Path: config.DBServicePath}
	// Generator routes.
	g := r.PathPrefix("/api/" + genAPIVer).Subrouter()
	g.Handle("/generate", r.generateSecret()).Methods("GET")

	// DB routes.
	d := r.PathPrefix("/api/" + dbAPIVer).Subrouter()
	d.Handle("/secrets/{id}", r.getSecret()).Methods("GET")
	d.Handle("/secrets", r.createSecret()).Methods("POST")
	// Init middleware for all db routes.
	amw := mw.AuthHeader{Token: config.DBAPIKey}
	hmw := mw.HeaderStrip{Exceptions: []string{"X-Passphrase"}}
	d.Use(hmw.Strip, amw.AddAuthHeader)
	// All other routes.
	r.PathPrefix("/").HandlerFunc(r.notFound)
	// Attach logger.
	r.Use(mw.Logger)

	return r
}
