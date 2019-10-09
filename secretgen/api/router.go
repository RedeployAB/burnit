package api

import (
	"github.com/RedeployAB/redeploy-secrets/common/middleware"
	"github.com/gorilla/mux"
)

var apiVer = "v0"

// Router represents a mux.Router.
type Router struct {
	*mux.Router
}

// NewRouter returns a mux Router.
func NewRouter() *Router {
	r := &Router{mux.NewRouter()}
	// Routes.
	r.HandleFunc("/api/"+apiVer+"/generate", r.generateSecret).Methods("GET")
	r.PathPrefix("/").HandlerFunc(r.notFound)
	r.Use(middleware.Logger)

	return r
}
