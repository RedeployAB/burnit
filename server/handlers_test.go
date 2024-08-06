package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/RedeployAB/burnit/secret"
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
				secrets: &mockSecretService{},
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
				secrets: &mockSecretService{},
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
				secrets: &mockSecretService{},
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
				secrets: &mockSecretService{},
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
				secrets: &mockSecretService{},
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
			server := &server{
				secrets: test.input.secrets,
			}

			rr := httptest.NewRecorder()
			req := test.input.req
			server.generateSecret().ServeHTTP(rr, req)

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
				secrets: &mockSecretService{
					secrets: []secret.Secret{
						{ID: "1", Value: "secret"},
					},
				},
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/secret/1", nil)
					req.SetPathValue("id", "1")
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
			name: "get secret - passphrase in header",
			input: struct {
				secrets secret.Service
				req     *http.Request
				path    string
			}{
				secrets: &mockSecretService{
					secrets: []secret.Secret{
						{ID: "1", Value: "secret", Passphrase: "passphrase"},
					},
				},
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/secret/1", nil)
					req.SetPathValue("id", "1")
					req.Header.Set("Passphrase", "passphrase")
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
			name: "get secret - passphrase in path",
			input: struct {
				secrets secret.Service
				req     *http.Request
				path    string
			}{
				secrets: &mockSecretService{
					secrets: []secret.Secret{
						{ID: "1", Value: "secret", Passphrase: "passphrase"},
					},
				},
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/secret/1", nil)
					req.SetPathValue("id", "1")
					req.SetPathValue("passphrase", "passphrase")
					return req
				}(),
				path: "/secret/1/passphrase",
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
			name: "get secret - invalid passphrase",
			input: struct {
				secrets secret.Service
				req     *http.Request
				path    string
			}{
				secrets: &mockSecretService{
					secrets: []secret.Secret{
						{ID: "1", Value: "secret", Passphrase: "passphrase"},
					},
				},
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/secret/1", nil)
					req.SetPathValue("id", "1")
					req.SetPathValue("passphrase", "")
					return req
				}(),
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusUnauthorized,
				body:   []byte(`{"statusCode":401,"error":"invalid passphrase"}` + "\n"),
			},
		},
		{
			name: "get secret - secret not found",
			input: struct {
				secrets secret.Service
				req     *http.Request
				path    string
			}{
				secrets: &mockSecretService{
					secrets: []secret.Secret{},
				},
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/secret/1", nil)
					req.SetPathValue("id", "1")
					return req
				}(),
				path: "/secret/1",
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusNotFound,
				body:   []byte(`{"statusCode":404,"error":"secret not found"}` + "\n"),
			},
		},
		{
			name: "get secret - error from service",
			input: struct {
				secrets secret.Service
				req     *http.Request
				path    string
			}{
				secrets: &mockSecretService{
					err: errSecretService,
				},
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/secret/1", nil)
					req.SetPathValue("id", "1")
					return req
				}(),
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusInternalServerError,
				body:   []byte(`{"statusCode":500,"error":"internal server error"}` + "\n"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := &server{
				secrets: test.input.secrets,
			}

			rr := httptest.NewRecorder()
			req := test.input.req
			server.getSecret().ServeHTTP(rr, req)

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
				secrets: &mockSecretService{},
				req:     httptest.NewRequest("POST", "/secret", strings.NewReader(`{"value":"1","ttl":300}`)),
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusCreated,
				body:   []byte(`{"id":"1","ttl":300}` + "\n"),
			},
		},
		{
			name: "create secret - error empty value",
			input: struct {
				secrets secret.Service
				req     *http.Request
			}{
				secrets: &mockSecretService{},
				req:     httptest.NewRequest("POST", "/secret", strings.NewReader(`{"value":"","ttl":300}`)),
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusBadRequest,
				body:   []byte(`{"statusCode":400,"error":"invalid request: value is required"}` + "\n"),
			},
		},
		{
			name: "create secret - error from service",
			input: struct {
				secrets secret.Service
				req     *http.Request
			}{
				secrets: &mockSecretService{
					err: errSecretService,
				},
				req: httptest.NewRequest("POST", "/secret", strings.NewReader(`{"value":"1","ttl":300}`)),
			},
			want: struct {
				status int
				body   []byte
			}{
				status: http.StatusInternalServerError,
				body:   []byte(`{"statusCode":500,"error":"internal server error"}` + "\n"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := &server{
				secrets: test.input.secrets,
			}

			rr := httptest.NewRecorder()
			req := test.input.req
			server.createSecret().ServeHTTP(rr, req)

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

type mockSecretService struct {
	secrets []secret.Secret
	err     error
}

func (s mockSecretService) Generate(length int, specialCharacters bool) string {
	var builder strings.Builder
	for i := 0; i < length; i++ {
		if specialCharacters && i%2 != 0 {
			builder.WriteString("?")
		} else {
			builder.WriteString("a")
		}
	}
	return builder.String()
}

func (s mockSecretService) Get(id, passphrase string) (secret.Secret, error) {
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

func (s *mockSecretService) Create(se secret.Secret) (secret.Secret, error) {
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

	secret := secret.Secret{ID: id, Value: se.Value, TTL: se.TTL}
	s.secrets = append(s.secrets, secret)
	return secret, nil
}

func (s mockSecretService) Delete(id string) error {
	return nil
}

func (s mockSecretService) DeleteExpired() error {
	return nil
}

func (s mockSecretService) Close() error {
	return nil
}

var (
	errSecretService = errors.New("secret service error")
)
