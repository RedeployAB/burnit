package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/RedeployAB/burnit/internal/config"
	"github.com/RedeployAB/burnit/internal/log"
	"github.com/RedeployAB/burnit/internal/server"
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
	flags, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			// Return nil and show help.
			return nil
		}
		return fmt.Errorf("could not parse flags: %w", err)
	}

	log.Info("Starting service.", "version", version.Version())
	cfg, err := config.New(config.WithFlags(flags))
	if err != nil {
		return fmt.Errorf("could not load configuration: %w", err)
	}
	log.Info("Configuration loaded.", "config", cfg)

	services, err := config.Setup(cfg)
	if err != nil {
		return fmt.Errorf("could not setup services: %w", err)
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
		server.WithUI(services.UI),
	)
	if err != nil {
		return fmt.Errorf("could setup server: %w", err)
	}

	if err := srv.Start(); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
