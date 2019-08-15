package api

import (
	"github.com/RedeployAB/redeploy-secrets/common/auth"
	mw "github.com/RedeployAB/redeploy-secrets/common/middleware"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

var apiVer = "v1"

// NewRouter returns a mux Router.
func NewRouter(ts auth.TokenStore, collection *mgo.Collection) *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/" + apiVer).Subrouter()
	// Routes.
	s.Handle("/secrets/{id}", readSecretHandler(collection)).Methods("GET")
	s.Handle("/secrets", createSecretHandler(collection)).Methods("POST")
	s.Handle("/secrets/{id}", updateSecretHandler(collection)).Methods("PUT")
	s.Handle("/secrets/{id}", deleteSecretHandler(collection)).Methods("DELETE")
	// All other routes.
	r.PathPrefix("/").HandlerFunc(notFoundHandler)
	// Attach middleware.
	amw := mw.Authentication{TokenStore: ts}
	r.Use(mw.Logger, amw.Authenticate)

	return r
}
