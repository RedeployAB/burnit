package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/RedeployAB/redeploy-secrets/secret-generator/secrets"
)

// HandleGenerateSecret handles requests for secret generation.
func HandleGenerateSecret(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	params := handleGenerateSecretQuery(query)
	secret := secrets.GenerateSecret(params.Length, params.SpecialCharacters)
	// Set headers and response.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// Responde with JSON.
	sr := SecretResponse{Secret: secret}
	if err := json.NewEncoder(w).Encode(&sr); err != nil {
		panic(err)
	}
}

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
