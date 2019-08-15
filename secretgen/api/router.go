package api

import (
	"github.com/RedeployAB/redeploy-secrets/common/middleware"
	"github.com/RedeployAB/redeploy-secrets/secretapi/api"
	"github.com/RedeployAB/redeploy-secrets/secretgen/config"
	"github.com/gorilla/mux"
)

var apiVer = "v1"

// NewRouter returns a mux Router.
func NewRouter(config config.Server) *mux.Router {
	r := mux.NewRouter()
	// Routes.
	r.HandleFunc("/api/"+apiVer+"/generate", generateSecretHandler).Methods("GET")
	r.PathPrefix("/").HandlerFunc(api.NotFoundHandler)
	r.Use(middleware.Logger)

	return r
}
