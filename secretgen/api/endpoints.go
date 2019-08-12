package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/RedeployAB/redeploy-secrets/secretgen/secrets"
)

// NotFoundHandler handles all non used routes.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

// GenerateSecretHandler handles requests for secret generation.
func GenerateSecretHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	params := handleGenerateSecretQuery(query)
	secret := secrets.GenerateRandomString(params.Length, params.SpecialCharacters)
	// Set headers and response.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// Respond with JSON.
	sr := SecretResponseBody{Data: secretData{Secret: secret}}
	if err := json.NewEncoder(w).Encode(&sr); err != nil {
		panic(err)
	}
}

// Handles query parameters to get length and special character
// options.
func handleGenerateSecretQuery(query url.Values) secretParams {
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

	return secretParams{Length: length, SpecialCharacters: spChar}
}

type secretParams struct {
	Length            int
	SpecialCharacters bool
}
