package main

import (
	"github.com/RedeployAB/burnit/config"
	"github.com/RedeployAB/burnit/server"
)

func main() {
	log := server.NewDefaultLogger()

	cfg, err := config.New()
	if err != nil {
		log.Error("Configuration error.", "error", err)
	}

	services, err := config.SetupServices(cfg.Services)
	if err != nil {
		log.Error("Services setup error.", "error", err)
	}

	srv := server.New(server.WithOptions(server.Options{
		Host:    cfg.Server.Host,
		Port:    cfg.Server.Port,
		Logger:  log,
		Secrets: services.Secrets,
		TLS: server.TLSConfig{
			Certificate: cfg.Server.TLS.CertFile,
			Key:         cfg.Server.TLS.KeyFile,
		},
	}))

	if err := srv.Start(); err != nil {
		log.Error("Server error.", "error", err)
	}
}
