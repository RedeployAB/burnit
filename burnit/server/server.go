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
	"github.com/RedeployAB/burnit/burnit/db"
	"github.com/RedeployAB/burnit/burnit/secret"
	"github.com/RedeployAB/burnit/common/auth"
	"github.com/gorilla/mux"
)

// server represents server with configuration.
type server struct {
	httpServer    *http.Server
	router        *mux.Router
	tls           tlsConfig
	middleware    middlewareConfig
	dbClient      db.Client
	secretService secret.Service
	tokenStore    auth.TokenStore
	driver        int
}

// Options represents options to be used with server.
type Options struct {
	Config     *config.Configuration
	Router     *mux.Router
	DBClient   db.Client
	Repository db.Repository
	TokenStore auth.TokenStore
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
		Addr:         opts.Config.Server.Host + ":" + opts.Config.Server.Port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      opts.Router,
	}

	var driver int
	switch opts.Config.Database.Driver {
	case "redis":
		driver = 1
	case "mongo":
		driver = 2
	}

	secretService := secret.NewService(
		opts.Repository,
		secret.Options{EncryptionKey: opts.Config.Server.Security.Encryption.Key},
	)

	var tlsCfg tlsConfig
	if len(opts.Config.Server.Security.TLS.Certificate) != 0 && len(opts.Config.Server.Security.TLS.Key) != 0 {
		srv.TLSConfig = config.NewTLSConfig()
		srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
		tlsCfg = tlsConfig{
			certificate: opts.Config.Server.Security.TLS.Certificate,
			key:         opts.Config.Server.Security.TLS.Key,
		}
	}

	var corsCfg cors
	if len(opts.Config.Server.Security.CORS.Origin) != 0 {
		corsCfg = cors{
			origin:  opts.Config.Server.Security.CORS.Origin,
			headers: config.CORSHeaders(),
			enabled: true,
		}
	}

	return &server{
		httpServer: srv,
		router:     opts.Router,
		tls:        tlsCfg,
		middleware: middlewareConfig{
			cors: corsCfg,
		},
		dbClient:      opts.DBClient,
		secretService: secretService,
		driver:        driver,
		tokenStore:    opts.TokenStore,
	}
}

// Start creates an http server and runs ListenAndServe().
func (s *server) Start() {
	s.routes(s.tokenStore)

	go func() {
		if err := s.listenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v\n", err)
		}
	}()

	var wg sync.WaitGroup
	cleanup := make(chan bool, 1)

	if s.driver == 2 {
		go s.cleanup(&wg, cleanup)
		wg.Add(1)
	}

	log.Printf("server listening on: %s\n", s.httpServer.Addr)
	s.shutdown(&wg, cleanup)
	log.Println("server has been stopped")
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

	log.Printf("shutting down server. reason: %s\n", sig.String())

	cleanup <- true
	wg.Wait()

	log.Println("closing connection to database...")
	if err := db.Close(s.dbClient); err != nil {
		log.Printf("database: %v", err)
	}
	log.Println("disconnected from database")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.httpServer.SetKeepAlivesEnabled(false)
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("server: %v\n", err)
	}
}

// cleanup runs the repository's DeleteExpired() to delete
// expited entries. Listens on bool receive channel and
// marks WaitGroup as done.
func (s *server) cleanup(wg *sync.WaitGroup, stop <-chan bool) {
	for {
		select {
		case <-stop:
			log.Println("stopping cleanup task")
			wg.Done()
			return
		case <-time.After(5 * time.Second):
			_, err := s.secretService.DeleteExpired()
			if err != nil {
				log.Printf("db cleanup: %v\n", err)
			}
		}
	}
}
