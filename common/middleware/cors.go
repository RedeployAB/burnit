package middleware

import (
	"net/http"
	"strings"
)

// CORSHandler is used to add headers
// for CORS requests.
type CORSHandler struct {
	Origin  string
	Headers map[string][]string
}

// Handle writes the headers from CORSHandler to the response
// and passes it on.
func (c *CORSHandler) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", c.Origin)
		for k, v := range c.Headers {
			w.Header().Set(k, strings.Join(v, ", "))
		}
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
