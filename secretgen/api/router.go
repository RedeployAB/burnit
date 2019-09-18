package api

import (
	"github.com/RedeployAB/redeploy-secrets/common/middleware"
	"github.com/RedeployAB/redeploy-secrets/secretgen/config"
	"github.com/gorilla/mux"
)

var apiVer = "v1"

// NewRouter returns a mux Router.
func NewRouter(config config.Server) *mux.Router {
	r := mux.NewRouter()
	// Routes.
	r.HandleFunc("/api/"+apiVer+"/generate", generateSecret).Methods("GET")
	r.PathPrefix("/").HandlerFunc(notFound)
	r.Use(middleware.Logger)

	return r
}
