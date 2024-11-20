package server

import (
	"net/http"

	"github.com/RedeployAB/burnit/internal/log"
	"github.com/RedeployAB/burnit/internal/middleware"
	"github.com/RedeployAB/burnit/internal/ui"
)

const (
	// defaultContentSecurityPolicy is the default content security policy.
	defaultContentSecurityPolicy = "default-src 'self';"
)

// routes sets up the routes for the server.
func (s *server) routes() {
	middlewares, shutdownFuncs := setupMiddlewares(s.log, s.rateLimiter, s.cors)
	s.shutdownFuncs = append(s.shutdownFuncs, shutdownFuncs...)

	// Secret router and handlers.
	secretRouter := http.NewServeMux()
	secretRouter.Handle("GET /secret", generateSecret(s.secrets, s.log))

	// Secrets router and handlers.
	secretsRouter := http.NewServeMux()
	secretsRouter.Handle("GET /secrets/{id}", getSecret(s.secrets, s.log))
	secretsRouter.Handle("POST /secrets", createSecret(s.secrets, s.log))
	secretsRouter.Handle("DELETE /secrets/{id}", deleteSecret(s.secrets, s.log))

	secretHandler := middleware.Chain(secretRouter, middlewares...)
	secretsHandler := middleware.Chain(secretsRouter, middlewares...)

	s.router.Handle("/secret", secretHandler)
	s.router.Handle("/secrets", secretsHandler)

	if s.ui == nil {
		s.router.Handle("/{$}", middleware.Logger(s.log)(index(nil, s.log)))
		s.router.Handle("/", middleware.Logger(s.log)(notFound(nil)))
		return
	}

	uiMiddlewares := setupUIMiddlewares(s.log, s.ui.RuntimeRender())

	fer := http.NewServeMux()
	fer.Handle("/ui/secrets", ui.CreateSecret(s.ui, s.secrets))
	fer.Handle("/ui/secrets/", ui.GetSecret(s.ui, s.secrets, s.log))
	fer.Handle("/ui/handlers/secret/get", middleware.HTMX(ui.HandlerGetSecret(s.ui, s.secrets, s.log)))
	fer.Handle("/ui/handlers/secret/create", middleware.HTMX(ui.HandlerCreateSecret(s.ui, s.secrets, s.log)))
	fer.Handle("/ui/", ui.NotFound(s.ui))

	uiHandler := middleware.Chain(fer, uiMiddlewares...)
	s.router.Handle("/ui/", uiHandler)

	s.router.Handle("/static/", middleware.Logger(s.log)(http.StripPrefix("/static/", ui.FileServer(s.ui.Static()))))
	s.router.Handle("/{$}", middleware.Logger(s.log)(index(s.ui, s.log)))
	s.router.Handle("/", middleware.Logger(s.log)(notFound(s.ui)))
}

// setupMiddlewares sets up the middlewares for the server.
func setupMiddlewares(log log.Logger, rl RateLimiter, c CORS) ([]middleware.Middleware, []func() error) {
	middlewares := []middleware.Middleware{middleware.Logger(log)}
	var shutdownFuncs []func() error
	if !rl.isEmpty() {
		mw, closeRateLimiter := middleware.RateLimiter(
			middleware.WithRateLimiterRate(rl.Rate),
			middleware.WithRateLimiterBurst(rl.Burst),
			middleware.WithRateLimiterTTL(rl.TTL),
			middleware.WithRateLimiterCleanupInterval(rl.CleanupInterval),
		)
		middlewares = append(middlewares, mw)
		shutdownFuncs = append(shutdownFuncs, closeRateLimiter)
	}
	if !c.isEmpty() {
		middlewares = append(middlewares, middleware.CORS(c.Origin))
	}
	middlewares = append(middlewares, middleware.Headers())
	return middlewares, shutdownFuncs
}

// setupUIMiddlewares sets up the middlewares for the UI.
func setupUIMiddlewares(log log.Logger, runtimeRender bool) []middleware.Middleware {
	middlewares := []middleware.Middleware{middleware.Logger(log)}

	var contentSecurityPolicy string
	if !runtimeRender {
		contentSecurityPolicy = defaultContentSecurityPolicy
	}
	middlewares = append(middlewares, middleware.Headers(func(o *middleware.HeadersOptions) {
		o.ContentSecurityPolicy = contentSecurityPolicy
	}))
	middlewares = append(middlewares, middleware.Compress())
	return middlewares
}
