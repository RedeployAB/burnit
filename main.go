package main

import (
	"os"

	"github.com/RedeployAB/burnit/config"
	"github.com/RedeployAB/burnit/server"
)

func main() {
	log := server.NewDefaultLogger()

	cfg, err := config.New()
	if err != nil {
		log.Error("Configuration error.", "error", err)
		os.Exit(1)
	}

	services, err := config.SetupServices(cfg.Services)
	if err != nil {
		log.Error("Services setup error.", "error", err)
		os.Exit(1)
	}

	srv, err := server.New(services.Secrets, server.WithOptions(server.Options{
		Host:   cfg.Server.Host,
		Port:   cfg.Server.Port,
		Logger: log,
		TLS: server.TLSConfig{
			Certificate: cfg.Server.TLS.CertFile,
			Key:         cfg.Server.TLS.KeyFile,
		},
		CORS: server.CORS{
			Origin: cfg.Server.CORS.Origin,
		},
	}))
	if err != nil {
		log.Error("Server setup error.", "error", err)
	}

	if err := srv.Start(); err != nil {
		log.Error("Server error.", "error", err)
	}
}
