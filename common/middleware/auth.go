package middleware

import (
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/common/auth"
	"github.com/RedeployAB/redeploy-secrets/common/httperror"
)

// AuthenticationMiddleware is used to authenticate
// requests against a stored application token.
type AuthenticationMiddleware struct {
	TokenStore auth.TokenStore
}

// Authenticate checks the incoming request for the header X-API-Key
// and verifies it against the Token Store.
func (amw *AuthenticationMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-API-Key")
		if token == "" {
			httperror.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ok := amw.TokenStore.Verify(token)
		if !ok {
			httperror.Error(w, "Forbidden", http.StatusForbidden)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// AuthHeaderMiddleware adds the X-API-Key header for
// each request.
type AuthHeaderMiddleware struct {
	Token string
}

// AddAuthHeader adds header containing API key/Token.
func (amw *AuthHeaderMiddleware) AddAuthHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("X-API-Key", amw.Token)
		next.ServeHTTP(w, r)
	})
}
