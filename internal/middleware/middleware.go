package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

// Middleware is a function that wraps an http.Handler and
// is intended to be used as a middleware.
type Middleware func(next http.Handler) http.Handler

// contextKey is a custom type for context keys.
type contextKey int

// Chain chains the given middlewares to the provided http.Handler
// in the order they are provided and returns the final http.Handler.
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// newUUID generates a new UUID.
var newUUID = func() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New().String()
	}
	return id.String()
}
