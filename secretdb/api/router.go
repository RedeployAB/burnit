package api

import (
	"github.com/RedeployAB/redeploy-secrets/common/auth"
	mw "github.com/RedeployAB/redeploy-secrets/common/middleware"
	"github.com/RedeployAB/redeploy-secrets/secretapi/api"
	"github.com/RedeployAB/redeploy-secrets/secretdb/config"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

var apiVer = "v1"

// NewRouter returns a mux Router.
func NewRouter(config config.Server, collection *mgo.Collection) *mux.Router {
	// Setup TokenStore.
	ts := auth.NewMemoryTokenStore()
	ts.Set(config.DBAPIKey, "app")

	r := mux.NewRouter()
	s := r.PathPrefix("/api/" + apiVer).Subrouter()
	// Routes.
	s.Handle("/secrets/{id}", ReadSecretHandler(collection)).Methods("GET")
	s.Handle("/secrets", CreateSecretHandler(collection)).Methods("POST")
	s.Handle("/secrets", CreateSecretHandler(collection)).Methods("POST")
	s.Handle("/secrets/{id}", UpdateSecretHandler(collection)).Methods("PUT")
	s.Handle("/secrets/{id}", DeleteSecretHandler(collection)).Methods("DELETE")
	// All other routes.
	r.PathPrefix("/").HandlerFunc(api.NotFoundHandler)
	// Attach middleware.
	amw := mw.AuthenticationMiddleware{TokenStore: ts}
	r.Use(mw.Logger, amw.Authenticate)

	return r
}
