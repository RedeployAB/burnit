package api

import (
	"github.com/RedeployAB/redeploy-secrets/common/auth"
	mw "github.com/RedeployAB/redeploy-secrets/common/middleware"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var apiVer = "v1"

// Router represents a mux.Router with a mongo client.
type Router struct {
	*mux.Router
	client *mongo.Client
}

// NewRouter returns a mux Router.
func NewRouter(ts auth.TokenStore, client *mongo.Client) *Router {
	r := &Router{mux.NewRouter(), client}
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
