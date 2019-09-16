package api

import (
	"github.com/RedeployAB/redeploy-secrets/common/auth"
	mw "github.com/RedeployAB/redeploy-secrets/common/middleware"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var apiVer = "v1"

// NewRouter returns a mux Router.
func NewRouter(ts auth.TokenStore, client *mongo.Client) *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/" + apiVer).Subrouter()
	// Routes.
	s.Handle("/secrets/{id}", getSecret(client)).Methods("GET")
	s.Handle("/secrets", createSecret(client)).Methods("POST")
	s.Handle("/secrets/{id}", updateSecret(client)).Methods("PUT")
	s.Handle("/secrets/{id}", deleteSecret(client)).Methods("DELETE")
	// All other routes.
	r.PathPrefix("/").HandlerFunc(notFound)
	// Attach middleware.
	amw := mw.Authentication{TokenStore: ts}
	r.Use(mw.Logger, amw.Authenticate)

	return r
}
