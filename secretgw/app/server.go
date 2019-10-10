package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/RedeployAB/burnit/secretgw/client"
	"github.com/RedeployAB/burnit/secretgw/config"
)

// Server represents server with configuration.
type Server struct {
	httpServer       *http.Server
	router           *mux.Router
	middlewareConfig middlewareConfig
	generatorService client.APIClient
	dbService        client.APIClient
}

type middlewareConfig struct {
	dbAPIkey string
}

// NewServer returns a configured Server.
func NewServer(config config.Configuration, r *mux.Router) *Server {
	srv := &http.Server{
		Addr:         "0.0.0.0:" + config.Port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	return &Server{
		httpServer: srv,
		router:     r,
		middlewareConfig: middlewareConfig{
			dbAPIkey: config.Server.DBAPIKey,
		},
		generatorService: client.APIClient{
			BaseURL: config.GeneratorBaseURL,
			Path:    config.GeneratorServicePath,
		},
		dbService: client.APIClient{
			BaseURL: config.DBBaseURL,
			Path:    config.DBServicePath,
		},
	}
}

// Start creates an http server and runs ListenAndServe().
func (s *Server) Start() {
	// Setup routes.
	s.routes()
	// Listen and Serve.
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server err: %v\n", err)
		}
	}()
	log.Printf("server listening on: %s", s.httpServer.Addr)
	s.gracefulShutdown()
}

// gracefulShutdown will shutdown server when interrupt or
// kill is received.
func (s *Server) gracefulShutdown() {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, os.Kill)
	sig := <-stop
	log.Printf("shutting down server. reason: %s\n", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Turn of SetKeepAlive when awaiting shutdown.
	s.httpServer.SetKeepAlivesEnabled(false)
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("could not shutdown server gracefully: %v", err)
	}
	log.Println("server has been stopped")
}
