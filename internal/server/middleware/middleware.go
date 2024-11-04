package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

// Middleware is a function that wraps an http.Handler and
// is intended to be used as a middleware.
type Middleware func(next http.Handler) http.Handler

// Chain chains the given middlewares to the provided http.Handler
// in the order they are provided and returns the final http.Handler.
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// resolveIP checks request for headers Forwarded, X-Forwarded-For, and X-Real-Ip
// and falls back to the RemoteAddr if none are found.
func resolveIP(r *http.Request) string {
	var addr string
	if f := r.Header.Get("Forwarded"); f != "" {
		for _, segment := range strings.Split(f, ",") {
			addr = strings.TrimPrefix(segment, "for=")
			break
		}
	} else if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		addr = strings.Split(xff, ",")[0]
	} else if xrip := r.Header.Get("X-Real-Ip"); xrip != "" {
		addr = xrip
	} else {
		addr = r.RemoteAddr
	}
	ip := strings.Split(addr, ":")[0]
	if net.ParseIP(ip) == nil {
		return "N/A"
	}
	return ip
}

// responseError represents an error response.
type responseError struct {
	StatusCode int    `json:"statusCode"`
	Err        string `json:"error"`
}

// Error returns the error message.
func (e *responseError) Error() string {
	return e.Err
}

// JSON returns the error as a JSON byte slice.
func (e responseError) JSON() []byte {
	b, _ := json.Marshal(e)
	return b
}
