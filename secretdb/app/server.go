package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/RedeployAB/burnit/common/auth"
	"github.com/RedeployAB/burnit/secretdb/config"
	"github.com/RedeployAB/burnit/secretdb/db"
	"github.com/gorilla/mux"
)

// Server represents server with configuration.
type Server struct {
	httpServer *http.Server
	router     *mux.Router
	connection db.Connection
	repository db.Repository
	tokenStore auth.TokenStore
}

// ServerOptions represents options to be used with server.
type ServerOptions struct {
	Config     config.Configuration
	Router     *mux.Router
	Connection db.Connection
	Repository db.Repository
	TokenStore auth.TokenStore
}

// NewServer returns a configured Server.
func NewServer(opts ServerOptions) *Server {
	srv := &http.Server{
		Addr:         "0.0.0.0:" + opts.Config.Server.Port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      opts.Router,
	}

	//repo := db.NewSecretRepository(opts.Connection, opts.Config.Server.Passphrase)

	return &Server{
		httpServer: srv,
		router:     opts.Router,
		connection: opts.Connection,
		repository: opts.Repository,
		tokenStore: opts.TokenStore,
	}
}

// Start creates an http server and runs ListenAndServe().
func (s *Server) Start() {
	// Setup WaitGroup and channel for cleanup job.
	var wg sync.WaitGroup
	cleanup := make(chan bool, 1)
	// Setup routes.
	s.routes(s.tokenStore)
	// Listen and Serve.
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server err: %v\n", err)
		}
	}()
	// Start cleanup expired entries go routine.
	go s.cleanup(&wg, cleanup)
	wg.Add(1)

	log.Printf("server listening on: %s", s.httpServer.Addr)
	s.shutdown(&wg, cleanup)
	log.Println("server has been stopped")
}

// shutdown will shutdown server when interrupt or
// kill is received.
func (s *Server) shutdown(wg *sync.WaitGroup, cleanup chan<- bool) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	sig := <-stop

	log.Printf("shutting down server. reason: %s\n", sig.String())
	// Awaiting cleanup job to stop before proceeding.
	cleanup <- true
	wg.Wait()

	log.Println("closing connection to database...")
	if err := db.Close(s.connection); err != nil {
		log.Printf("database: %v", err)
	}
	log.Println("disonnected from database")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Turn of SetKeepAlive when awaiting shutdown.
	s.httpServer.SetKeepAlivesEnabled(false)
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("could not shutdown server gracefully: %v", err)
	}
}

// cleanup runs the repository's DeleteExpired() to delete
// expited entries. Listens on bool receive channel and
// marks WaitGroup as done.
func (s *Server) cleanup(wg *sync.WaitGroup, stop <-chan bool) {
	for {
		select {
		case <-stop:
			log.Println("stopping cleanup task")
			wg.Done()
			return
		case <-time.After(5 * time.Second):
			_, err := s.repository.DeleteExpired()
			if err != nil {
				log.Printf("error in db expired cleanup: %v\n", err)
			}
		}
	}
}
