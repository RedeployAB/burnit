package server

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RedeployAB/burnit/burnitgw/config"
	"github.com/RedeployAB/burnit/burnitgw/service/db"
	"github.com/RedeployAB/burnit/burnitgw/service/generator"
	"github.com/RedeployAB/burnit/burnitgw/service/request"
	"github.com/gorilla/mux"
)

// Server represents server with configuration.
type Server struct {
	httpServer       *http.Server
	router           *mux.Router
	middlewareConfig middlewareConfig
	generatorService generator.Service
	dbService        db.Service
	tlsConfig        tlsConfig
}

type tlsConfig struct {
	certificate string
	key         string
}

type middlewareConfig struct {
	dbAPIkey string
}

// New returns a configured Server.
func New(conf *config.Configuration, r *mux.Router) *Server {
	srv := &http.Server{
		Addr:         "0.0.0.0:" + conf.Port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	generatorService := generator.NewService(
		request.NewClient(
			conf.GeneratorAddress,
			conf.GeneratorServicePath,
		),
	)

	dbService := db.NewService(
		request.NewClient(
			conf.DBAddress,
			conf.DBServicePath,
		),
	)

	var tlsCfg tlsConfig
	if len(conf.TLS.Certificate) > 0 && len(conf.TLS.Key) > 0 {
		srv.TLSConfig = config.NewTLSConfig()
		srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
		tlsCfg = tlsConfig{
			certificate: conf.TLS.Certificate,
			key:         conf.TLS.Key,
		}
	}

	return &Server{
		httpServer: srv,
		router:     r,
		middlewareConfig: middlewareConfig{
			dbAPIkey: conf.Server.DBAPIKey,
		},
		generatorService: generatorService,
		dbService:        dbService,
		tlsConfig:        tlsCfg,
	}
}

// Start creates an http server and runs ListenAndServe().
func (s *Server) Start() {
	s.routes()

	go func() {
		if err := s.listenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v\n", err)
		}
	}()
	log.Printf("server listening on: %s\n", s.httpServer.Addr)
	s.shutdown()
	log.Println("server has been stopped")
}

// listenAndServe wraps around httpServer.Server ListenAndServe and
// ListenAndServeTLS depending on TLS configuration.
func (s *Server) listenAndServe() error {
	if (tlsConfig{}) != s.tlsConfig {
		return s.httpServer.ListenAndServeTLS(s.tlsConfig.certificate, s.tlsConfig.key)
	}
	return s.httpServer.ListenAndServe()
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
		log.Fatalf("server: %v", err)
	}
}
