package ui

import (
	"errors"
	"regexp"
	"strings"

	"github.com/RedeployAB/burnit/internal/secret"
)

// secretCreateResponse is the response data for a create secret request.
type secretCreateResponse struct {
	BaseURL        string
	ID             string
	Passphrase     string
	PassphraseHash string
	CSRFToken      string
}

// secretGetResponse is the response data for a get secret request.
type secretGetResponse struct {
	ID             string
	PassphraseHash string
	Value          string
	CSRFToken      string
}

// errorResponse is the response data for an error.
type errorResponse struct {
	RequestID string
	Title     string
	Message   string
}

// formatErrorMessage formats an error message.
func formatErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	re := regexp.MustCompile(`^.*:\s+`)
	msg = re.ReplaceAllString(msg, "")
	if !strings.HasSuffix(msg, ".") {
		msg += "."
	}
	return strings.ToUpper(msg[:1]) + msg[1:]
}

// isSecretBadRequestError returns true if the error is a bad request error.
func isSecretBadRequestError(err error) bool {
	errs := []error{
		secret.ErrValueInvalid,
		secret.ErrInvalidPassphrase,
		secret.ErrValueTooManyCharacters,
		secret.ErrInvalidExpirationTime,
		secret.ErrPassphraseInvalid,
		secret.ErrPassphraseTooFewCharacters,
		secret.ErrPassphraseTooManyCharacters,
	}

	for _, e := range errs {
		if errors.Is(err, e) {
			return true
		}
	}
	return false
}
