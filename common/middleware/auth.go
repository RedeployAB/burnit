package middleware

import (
	"net/http"

	"github.com/RedeployAB/burnit/common/auth"
	"github.com/RedeployAB/burnit/common/httperror"
)

// Authentication is used to authenticate
// requests against a stored application token.
type Authentication struct {
	TokenStore auth.TokenStore
}

// Authenticate checks the incoming request for the header API-Key
// and verifies it against the Token Store.
func (amw *Authentication) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("API-Key")
		if len(token) == 0 {
			httperror.Write(w, http.StatusUnauthorized, "", "")
			return
		}

		ok := amw.TokenStore.Verify(token)
		if !ok {
			httperror.Write(w, http.StatusForbidden, "", "")
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// AuthHeader adds the API-Key header for
// each request.
type AuthHeader struct {
	Token string
}

// AddAuthHeader adds header containing API key/Token.
func (amw *AuthHeader) AddAuthHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("API-Key", amw.Token)
		next.ServeHTTP(w, r)
	})
}
