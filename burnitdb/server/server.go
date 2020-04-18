package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/RedeployAB/burnit/burnitdb/config"
	"github.com/RedeployAB/burnit/burnitdb/db"
	"github.com/RedeployAB/burnit/common/auth"
	"github.com/RedeployAB/burnit/common/security"
	"github.com/gorilla/mux"
)

// Server represents server with configuration.
type Server struct {
	httpServer  *http.Server
	router      *mux.Router
	dbClient    db.Client
	repository  db.Repository
	tokenStore  auth.TokenStore
	compareHash func(hash, s string) bool
}

// Options represents options to be used with server.
type Options struct {
	Config     config.Configuration
	Router     *mux.Router
	DBClient   db.Client
	Repository db.Repository
	TokenStore auth.TokenStore
}

// New returns a configured Server.
func New(opts Options) *Server {
	srv := &http.Server{
		Addr:         "0.0.0.0:" + opts.Config.Server.Port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      opts.Router,
	}

	var compareHash func(hash, s string) bool
	switch opts.Config.Server.Security.HashMethod {
	case "md5":
		compareHash = security.CompareMD5Hash
	case "bcrypt":
		compareHash = security.CompareBcryptHash
	}

	return &Server{
		httpServer:  srv,
		router:      opts.Router,
		dbClient:    opts.DBClient,
		repository:  opts.Repository,
		tokenStore:  opts.TokenStore,
		compareHash: compareHash,
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

	if s.repository.GetDriver() == "mongo" {
		// Start cleanup expired entries go routine.
		go s.cleanup(&wg, cleanup)
		wg.Add(1)
	}

	log.Printf("server listening on: %s", s.httpServer.Addr)
	s.shutdown(&wg, cleanup)
	log.Println("server has been stopped")
}

// shutdown will shutdown server when interrupt or
// kill is received.
func (s *Server) shutdown(wg *sync.WaitGroup, cleanup chan<- bool) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, os.Kill)
	sig := <-stop

	log.Printf("shutting down server. reason: %s\n", sig.String())
	// Awaiting cleanup job to stop before proceeding.
	cleanup <- true
	wg.Wait()

	log.Println("closing connection to database...")
	if err := db.Close(s.dbClient); err != nil {
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
