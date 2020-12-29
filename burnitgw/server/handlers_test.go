package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RedeployAB/burnit/burnitgw/service/db"
	"github.com/RedeployAB/burnit/burnitgw/service/generator"
	"github.com/gorilla/mux"
)

type mockClient struct {
	address string
	path    string
	mode    string
}

type mockGeneratorService struct {
	client mockClient
}

func (s mockGeneratorService) Generate(r *http.Request) (*generator.Secret, error) {
	var secret *generator.Secret
	var err error
	switch s.client.mode {
	case "gen-success":
		secret = &generator.Secret{}
		secret.Value = "value"
	case "gen-fail":
		err = errors.New("call to api failed")
	}
	return secret, err
}

type mockDbService struct {
	client mockClient
}

func (s mockDbService) Get(r *http.Request, params map[string]string) (*db.Secret, error) {
	var secret *db.Secret
	var err error
	switch s.client.mode {
	case "db-get-success":
		secret = &db.Secret{}
		secret.ID = "1234"
		secret.Value = "value"
	case "db-get-fail":
		err = errors.New("call to api failed")
	}
	return secret, err
}

func (s mockDbService) Create(r *http.Request) (*db.Secret, error) {
	var secret *db.Secret
	var err error
	switch s.client.mode {
	case "db-create-success":
		secret = &db.Secret{}
		secret.ID = "4321"
		secret.Value = "value"
	case "db-create-fail":
		err = errors.New("call to api failed")
	}
	return secret, err
}

func SetupServer(mode string) Server {

	generatorService := mockGeneratorService{
		mockClient{
			mode: mode,
		},
	}

	dbService := mockDbService{
		mockClient{
			mode: mode,
		},
	}

	srv := Server{
		router:           mux.NewRouter(),
		generatorService: generatorService,
		dbService:        dbService,
	}
	srv.routes()
	return srv
}

func TestGenerateSecret(t *testing.T) {
	var tests = []struct {
		mode       string
		wantCode   int
		wantSecret string
	}{
		{mode: "gen-success", wantCode: 200, wantSecret: "value"},
		{mode: "gen-fail", wantCode: 500},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("GET", "/secret", nil)
		res := httptest.NewRecorder()
		SetupServer(test.mode).router.ServeHTTP(res, req)

		if res.Code != test.wantCode {
			t.Errorf("status code incorrect, got: %d, want: %d", res.Code, test.wantCode)
		}

		if test.mode == "gen-success" {
			var rb *generator.Secret
			b, err := ioutil.ReadAll(res.Body)
			if err = json.Unmarshal(b, &rb); err != nil {
				t.Fatalf("error in test: %v", err)
			}

			if rb.Value != test.wantSecret {
				t.Errorf("response incorrect, got: %s, want: %s", rb.Value, test.wantSecret)
			}
		}
	}
}

func TestGetSecret(t *testing.T) {
	var tests = []struct {
		mode       string
		param      string
		wantCode   int
		wantSecret string
	}{
		{mode: "db-get-success", param: "1234", wantCode: 200, wantSecret: "value"},
		{mode: "db-get-fail", param: "1234", wantCode: 500},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("GET", "/secrets/"+test.param, nil)
		res := httptest.NewRecorder()
		SetupServer(test.mode).router.ServeHTTP(res, req)

		if res.Code != test.wantCode {
			t.Errorf("status code incorrect, got: %d, want: %d", res.Code, test.wantCode)
		}

		if test.mode == "db-get-success" {
			var rb *db.Secret
			b, err := ioutil.ReadAll(res.Body)
			if err = json.Unmarshal(b, &rb); err != nil {
				t.Fatalf("error in test: %v", err)
			}

			if rb.ID != test.param {
				t.Errorf("response incorrect, got: %s, want: %s", rb.ID, test.param)
			}
			if rb.Value != test.wantSecret {
				t.Errorf("response incorrect, got: %s, want: %s", rb.Value, test.wantSecret)
			}
		}
	}
}

func TestCreateSecret(t *testing.T) {
	jsonStr := []byte(`{"value":"value"}`)
	malformedJSONstr := []byte(`{"value)`)

	var tests = []struct {
		mode       string
		body       []byte
		wantCode   int
		wantID     string
		wantSecret string
	}{
		{mode: "db-create-success", body: jsonStr, wantCode: 201, wantID: "4321", wantSecret: "value"},
		{mode: "db-create-fail", body: malformedJSONstr, wantCode: 500},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("POST", "/secrets", bytes.NewBuffer(test.body))
		res := httptest.NewRecorder()
		SetupServer(test.mode).router.ServeHTTP(res, req)

		if res.Code != test.wantCode {
			t.Errorf("status code incorrect, got: %d, want: %d", res.Code, test.wantCode)
		}

		if test.mode == "db-create-success" {
			var rb *db.Secret
			b, err := ioutil.ReadAll(res.Body)
			if err = json.Unmarshal(b, &rb); err != nil {
				t.Fatalf("error in test: %v", err)
			}

			if rb.ID != test.wantID {
				t.Errorf("response incorrect, got: %s, want: %s", rb.ID, test.wantID)
			}
			if rb.Value != test.wantSecret {
				t.Errorf("response incorrect, got: %s, want: %s", rb.Value, test.wantSecret)
			}
		}
	}
}

func TestNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	SetupServer("").router.ServeHTTP(res, req)

	expectedCode := 404
	if res.Code != expectedCode {
		t.Errorf("Status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}
}
