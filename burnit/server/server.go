package server

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/RedeployAB/burnit/burnit/config"
	"github.com/RedeployAB/burnit/burnit/secret"
	"github.com/gorilla/mux"
)

// server represents server with configuration.
type server struct {
	httpServer    *http.Server
	router        *mux.Router
	tls           tlsConfig
	configuration *config.Configuration
	middleware    middlewareConfig
	secrets       secret.Service
}

// Options represents options to be used with server.
type Options struct {
	Router        *mux.Router
	Configuration *config.Configuration
	Secrets       secret.Service
}

// tls contains configuration for TLS.
type tlsConfig struct {
	certificate string
	key         string
}

// middleware contains middleware.
type middlewareConfig struct {
	cors cors
}

// cors contains settings for CORS.
type cors struct {
	enabled bool
	origin  string
	headers http.Header
}

// middlewareConfig struct

// New returns a configured Server.
func New(opts Options) *server {
	srv := &http.Server{
		Addr:         opts.Configuration.Server.Host + ":" + opts.Configuration.Server.Port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      opts.Router,
	}

	var tlsCfg tlsConfig
	if len(opts.Configuration.Server.Security.TLS.Certificate) != 0 && len(opts.Configuration.Server.Security.TLS.Key) != 0 {
		srv.TLSConfig = config.NewTLSConfig()
		srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
		tlsCfg = tlsConfig{
			certificate: opts.Configuration.Server.Security.TLS.Certificate,
			key:         opts.Configuration.Server.Security.TLS.Key,
		}
	}

	var corsCfg cors
	if len(opts.Configuration.Server.Security.CORS.Origin) != 0 {
		corsCfg = cors{
			origin:  opts.Configuration.Server.Security.CORS.Origin,
			headers: config.CORSHeaders(),
			enabled: true,
		}
	}

	return &server{
		httpServer:    srv,
		router:        opts.Router,
		tls:           tlsCfg,
		configuration: opts.Configuration,
		middleware: middlewareConfig{
			cors: corsCfg,
		},
		secrets: opts.Secrets,
	}
}

// Start creates an http server and runs ListenAndServe().
func (s *server) Start() {
	s.routes()

	go func() {
		if err := s.listenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server: %v.\n", err)
		}
	}()

	var wg sync.WaitGroup
	cleanup := make(chan bool, 1)

	if s.configuration.Database.Driver == config.DatabaseDriverMongo {
		go s.cleanup(&wg, cleanup)
		wg.Add(1)
	}

	log.Printf("Server listening on: %s.\n", s.httpServer.Addr)
	s.shutdown(&wg, cleanup)
	log.Println("Server has been stopped.")
}

// listenAndServer wraps around httpServer.Server methods ListenAndServer
// and ListenAndServeTLS depending on TLS configuration.
func (s *server) listenAndServe() error {
	if (tlsConfig{}) != s.tls {
		return s.httpServer.ListenAndServeTLS(s.tls.certificate, s.tls.key)
	}
	return s.httpServer.ListenAndServe()
}

// shutdown will shutdown server when interrupt or
// kill is received.
func (s *server) shutdown(wg *sync.WaitGroup, cleanup chan<- bool) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	log.Printf("Shutting down server. Reason: %s.\n", sig.String())

	cleanup <- true
	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Stopping service...")
	if err := s.secrets.Stop(); err != nil {
		log.Printf("Service: %v.\n", err)
	}
	log.Println("Service stopped.")

	s.httpServer.SetKeepAlivesEnabled(false)
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server: %v\n", err)
	}
}

// cleanup runs the repository's DeleteExpired() to delete
// expited entries. Listens on bool receive channel and
// marks WaitGroup as done.
func (s *server) cleanup(wg *sync.WaitGroup, stop <-chan bool) {
	for {
		select {
		case <-stop:
			log.Println("Stopping cleanup task...")
			wg.Done()
			return
		case <-time.After(5 * time.Second):
			_, err := s.secrets.DeleteExpired()
			if err != nil {
				log.Printf("Database cleanup: %v.\n", err)
			}
		}
	}
}
