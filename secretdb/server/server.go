package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/RedeployAB/burnit/common/auth"
	"github.com/RedeployAB/burnit/secretdb/api"
	"github.com/RedeployAB/burnit/secretdb/config"
	"github.com/RedeployAB/burnit/secretdb/db"
)

// Server represents server with configuration.
type Server struct {
	*http.Server
	DB   *db.DB
	Port string
}

// NewServer returns a configured Server.
func NewServer(config config.Configuration, tokenStore auth.TokenStore, db *db.DB) *Server {

	r := api.NewRouter(config, tokenStore, db)

	srv := &http.Server{
		Addr:         "0.0.0.0:" + config.Server.Port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	return &Server{srv, db, config.Server.Port}
}

// Start creates an http server and runs ListenAndServe().
func (s *Server) Start() {
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server err: %v\n", err)
		}
	}()
	log.Printf("server listening on: %s", s.Port)
	s.gracefulShutdown()
}

// gracefulShutdown will shutdown server when interrupt or
// kill is received.
func (s *Server) gracefulShutdown() {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, os.Kill)
	sig := <-stop

	log.Printf("shutting down server. reason: %s\n", sig.String())
	log.Println("closing connection to database...")
	if err := s.DB.Close(); err != nil {
		log.Printf("database: %v", err)
	}
	log.Println("disonnected from database.")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Turn of SetKeepAlive when awaiting shutdown.
	s.SetKeepAlivesEnabled(false)
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("could not shutdown server gracefully: %v", err)
	}
	log.Println("server has been stopped")
}
