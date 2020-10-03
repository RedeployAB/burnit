package server

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/RedeployAB/burnit/burnitgen/secret"
	"github.com/RedeployAB/burnit/common/httperror"
)

// notFound handles all non used routes.
func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	httperror.Error(w, http.StatusNotFound)
}

// generateSecret handles requests for secret generation.
func (s *Server) generateSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		s := secret.Generate(parseGenerateSecretQuery(r.URL.Query()))

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(newSecretResponse(s)); err != nil {
			panic(err)
		}
	})
}

// parseGenerateSecretQuery parses query parameters to get length and special character
// options.
func parseGenerateSecretQuery(query url.Values) (int, bool) {
	lParam := query.Get("length")
	scParam := query.Get("specialchars")

	l, err := strconv.Atoi(lParam)
	if err != nil {
		l = 16
	}
	sc, err := strconv.ParseBool(scParam)
	if err != nil {
		sc = false
	}

	return l, sc
}
