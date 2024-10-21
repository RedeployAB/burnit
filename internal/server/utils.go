package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
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

// parseGenerateSecretQuery parses the query parameters for length
// and special characters.
func parseGenerateSecretQuery(v url.Values) (int, bool) {
	var length string
	if l, ok := v["length"]; ok {
		length = l[0]
	}
	if l, ok := v["l"]; ok {
		length = l[0]
	}

	var specialCharacters string
	if sc, ok := v["specialCharacters"]; ok {
		specialCharacters = sc[0]
	}
	if sc, ok := v["sc"]; ok {
		specialCharacters = sc[0]
	}

	l, err := strconv.Atoi(length)
	if err != nil {
		l = defaultLength
	}
	sc, err := strconv.ParseBool(specialCharacters)
	if err != nil {
		sc = false
	}

	return l, sc
}

// resolveIP checks request for headers Forwarded, X-Forwarded-For, and X-Real-Ip
// and falls back to the RemoteAddr if none are found.
func resolveIP(r *http.Request) string {
	var addr string
	if f := r.Header.Get("Forwarded"); f != "" {
		for _, segment := range strings.Split(f, ",") {
			addr = strings.TrimPrefix(segment, "for=")
			break
		}
	} else if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		addr = strings.Split(xff, ",")[0]
	} else if xrip := r.Header.Get("X-Real-Ip"); xrip != "" {
		addr = xrip
	} else {
		addr = r.RemoteAddr
	}
	ip := strings.Split(addr, ":")[0]
	if net.ParseIP(ip) == nil {
		return "N/A"
	}
	return ip
}

// extractIDAndPassphrase extracts the ID and passphrase from the path.
func extractIDAndPassphrase(route, path string) (string, string, error) {
	path = strings.TrimPrefix(path, route)
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		return "", "", ErrInvalidPath
	}
	if len(parts) == 1 {
		return parts[0], "", nil
	}
	return parts[0], parts[1], nil
}

// decodeBase64 decodes a base64 string.
func decodeBase64(s string) (string, error) {
	s = strings.Replace(s, "=", "", -1)

	var encoding *base64.Encoding
	re := regexp.MustCompile(`[/+]`)
	if !re.MatchString(s) {
		encoding = base64.RawURLEncoding
	} else {
		encoding = base64.RawStdEncoding
	}

	b, err := encoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
