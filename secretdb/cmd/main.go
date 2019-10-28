package main

import (
	"log"

	"github.com/RedeployAB/burnit/common/auth"
	"github.com/RedeployAB/burnit/secretdb/app"
	"github.com/RedeployAB/burnit/secretdb/config"
	"github.com/RedeployAB/burnit/secretdb/db"
	"github.com/gorilla/mux"
)

func main() {
	// Setup configuration.
	conf := config.Configure()
	ts := auth.NewMemoryTokenStore()
	ts.Set(conf.Server.DBAPIKey, "app")
	// Connect to database.
	log.Printf("connecting to db server: %s...\n", conf.Database.Address)
	connection, err := db.Connect(conf.Database)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// Setup repositories.
	r := mux.NewRouter()
	srv := app.NewServer(app.ServerOptions{
		Config:     config,
		Router:     r,
		Connection: connection,
		TokenStore: ts,
	})
	// Listen and serve.
	srv.Start()
}
