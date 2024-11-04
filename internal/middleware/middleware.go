package middleware

import "net/http"

// Middleware is a function that wraps an http.Handler and
// is intended to be used as a middleware.
type Middleware func(next http.Handler) http.Handler

// Chain chains the given middlewares to the provided http.Handler
// in the order they are provided and returns the final http.Handler.
func Chain(h http.Handler, middlewares ...func(next http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
