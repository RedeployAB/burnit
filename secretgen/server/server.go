package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/RedeployAB/redeploy-secrets/secretgen/config"
	"github.com/gorilla/mux"
)

// Server represents server with configuration.
type Server struct {
	*http.Server
	Port string
}

// NewServer returns a configured Server.
func NewServer(config config.Configuration, r *mux.Router) *Server {
	srv := &Server{
		&http.Server{
			Addr:         "0.0.0.0:" + config.Port,
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      r,
		}, config.Port,
	}

	return srv
}

// Start creates an http server and runs ListenAndServe().
func (s *Server) Start() {
	log.Printf("server listening on: %s", s.Port)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server err: %v\n", err)
	}
}

// AddShutdownHook awaits os.Interrupt and os.Kill to start
// shutdown on server.
func (s *Server) AddShutdownHook(done chan<- bool) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Turn of SetKeepAlive when awaiting shutdown.
	s.SetKeepAlivesEnabled(false)
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("could not shutdown server gracefully: %v", err)
	}
	close(done)
}
