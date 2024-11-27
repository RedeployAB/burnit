package server

import (
	"strconv"
	"time"

	"github.com/RedeployAB/burnit/internal/log"
	"github.com/RedeployAB/burnit/internal/session"
	"github.com/RedeployAB/burnit/internal/ui"
)

// WithOptions configures the server with the given Options.
func WithOptions(options Options) Option {
	return func(s *server) {
		if options.Router != nil {
			s.router = options.Router
			s.httpServer.Handler = s.router
		}
		if options.Logger != nil {
			s.log = options.Logger
		}
		if len(options.Host) > 0 || options.Port > 0 {
			s.httpServer.Addr = options.Host + ":" + strconv.Itoa(options.Port)
		}
		if options.ReadTimeout > 0 {
			s.httpServer.ReadTimeout = options.ReadTimeout
		}
		if options.WriteTimeout > 0 {
			s.httpServer.WriteTimeout = options.WriteTimeout
		}
		if options.IdleTimeout > 0 {
			s.httpServer.IdleTimeout = options.IdleTimeout
		}
		if !options.TLS.isEmpty() {
			s.tls = options.TLS
		}
		if !options.CORS.isEmpty() {
			s.cors = options.CORS
		}
		if !options.RateLimiter.isEmpty() {
			s.rateLimiter = options.RateLimiter
		}
	}
}

// WithLogger configures the server with the given logger.
func WithLogger(log log.Logger) Option {
	return func(s *server) {
		if log != nil {
			s.log = log
		}
	}
}

// WithRouter configures the server with the given router.
func WithRouter(router *router) Option {
	return func(s *server) {
		if router != nil {
			s.router = router
			s.httpServer.Handler = s.router
		}
	}
}

// WithAddress configures the server with the given address.
func WithAddress(addr string) Option {
	return func(s *server) {
		if len(addr) > 0 {
			s.httpServer.Addr = addr
		}
	}
}

// WithTimeouts configures the server with the given timeouts.
func WithTimeouts(read, write, idle time.Duration) Option {
	return func(s *server) {
		if read > 0 {
			s.httpServer.ReadTimeout = read
		}
		if write > 0 {
			s.httpServer.WriteTimeout = write
		}
		if idle > 0 {
			s.httpServer.IdleTimeout = idle
		}
	}
}

// WithTLS configures the server with the given TLS configuration.
func WithTLS(tls TLSConfig) Option {
	return func(s *server) {
		if !tls.isEmpty() {
			s.tls = tls
		}
	}
}

// WithCORS configures the server with the given CORS configuration.
func WithCORS(cors CORS) Option {
	return func(s *server) {
		if !cors.isEmpty() {
			s.cors = cors
		}
	}
}

// WithRateLimiter configures the server with the given rate limiter.
func WithRateLimiter(rateLimiter RateLimiter) Option {
	return func(s *server) {
		if !rateLimiter.isEmpty() {
			s.rateLimiter = rateLimiter
		}
	}
}

// WithUI configures the server with the given UI.
func WithUI(ui ui.UI) Option {
	return func(s *server) {
		if ui != nil {
			s.ui = ui
		}
	}
}

// WithSessionStore configures the server with the given session store.
func WithSessionStore(sessions session.Store) Option {
	return func(s *server) {
		if sessions != nil {
			s.sessions = sessions
		}
	}
}
