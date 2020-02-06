package app

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/RedeployAB/burnit/common/httperror"
	"github.com/RedeployAB/burnit/secretgen/internal/secret"
)

// notFound handles all non used routes.
func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	httperror.Error(w, http.StatusNotFound)
}

// generateSecret handles requests for secret generation.
func (s *Server) generateSecret(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	length, specialChars := parseGenerateSecretQuery(query)
	secret := secret.Generate(length, specialChars)
	// Set headers and response.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response(secret)); err != nil {
		panic(err)
	}
}

// parseGenerateSecretQuery parses query parameters to get length and special character
// options.
func parseGenerateSecretQuery(query url.Values) (int, bool) {
	lengthParam := query.Get("length")
	specialCharsParam := query.Get("specialchars")

	length, err := strconv.Atoi(lengthParam)
	if err != nil {
		length = 16
	}
	specialChars, err := strconv.ParseBool(specialCharsParam)
	if err != nil {
		specialChars = false
	}

	return length, specialChars
}
