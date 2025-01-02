package server

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/RedeployAB/burnit/internal/secret"
	"github.com/google/go-cmp/cmp"
)

func TestServer_generateSecret(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			secrets secret.Service
			req     *http.Request
		}
		want struct {
			status int
			body   []byte
		}
	}{
		{
			name: "generate secret - 8 characters",
			input: struct {
				secrets secret.Service
				req     *http.Request
			}{
				secrets: &stubSecretService{},
				req:     httptest.NewRequest("GET", "/secret?length=8", nil),
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusOK,
				body:   []byte(`{"value":"aaaaaaaa"}` + "\n"),
			},
		},
		{
			name: "generate secret - 8 characters with special characters",
			input: struct {
				secrets secret.Service
				req     *http.Request
			}{
				secrets: &stubSecretService{},
				req:     httptest.NewRequest("GET", "/secret?length=8&specialCharacters=true", nil),
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusOK,
				body:   []byte(`{"value":"a?a?a?a?"}` + "\n"),
			},
		},
		{
			name: "generate secret - 16 characters",
			input: struct {
				secrets secret.Service
				req     *http.Request
			}{
				secrets: &stubSecretService{},
				req:     httptest.NewRequest("GET", "/secret?length=16", nil),
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusOK,
				body:   []byte(`{"value":"aaaaaaaaaaaaaaaa"}` + "\n"),
			},
		},
		{
			name: "generate secret - 16 characters with special characters",
			input: struct {
				secrets secret.Service
				req     *http.Request
			}{
				secrets: &stubSecretService{},
				req:     httptest.NewRequest("GET", "/secret?length=16&specialCharacters=true", nil),
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusOK,
				body:   []byte(`{"value":"a?a?a?a?a?a?a?a?"}` + "\n"),
			},
		},
		{
			name: "generate secret - 8 characters and text/plain accept header",
			input: struct {
				secrets secret.Service
				req     *http.Request
			}{
				secrets: &stubSecretService{},
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/secret?length=8", nil)
					req.Header.Set("Accept", "text/plain")
					return req
				}(),
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusOK,
				body:   []byte("aaaaaaaa"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := test.input.req

			generateSecret(test.input.secrets, &stubLogger{}).ServeHTTP(rr, req)

			gotCode := rr.Code
			gotBody := rr.Body.Bytes()

			if diff := cmp.Diff(test.want.status, gotCode); diff != "" {
				t.Errorf("generateSecret() = unexpected status code (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.want.body, gotBody); diff != "" {
				t.Errorf("generateSecret() = unexpected body (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestServer_getSecret(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			secrets secret.Service
			req     *http.Request
			path    string
		}
		want struct {
			status int
			body   []byte
		}
	}{
		{
			name: "get secret",
			input: struct {
				secrets secret.Service
				req     *http.Request
				path    string
			}{
				secrets: &stubSecretService{
					secrets: []secret.Secret{
						{ID: "1", Value: "secret"},
					},
				},
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/secrets/1", nil)
					req.SetPathValue("id", "1")
					req.Header.Set("Passphrase", base64.StdEncoding.EncodeToString([]byte("passphrase")))
					return req
				}(),
				path: "/secret/1",
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusOK,
				body:   []byte(`{"value":"secret"}` + "\n"),
			},
		},
		{
			name: "get secret - passphrase required",
			input: struct {
				secrets secret.Service
				req     *http.Request
				path    string
			}{
				secrets: &stubSecretService{
					secrets: []secret.Secret{
						{ID: "1", Value: "secret", Passphrase: "passphrase"},
					},
					err: secret.ErrInvalidPassphrase,
				},
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/secrets/1", nil)
					req.SetPathValue("id", "1")
					return req
				}(),
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusUnauthorized,
				body:   []byte(`{"statusCode":401,"code":"PassphraseRequired","error":"passphrase required"}` + "\n"),
			},
		},
		{
			name: "get secret - invalid passphrase",
			input: struct {
				secrets secret.Service
				req     *http.Request
				path    string
			}{
				secrets: &stubSecretService{
					secrets: []secret.Secret{
						{ID: "1", Value: "secret", Passphrase: "passphrase"},
					},
					err: secret.ErrInvalidPassphrase,
				},
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/secrets/1", nil)
					req.SetPathValue("id", "1")
					req.Header.Set("Passphrase", base64.StdEncoding.EncodeToString([]byte("invalid")))
					return req
				}(),
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusUnauthorized,
				body:   []byte(`{"statusCode":401,"code":"InvalidPassphrase","error":"invalid passphrase"}` + "\n"),
			},
		},
		{
			name: "get secret - secret not found",
			input: struct {
				secrets secret.Service
				req     *http.Request
				path    string
			}{
				secrets: &stubSecretService{
					secrets: []secret.Secret{},
				},
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/secrets/1", nil)
					req.SetPathValue("id", "1")
					req.Header.Set("Passphrase", base64.StdEncoding.EncodeToString([]byte("invalid")))
					return req
				}(),
				path: "/secret/1",
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusNotFound,
				body:   []byte(`{"statusCode":404,"code":"SecretNotFound","error":"secret not found"}` + "\n"),
			},
		},
		{
			name: "get secret - error from service",
			input: struct {
				secrets secret.Service
				req     *http.Request
				path    string
			}{
				secrets: &stubSecretService{
					err: errSecretService,
				},
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/secrets/1", nil)
					req.SetPathValue("id", "1")
					req.Header.Set("Passphrase", base64.StdEncoding.EncodeToString([]byte("passphrase")))
					return req
				}(),
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusInternalServerError,
				body:   []byte(`{"statusCode":500,"code":"ServerError","error":"internal server error"}` + "\n"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := test.input.req

			getSecret(test.input.secrets, &stubLogger{}).ServeHTTP(rr, req)

			gotCode := rr.Code
			gotBody := rr.Body.Bytes()

			if diff := cmp.Diff(test.want.status, gotCode); diff != "" {
				t.Errorf("getSecret() = unexpected status code (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.want.body, gotBody); diff != "" {
				t.Errorf("getSecret() = unexpected body (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestServer_createSecret(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			secrets secret.Service
			req     *http.Request
		}
		want struct {
			status int
			body   []byte
		}
	}{
		{
			name: "create secret",
			input: struct {
				secrets secret.Service
				req     *http.Request
			}{
				secrets: &stubSecretService{},
				req:     httptest.NewRequest("POST", "/secret", strings.NewReader(`{"value":"1","ttl":"1h"}`)),
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusCreated,
				body:   []byte(`{"id":"1","passphrase":"passphrase","ttl":"1h0m0s"}` + "\n"),
			},
		},
		{
			name: "create secret - error empty value",
			input: struct {
				secrets secret.Service
				req     *http.Request
			}{
				secrets: &stubSecretService{},
				req:     httptest.NewRequest("POST", "/secret", strings.NewReader(`{"value":"","ttl":"1h"}`)),
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusBadRequest,
				body:   []byte(`{"statusCode":400,"code":"InvalidRequest","error":"invalid request: value is required"}` + "\n"),
			},
		},
		{
			name: "create secret - error from service",
			input: struct {
				secrets secret.Service
				req     *http.Request
			}{
				secrets: &stubSecretService{
					err: errSecretService,
				},
				req: httptest.NewRequest("POST", "/secret", strings.NewReader(`{"value":"1","ttl":"1h"}`)),
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusInternalServerError,
				body:   []byte(`{"statusCode":500,"code":"ServerError","error":"internal server error"}` + "\n"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := test.input.req

			createSecret(test.input.secrets, &stubLogger{}).ServeHTTP(rr, req)

			gotCode := rr.Code
			gotBody := rr.Body.Bytes()

			if diff := cmp.Diff(test.want.status, gotCode); diff != "" {
				t.Errorf("createSecret() = unexpected status code (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.want.body, gotBody); diff != "" {
				t.Errorf("createSecret() = unexpected body (-want +got)\n%s\n", diff)
			}
		})
	}
}

type stubSecretService struct {
	secrets []secret.Secret
	err     error
}

func (s stubSecretService) Start() error {
	return nil
}

func (s stubSecretService) Generate(options ...secret.GenerateOption) string {
	opts := secret.GenerateOptions{}
	for _, option := range options {
		option(&opts)
	}

	var builder strings.Builder
	for i := 0; i < opts.Length; i++ {
		if opts.SpecialCharacters && i%2 != 0 {
			builder.WriteString("?")
		} else {
			builder.WriteString("a")
		}
	}
	return builder.String()
}

func (s stubSecretService) Get(id, passphrase string, options ...secret.GetOption) (secret.Secret, error) {
	if s.err != nil {
		return secret.Secret{}, s.err
	}

	var sec secret.Secret
	for _, s := range s.secrets {
		if s.ID == id {
			sec = s
			break
		}
	}
	if sec == (secret.Secret{}) {
		return secret.Secret{}, secret.ErrSecretNotFound
	}

	if len(sec.Passphrase) == 0 {
		return sec, nil
	}
	if sec.Passphrase != passphrase {
		return secret.Secret{}, secret.ErrInvalidPassphrase
	}

	return sec, nil
}

func (s *stubSecretService) Create(se secret.Secret) (secret.Secret, error) {
	if s.err != nil {
		return secret.Secret{}, s.err
	}

	if s.secrets == nil {
		s.secrets = make([]secret.Secret, 0)
	}

	var ids []string
	for _, secret := range s.secrets {
		ids = append(ids, secret.ID)
	}
	sort.Strings(ids)

	var id string
	if len(ids) == 0 {
		id = "1"
	} else {
		last := ids[len(ids)-1]
		lastNum, _ := strconv.Atoi(last)
		lastNum++
		id = strconv.Itoa(lastNum)
	}

	secret := secret.Secret{ID: id, Value: se.Value, Passphrase: "passphrase", TTL: se.TTL}
	s.secrets = append(s.secrets, secret)
	return secret, nil
}

func (s stubSecretService) Delete(id string, options ...secret.DeleteOption) error {
	return nil
}

func (s stubSecretService) Close() error {
	return nil
}

func (s stubSecretService) Cleanup() chan error {
	return nil
}

var (
	errSecretService = errors.New("secret service error")
)
