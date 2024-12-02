package middleware

import (
	"context"
	"net/http"
)

const (
	// contextKeyRequestID is the context key for the request ID.
	contextKeyRequestID contextKey = 0
)

// RequestID is a middleware that sets a unique request ID in the request context.
func RequestID() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(setRequestID(r.Context(), newUUID())))
		})
	}
}

// setRequestID sets the request ID in the request context.
func setRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKeyRequestID, id)
}

// getRequestID returns the request ID from the request context.
func getRequestID(ctx context.Context) string {
	val := ctx.Value(contextKeyRequestID)
	id, ok := val.(string)
	if !ok {
		return ""
	}
	return id
}

// RequestIDFromContext returns the request ID from the context.
func RequestIDFromContext(ctx context.Context) string {
	return getRequestID(ctx)
}
