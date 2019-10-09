package api

import (
	"github.com/RedeployAB/redeploy-secrets/common/auth"
	mw "github.com/RedeployAB/redeploy-secrets/common/middleware"
	"github.com/RedeployAB/redeploy-secrets/secretdb/config"
	"github.com/RedeployAB/redeploy-secrets/secretdb/db"
	"github.com/gorilla/mux"
)

var apiVer = "v0"

// Router represents a mux.Router with a mongo client.
type Router struct {
	*mux.Router
	Config config.Configuration
	DB     *db.DB
}

// NewRouter returns a mux Router.
func NewRouter(config config.Configuration, ts auth.TokenStore, db *db.DB) *Router {
	r := &Router{mux.NewRouter(), config, db}
	s := r.PathPrefix("/api/" + apiVer).Subrouter()
	// Routes.
	s.Handle("/secrets/{id}", r.getSecret()).Methods("GET")
	s.Handle("/secrets", r.createSecret()).Methods("POST")
	s.Handle("/secrets/{id}", r.updateSecret()).Methods("PUT")
	s.Handle("/secrets/{id}", r.deleteSecret()).Methods("DELETE")
	s.Handle("/maintenance/cleanup", r.deleteExpiredSecrets()).Methods("DELETE")
	// All other routes.
	r.PathPrefix("/").HandlerFunc(r.notFound)
	// Attach middleware.
	amw := mw.Authentication{TokenStore: ts}
	r.Use(mw.Logger, amw.Authenticate)

	return r
}
