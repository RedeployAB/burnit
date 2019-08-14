package middleware

import "net/http"

// Middleware takes and returns an http.Handler.
type Middleware func(http.Handler) http.Handler

// Chain takes an http.Handler as argument to process last, after
// processing the provided middleware.
func Chain(h http.Handler, mw ...Middleware) http.Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}
	return h
}
