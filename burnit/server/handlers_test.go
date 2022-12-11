package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RedeployAB/burnit/burnit/secret"
	"github.com/gorilla/mux"
)

var id1 = "507f1f77bcf86cd799439011"
var apiKey = "ABCDEF"

// Mock to handle service actions in handler tests.
type mockSecretService struct {
	action string
	mode   string
}

func (svc mockSecretService) Get(id, passphrase string) (*secret.Secret, error) {
	// Return different results based on underlying structs
	// state.
	var sec *secret.Secret
	var err error

	switch svc.mode {
	case "find-success":
		sec = &secret.Secret{ID: id1, Value: "value"}
		err = nil
	case "find-not-found":
		sec = nil
		err = nil
	case "find-passphrase-success":
		sec = &secret.Secret{ID: id1, Value: "value"}
		err = nil
	case "find-passphrase-error":
		sec = &secret.Secret{ID: id1, Value: ""}
		err = nil
	case "find-error":
		sec = nil
		err = errors.New("find error")
	case "find-delete-error":
		sec = &secret.Secret{ID: id1, Value: "value"}
		err = nil
	}
	return sec, err
}

func (svc mockSecretService) Create(s *secret.Secret) (*secret.Secret, error) {
	var sec *secret.Secret
	var err error

	switch svc.mode {
	case "insert-success":
		sec = &secret.Secret{ID: id1, Value: "value"}
		err = nil
	case "insert-error":
		sec = nil
		err = errors.New("insert error")
	}
	return sec, err
}

func (svc mockSecretService) DeleteExpired() (int64, error) {
	return 0, nil
}

func (svc mockSecretService) Generate(l int, sc bool) *secret.Secret {
	return nil
}

func (svc mockSecretService) Delete(id string) (int64, error) {
	var result int64
	var err error

	if svc.action == "find" && svc.mode == "find-delete-error" {
		result = 0
		err = errors.New("delete error")
	} else if svc.action == "delete" {
		switch svc.mode {
		case "delete-success":
			result = 1
			err = nil
		case "delete-not-found":
			result = 0
			err = nil
		case "delete-error":
			result = -10
			err = errors.New("db delete error")
		}
	}
	return result, err
}

func (svc mockSecretService) Start() error {
	return nil
}

func (svc mockSecretService) Stop() error {
	return nil
}

// The different methods on the handler will require
// states. When creating
func SetupServer(action, mode string) *server {
	service := &mockSecretService{action: action, mode: mode}

	r := mux.NewRouter()
	srv := &server{
		router:  r,
		secrets: service,
	}
	srv.routes()
	return srv
}

func TestGetSecret(t *testing.T) {
	authHeaders := map[string]string{
		"api-key": apiKey,
	}

	authPassphraseSuccess := map[string]string{
		"api-key":    apiKey,
		"passphrase": "passphrase",
	}

	authPassphraseFail := map[string]string{
		"api-key":    apiKey,
		"passphrase": "notpassphrase",
	}

	var tests = []struct {
		name    string
		mode    string
		headers map[string]string
		param   string
		want    int
	}{
		{
			name:    "find",
			mode:    "find-success",
			headers: authHeaders,
			param:   id1,
			want:    200,
		},
		{
			name:    "find - invalid oid",
			mode:    "find-invalid-oid",
			headers: authHeaders,
			param:   "1234",
			want:    404,
		},
		{
			name:    "find - not found",
			mode:    "find-not-found",
			headers: authHeaders,
			param:   id1,
			want:    404,
		},
		{
			name:    "find - error",
			mode:    "find-error",
			headers: authHeaders,
			param:   id1,
			want:    500,
		},
		{
			name:    "find - delete error",
			mode:    "find-delete-error",
			headers: authHeaders,
			param:   id1,
			want:    500,
		},
		{
			name:    "find - passphrase success",
			mode:    "find-passphrase-success",
			headers: authPassphraseSuccess,
			param:   id1,
			want:    200,
		},
		{
			name:    "find - passphrase error",
			mode:    "find-passphrase-error",
			headers: authPassphraseFail,
			param:   id1,
			want:    401,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/secrets/"+test.param, nil)
			for k, v := range test.headers {
				req.Header.Add(k, v)
			}
			res := httptest.NewRecorder()
			SetupServer("find", test.mode).router.ServeHTTP(res, req)

			if res.Code != test.want {
				t.Errorf("status code incorrect, got: %d, want: %d", res.Code, test.want)
			}
		})
	}
}

func TestCreateSecret(t *testing.T) {
	authHeaders := map[string]string{
		"api-key": apiKey,
	}

	jsonStr := []byte(`{"value":"value"}`)
	malformedJSONStr := []byte(`{"value":"value}`)

	var tests = []struct {
		name         string
		mode         string
		headers      map[string]string
		body         []byte
		want         int
		wantLocation string
	}{
		{
			name:         "insert",
			mode:         "insert-success",
			headers:      authHeaders,
			body:         jsonStr,
			want:         201,
			wantLocation: "/secrets/" + id1,
		},
		{
			name:         "insert - malformed JSON",
			mode:         "insert-success",
			headers:      authHeaders,
			body:         malformedJSONStr,
			want:         400,
			wantLocation: "",
		},
		{
			name:         "insert - error",
			mode:         "insert-error",
			headers:      authHeaders,
			body:         jsonStr,
			want:         500,
			wantLocation: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/secrets", bytes.NewBuffer(test.body))
			for k, v := range test.headers {
				req.Header.Add(k, v)
			}
			res := httptest.NewRecorder()
			SetupServer("create", test.mode).router.ServeHTTP(res, req)

			if res.Code != test.want {
				t.Errorf("status code was incorrect, got: %d, want: %d", res.Code, test.want)
			}
			locHdr := res.Header().Get("Location")
			if locHdr != test.wantLocation {
				t.Errorf("location header was incorrect, got: %s, want: %s", locHdr, test.wantLocation)
			}
		})
	}
}

func TestDeleteSecret(t *testing.T) {
	authHeaders := map[string]string{
		"api-key": apiKey,
	}

	var tests = []struct {
		name    string
		mode    string
		headers map[string]string
		param   string
		want    int
	}{
		{
			name:    "delete",
			mode:    "delete-success",
			headers: authHeaders,
			param:   id1,
			want:    200,
		},
		{
			name:    "delete - not found",
			mode:    "delete-not-found",
			headers: authHeaders,
			param:   id1,
			want:    404,
		},
		{
			name:    "delete - error",
			mode:    "delete-error",
			headers: authHeaders,
			param:   id1,
			want:    500,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/secrets/"+test.param, nil)
			for k, v := range test.headers {
				req.Header.Add(k, v)
			}
			res := httptest.NewRecorder()
			SetupServer("delete", test.mode).router.ServeHTTP(res, req)

			if res.Code != test.want {
				t.Errorf("status code incorrect, got: %d, want: %d", res.Code, test.want)
			}
		})
	}
}

func TestNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("", "").router.ServeHTTP(res, req)

	if res.Code != 404 {
		t.Errorf("Status code incorrect, got: %d, want: 404", res.Code)
	}
}
