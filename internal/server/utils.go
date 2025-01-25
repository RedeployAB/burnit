package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/RedeployAB/burnit/internal/api"
)

const (
	// defaultLength is the default length of a secret.
	defaultLength = 16
)

// Validator is an interface that can be implemented by types that need to be validated.
type Validator interface {
	Valid(ctx context.Context) (errors map[string]string)
}

// encode writes the response as JSON to the response writer.
func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set(contentType, contentTypeJSON)
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("%w: %s", ErrMalformedRequest, err)
	}
	return nil
}

// decode reads the request body as JSON and decodes it into the given value.
func decode[T Validator](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		var syntaxError *json.SyntaxError
		if errors.As(err, &syntaxError) {
			err = ErrMalformedRequest
		}
		if errors.Is(err, io.EOF) {
			err = ErrEmptyRequest
		}
		if errors.Is(err, api.ErrInvalidTimeFormat) {
			err = fmt.Errorf("%w: %s", ErrInvalidRequest, err)
		}
		return v, err
	}
	if errors := v.Valid(r.Context()); len(errors) > 0 {
		var errs []string
		for _, v := range errors {
			errs = append(errs, v)
		}
		return v, fmt.Errorf("%w: %s", ErrInvalidRequest, strings.Join(errs, ", "))
	}
	return v, nil
}

// writeValue writes the value as text to the response writer.
func writeValue(w http.ResponseWriter, value string) {
	w.Header().Set(contentType, contentTypeText)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}
