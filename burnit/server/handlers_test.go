package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RedeployAB/burnit/burnit/config"
	"github.com/RedeployAB/burnit/burnit/secret"
	"github.com/RedeployAB/burnit/common/auth"
	"github.com/gorilla/mux"
)

var id1 = "507f1f77bcf86cd799439011"
var apiKey = "ABCDEF"

// Mock to handle service actions in handler tests.
type mockSecretService struct {
	action string
	mode   string
}

func (r mockSecretService) Get(id, passphrase string) (*secret.Secret, error) {
	// Return different results based on underlying structs
	// state.
	var sec *secret.Secret
	var err error

	switch r.mode {
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

func (r mockSecretService) Create(s *secret.Secret) (*secret.Secret, error) {
	var sec *secret.Secret
	var err error

	switch r.mode {
	case "insert-success":
		sec = &secret.Secret{ID: id1, Value: "value"}
		err = nil
	case "insert-error":
		sec = nil
		err = errors.New("insert error")
	}
	return sec, err
}

func (r mockSecretService) DeleteExpired() (int64, error) {
	return 0, nil
}

func (r mockSecretService) Generate(l int, sc bool) *secret.Secret {
	return nil
}

func (r mockSecretService) Delete(id string) (int64, error) {
	var result int64
	var err error

	if r.action == "find" && r.mode == "find-delete-error" {
		result = 0
		err = errors.New("delete error")
	} else if r.action == "delete" {
		switch r.mode {
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

// The different methods on the handler will require
// states. When creating
func SetupServer(action, mode string) *server {
	conf := config.Configuration{}
	conf.Server.Security.APIKey = "ABCDEF"

	client := &mockClient{}
	service := &mockSecretService{action: action, mode: mode}
	ts := auth.NewMemoryTokenStore()
	ts.Set(conf.Server.Security.APIKey, "server")

	r := mux.NewRouter()
	srv := &server{
		router:        r,
		dbClient:      client,
		secretService: service,
		tokenStore:    ts,
	}
	srv.routes(ts)
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
		mode    string
		headers map[string]string
		param   string
		want    int
	}{
		{mode: "find-success", headers: authHeaders, param: id1, want: 200},
		{mode: "find-invalid-oid", headers: authHeaders, param: "1234", want: 404},
		{mode: "find-not-found", headers: authHeaders, param: id1, want: 404},
		{mode: "find-error", headers: authHeaders, param: id1, want: 500},
		{mode: "find-delete-error", headers: authHeaders, param: id1, want: 500},
		{mode: "find-passphrase-success", headers: authPassphraseSuccess, param: id1, want: 200},
		{mode: "find-passphrase-error", headers: authPassphraseFail, param: id1, want: 401},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("GET", "/secrets/"+test.param, nil)
		for k, v := range test.headers {
			req.Header.Add(k, v)
		}
		res := httptest.NewRecorder()
		SetupServer("find", test.mode).router.ServeHTTP(res, req)

		if res.Code != test.want {
			t.Errorf("status code incorrect, got: %d, want: %d", res.Code, test.want)
		}
	}
}

func TestCreateSecret(t *testing.T) {
	authHeaders := map[string]string{
		"api-key": apiKey,
	}

	jsonStr := []byte(`{"value":"value"}`)
	malformedJSONStr := []byte(`{"value":"value}`)

	var tests = []struct {
		mode         string
		headers      map[string]string
		body         []byte
		want         int
		wantLocation string
	}{
		{mode: "insert-success", headers: authHeaders, body: jsonStr, want: 201, wantLocation: "/secrets/" + id1},
		{mode: "insert-success", headers: authHeaders, body: malformedJSONStr, want: 400, wantLocation: ""},
		{mode: "insert-error", headers: authHeaders, body: jsonStr, want: 500, wantLocation: ""},
	}

	for _, test := range tests {
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
	}
}

func TestUpdateSecret(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/secrets/"+id1, nil)
	req.Header.Add("api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("", "").router.ServeHTTP(res, req)

	expectedCode := 501
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}
}

func TestDeleteSecret(t *testing.T) {
	authHeaders := map[string]string{
		"api-key": apiKey,
	}

	var tests = []struct {
		mode    string
		headers map[string]string
		param   string
		want    int
	}{
		{mode: "delete-success", headers: authHeaders, param: id1, want: 200},
		{mode: "delete-not-found", headers: authHeaders, param: id1, want: 404},
		{mode: "delete-error", headers: authHeaders, param: id1, want: 500},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("DELETE", "/secrets/"+test.param, nil)
		for k, v := range test.headers {
			req.Header.Add(k, v)
		}
		res := httptest.NewRecorder()
		SetupServer("delete", test.mode).router.ServeHTTP(res, req)

		if res.Code != test.want {
			t.Errorf("status code incorrect, got: %d, want: %d", res.Code, test.want)
		}
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
