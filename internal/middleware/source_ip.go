package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
)

const (
	// contextKeySourceIP is the context key for the source IP address.
	contextKeySourceIP contextKey = 1
)

const (
	// SourceIPNotAvailable is the value for when the source IP is not available.
	SourceIPNotAvailable = "N/A"
)

// SourceIP is a middleware that sets the source IP address in the request context.
func SourceIP() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := resolveIP(r)
			next.ServeHTTP(w, r.WithContext(setSourceIP(r.Context(), ip)))
		})
	}
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
		return SourceIPNotAvailable
	}
	return ip
}

// setSourceIP sets the source IP address in the request context.
func setSourceIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, contextKeySourceIP, ip)
}

// getSourceIP returns the source IP address from the request context.
func getSourceIP(ctx context.Context) string {
	val := ctx.Value(contextKeySourceIP)
	ip, ok := val.(string)
	if !ok {
		return SourceIPNotAvailable
	}
	return ip
}
