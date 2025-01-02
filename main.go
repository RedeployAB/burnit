package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/RedeployAB/burnit/internal/config"
	"github.com/RedeployAB/burnit/internal/log"
	"github.com/RedeployAB/burnit/internal/server"
	"github.com/RedeployAB/burnit/internal/session"
	"github.com/RedeployAB/burnit/internal/ui"
	"github.com/RedeployAB/burnit/internal/version"
)

func main() {
	log := log.New()
	if err := run(log); err != nil {
		log.Error("Server error.", "error", err)
		os.Exit(1)
	}
}

// run the application.
func run(log log.Logger) error {
	log.Info("Starting service.", "version", version.Version())
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("could not load configuration: %w", err)
	}
	log.Info("Configuration loaded.", "config", cfg)

	services, err := config.SetupServices(cfg.Services)
	if err != nil {
		return fmt.Errorf("could not setup services: %w", err)
	}

	var ui ui.UI
	var sessions session.Service
	if cfg.Server.BackendOnly == nil || !*cfg.Server.BackendOnly {
		ui, sessions, err = config.SetupUI(cfg.UI)
		if err != nil {
			return fmt.Errorf("could not setup ui: %w", err)
		}
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
			TTL:             cfg.Server.RateLimiter.TTL,
			CleanupInterval: cfg.Server.RateLimiter.CleanupInterval,
		}),
		server.WithUI(ui),
		server.WithSessionService(sessions),
	)
	if err != nil {
		return fmt.Errorf("could setup server: %w", err)
	}

	if err := srv.Start(); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
