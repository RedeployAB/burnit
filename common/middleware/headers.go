package middleware

import (
	"net/http"
	"strings"
)

// HeaderStrip is a middleware for handling
// headers in an HTTP request. This configurable
// struct exists to be compatible with Gorilla mux's
// Use() method.
type HeaderStrip struct {
	Exceptions []string
}

// Strip strips headers from requests.
// Keeps the header specified in Exception.
func (h *HeaderStrip) Strip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header != nil {
			header := http.Header{}
			for k, v := range r.Header {
				for _, e := range h.Exceptions {
					if strings.EqualFold(k, e) {
						header.Add(k, v[0])
					}
				}
			}
			r.Header = header
		}
		next.ServeHTTP(w, r)
	})
}
