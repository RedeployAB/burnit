package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/RedeployAB/redeploy-secrets/secretgen/internal"
)

// notFound handles all non used routes.
func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

// generateSecret handles requests for secret generation.
func generateSecret(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	l, sc := parseGenerateSecretQuery(query)
	s := internal.GenerateRandomString(l, sc)
	// Set headers and response.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// Respond with JSON.
	sr := secretResponseBody{Data: secretResponse{Secret: s}}
	if err := json.NewEncoder(w).Encode(&sr); err != nil {
		panic(err)
	}
}

// parseGenerateSecretQuery parses query parameters to get length and special character
// options.
func parseGenerateSecretQuery(query url.Values) (int, bool) {
	lengthParam := query.Get("length")
	spCharParam := query.Get("specialchars")

	length, err := strconv.Atoi(lengthParam)
	if err != nil {
		length = 16
	}
	spChar, err := strconv.ParseBool(spCharParam)
	if err != nil {
		spChar = false
	}

	return length, spChar
}
