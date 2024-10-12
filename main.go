package main

import (
	"os"
	"strconv"

	"github.com/RedeployAB/burnit/config"
	"github.com/RedeployAB/burnit/log"
	"github.com/RedeployAB/burnit/server"
)

func main() {
	log := log.New()

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

	srv, err := server.New(
		services.Secrets,
		server.WithAddress(cfg.Server.Host+":"+strconv.Itoa(cfg.Server.Port)),
		server.WithLogger(log),
		server.WithTLS(server.TLSConfig{Certificate: cfg.Server.TLS.CertFile, Key: cfg.Server.TLS.KeyFile}),
		server.WithCORS(server.CORS{Origin: cfg.Server.CORS.Origin}),
		server.WithRateLimiter(server.RateLimiter{
			Rate:            float64(cfg.Server.RateLimiter.Rate),
			Burst:           cfg.Server.RateLimiter.Burst,
			CleanupInterval: cfg.Server.RateLimiter.CleanupInterval,
			TTL:             cfg.Server.RateLimiter.TTL,
		}),
	)
	if err != nil {
		log.Error("Server setup error.", "error", err)
		os.Exit(1)
	}

	if err := srv.Start(); err != nil {
		log.Error("Server error.", "error", err)
		os.Exit(1)
	}
}

// run the application.
func run(log log.Logger) error {
	return nil
}
