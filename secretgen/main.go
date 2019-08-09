package main

import (
	"log"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/secretgen/api"
	"github.com/RedeployAB/redeploy-secrets/secretgen/config"
	"github.com/gorilla/mux"
)

var conf config.Config
var apiVer = "v1"

func init() {
	conf = config.Configure()
	log.Printf("Server listening on: %s\n", conf.Port)
}

func main() {
	r := mux.NewRouter()
	// Routes.
	r.HandleFunc("/api/"+apiVer+"/secret", api.GenerateSecretHandler).Methods("GET")
	r.PathPrefix("/").HandlerFunc(api.NotFoundHandler)
	// Start server.
	log.Fatal(http.ListenAndServe(":"+conf.Port, r))
}
