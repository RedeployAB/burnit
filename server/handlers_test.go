package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RedeployAB/burnit/secret"
	"github.com/google/go-cmp/cmp"
)

func TestGenerateSecret(t *testing.T) {
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
	return secret.Secret{}, nil
}

func (s mockSecretService) Create(se secret.Secret) (secret.Secret, error) {
	return secret.Secret{}, nil
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
