package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RedeployAB/burnit/internal/log"
	"github.com/RedeployAB/burnit/internal/secret"
	"github.com/RedeployAB/burnit/internal/session"
	"github.com/RedeployAB/burnit/internal/ui"
)

// Defaults for server configuration.
const (
	defaultHost         = "0.0.0.0"
	defaultPort         = "3000"
	defaultReadTimeout  = 15 * time.Second
	defaultWriteTimeout = 15 * time.Second
	defaultIdleTimeout  = 30 * time.Second
)

// server holds an http.Server, a router and it's configured options.
type server struct {
	httpServer    *http.Server
	router        *router
	secrets       secret.Service
	ui            ui.UI
	sessions      session.Service
	tls           TLSConfig
	rateLimiter   RateLimiter
	log           log.Logger
	cors          CORS
	shutdownFuncs []func() error
	stopCh        chan os.Signal
	errCh         chan error
}

// TLSConfig holds the configuration for the server's TLS settings.
type TLSConfig struct {
	Certificate string
	Key         string
}

// isEmpty returns true if the TLSConfig is empty.
func (c TLSConfig) isEmpty() bool {
	return len(c.Certificate) == 0 && len(c.Key) == 0
}

// RateLimiter holds the configuration for the server's rate limiter settings.
type RateLimiter struct {
	Rate            float64
	Burst           int
	TTL             time.Duration
	CleanupInterval time.Duration
}

// isEmpty returns true if the RateLimiter is empty.
func (r RateLimiter) isEmpty() bool {
	return r.Rate == 0 && r.Burst == 0 && r.TTL == 0 && r.CleanupInterval == 0
}

// CORS holds the configuration for the server's CORS settings.
type CORS struct {
	Origin string
}

// isEmpty returns true if the CORS is empty.
func (c CORS) isEmpty() bool {
	return len(c.Origin) == 0
}

// Options holds the configuration for the server.
type Options struct {
	Router       *router
	Host         string
	Port         int
	TLS          TLSConfig
	Logger       log.Logger
	RateLimiter  RateLimiter
	CORS         CORS
	BackendOnly  bool
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Option is a function that configures the server.
type Option func(*server)

// New returns a new server.
func New(secrets secret.Service, options ...Option) (*server, error) {
	if secrets == nil {
		return nil, fmt.Errorf("secrets service is nil")
	}

	s := &server{
		httpServer: &http.Server{
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
			IdleTimeout:  defaultIdleTimeout,
		},
		secrets: secrets,
		stopCh:  make(chan os.Signal),
		errCh:   make(chan error),
	}
	for _, option := range options {
		option(s)
	}

	if s.router == nil {
		s.router = NewRouter()
	}
	if s.log == nil {
		s.log = log.New()
	}
	if len(s.httpServer.Addr) == 0 {
		s.httpServer.Addr = defaultHost + ":" + defaultPort
	}
	if s.httpServer.Handler == nil {
		s.httpServer.Handler = s.router
	}
	if s.sessions != nil {
		s.shutdownFuncs = append(s.shutdownFuncs, s.sessions.Close)
	}

	return s, nil
}

// Start the server.
func (s server) Start() error {
	defer func() {
		close(s.errCh)
		close(s.stopCh)
	}()

	go func() {
		for err := range s.secrets.Cleanup() {
			s.log.Error("Could not cleanup secrets", "error", err)
		}
	}()

	if s.sessions != nil {
		go func() {
			for err := range s.sessions.Cleanup() {
				s.log.Error("Could not cleanup sessions", "error", err)
			}
		}()
	}

	s.routes()
	go func() {
		if err := s.listenAndServe(); err != nil && err != http.ErrServerClosed {
			s.errCh <- err
		}
	}()

	go func() {
		s.stop()
	}()

	s.log.Info("Server started.", "address", s.httpServer.Addr)
	for {
		select {
		case err := <-s.errCh:
			return err
		case sig := <-s.stopCh:
			s.log.Info("Server stopped.", "reason", sig.String())
			return nil
		}
	}
}

// listenAndServe wraps around http.Server ListenAndServe and
// ListenAndServeTLS depending on TLS configuration.
func (s *server) listenAndServe() error {
	if !s.tls.isEmpty() {
		s.httpServer.TLSConfig = newTLSConfig()
		return s.httpServer.ListenAndServeTLS(s.tls.Certificate, s.tls.Key)
	}
	return s.httpServer.ListenAndServe()
}

// stop the server.
func (s server) stop() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := s.secrets.Close(); err != nil {
		s.errCh <- err
	}

	for _, fn := range s.shutdownFuncs {
		if err := fn(); err != nil {
			s.errCh <- err
		}
	}

	s.httpServer.SetKeepAlivesEnabled(false)
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.errCh <- err
	}

	s.stopCh <- sig
}

// newTLSConfig returns a new tls.Config.
func newTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS13,
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
		},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}
