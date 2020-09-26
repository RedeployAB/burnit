package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RedeployAB/burnit/burnitgen/config"
)

// Server represents server with configuration.
type Server struct {
	httpServer *http.Server
	router     *http.ServeMux
}

// New returns a configured Server.
func New(conf *config.Configuration, r *http.ServeMux) *Server {
	srv := &http.Server{
		Addr:         "0.0.0.0:" + conf.Port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	return &Server{httpServer: srv, router: r}
}

// Start creates an http server and runs ListenAndServe().
func (s *Server) Start() {
	s.routes()

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server err: %v\n", err)
		}
	}()
	log.Printf("server listening on: %s", s.httpServer.Addr)
	s.shutdown()
	log.Println("server has been stopped")
}

// shutdown will shutdown server when interrupt or
// kill is received.
func (s *Server) shutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, os.Kill)
	sig := <-stop

	log.Printf("shutting down server. reason: %s\n", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.httpServer.SetKeepAlivesEnabled(false)
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("could not shutdown server gracefully: %v", err)
	}
}
