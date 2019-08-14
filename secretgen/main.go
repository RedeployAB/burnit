package main

import (
	"log"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/secretgen/api"
	"github.com/RedeployAB/redeploy-secrets/secretgen/config"
	"github.com/RedeployAB/redeploy-secrets/secretgen/middleware"
	"github.com/gorilla/mux"
)

var apiVer = "v1"
var port = config.Config.Port

func main() {
	r := mux.NewRouter()
	// Routes.
	r.HandleFunc("/api/"+apiVer+"/secret", api.GenerateSecretHandler).Methods("GET")
	r.PathPrefix("/").HandlerFunc(api.NotFoundHandler)
	r.Use(middleware.Logger)
	// Start server.
	log.Printf("server listening on: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
