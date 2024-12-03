package server

import (
	"errors"
	"net/http"

	"github.com/RedeployAB/burnit/internal/api"
	"github.com/RedeployAB/burnit/internal/secret"
	"github.com/RedeployAB/burnit/internal/security"
)

var (
	// ErrEmptyRequest is returned when the request body is empty.
	ErrEmptyRequest = errors.New("empty request")
	// ErrMalformedRequest is returned when the request body is invalid.
	ErrMalformedRequest = errors.New("malformed request")
	// ErrInvalidRequest is returned when the request is invalid.
	ErrInvalidRequest = errors.New("invalid request")
	// ErrInvalidPath is returned when the path is invalid.
	ErrInvalidPath = errors.New("invalid path")
)

var (
	// ErrPassphraseRequired is returned when the passphrase is required.
	ErrPassphraseRequired = errors.New("passphrase required")
	// ErrPassphraseNotBase64 is returned when the passphrase is not base64 encoded.
	ErrPassphraseNotBase64 = errors.New("passphrase should be base64 encoded")
)

// writeError writes an error response to the caller.
func writeError(w http.ResponseWriter, statusCode int, code string, err error) {
	if err == nil {
		err = errors.New("internal server error")
	}

	respErr := api.Error{
		StatusCode: statusCode,
		Code:       code,
		Err:        err.Error(),
	}

	if err := encode(w, statusCode, respErr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// writeServerError writes a server error response to the caller.
func writeServerError(w http.ResponseWriter, requestID string) {
	respErr := api.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       "ServerError",
		Err:        "internal server error",
		RequestID:  requestID,
	}

	if err := encode(w, http.StatusInternalServerError, respErr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// errorCode returns the status code for the given error.
func errorCode(err error) (int, string) {
	for statusCode, errs := range errorCodeMaps {
		for e, code := range errs {
			if errors.Is(err, e) {
				return statusCode, code
			}
		}
	}
	return 0, ""
}

// errorCodeMaps maps errors to status codes.
var errorCodeMaps = map[int]map[error]string{
	http.StatusBadRequest: {
		ErrEmptyRequest:                       "EmptyRequest",
		ErrInvalidRequest:                     "InvalidRequest",
		ErrMalformedRequest:                   "MalformedRequest",
		ErrPassphraseNotBase64:                "PassphraseNotBase64",
		secret.ErrInvalidExpirationTime:       "InvalidExpirationTime",
		secret.ErrValueInvalid:                "ValueInvalid",
		secret.ErrValueTooManyCharacters:      "ValueTooManyCharacters",
		secret.ErrPassphraseInvalid:           "PassphraseInvalid",
		secret.ErrPassphraseTooFewCharacters:  "PassphraseTooFewCharacters",
		secret.ErrPassphraseTooManyCharacters: "PassphraseTooManyCharacters",
		security.ErrInvalidBase64:             "InvalidBase64",
	},
	http.StatusUnauthorized: {
		ErrPassphraseRequired:       "PassphraseRequired",
		secret.ErrInvalidPassphrase: "InvalidPassphrase",
	},
	http.StatusNotFound: {
		secret.ErrSecretNotFound: "SecretNotFound",
	},
}
