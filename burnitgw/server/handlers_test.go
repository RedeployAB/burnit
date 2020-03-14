package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RedeployAB/burnit/burnitgw/internal/request"
	"github.com/gorilla/mux"
)

type mockClient struct {
	baseURL string
	path    string
	mode    string
}

type mockGenerateResponse struct {
	Secret string
}

type mockGenerateFullResponse struct {
	Data mockGenerateResponse
}

type mockDBResponse struct {
	ID        string
	Secret    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type mockDBFullResponse struct {
	Data mockDBResponse
}

func (c mockClient) Request(o request.Options) (request.ResponseBody, error) {

	resBody := request.ResponseBody{}
	var err error

	switch c.mode {
	case "gen-success":
		resBody.Data = mockGenerateResponse{Secret: "value"}
		err = nil
	case "gen-fail":
		resBody.Data = mockGenerateResponse{Secret: "fail"}
		err = errors.New("call to api failed")
	case "db-get-success":
		resBody.Data = mockDBResponse{
			ID:     "1234",
			Secret: "value",
		}
		err = nil
	case "db-get-fail":
		resBody.Data = mockDBResponse{}
		err = errors.New("call to api failed")
	case "db-create-success":
		resBody.Data = mockDBResponse{
			ID:     "4321",
			Secret: "value",
		}
		err = nil
	case "db-create-fail":
		resBody.Data = mockDBResponse{}
		err = errors.New("call to api failed")
	}

	return resBody, err
}

func SetupServer(mode string) Server {

	srv := Server{
		router:           mux.NewRouter(),
		generatorService: mockClient{mode: mode},
		dbService:        mockClient{mode: mode},
	}
	srv.routes()
	return srv
}

func TestGenerateSecretSuccess(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/generate", nil)
	res := httptest.NewRecorder()
	SetupServer("gen-success").router.ServeHTTP(res, req)

	if res.Code != 200 {
		t.Errorf("status code incorrect, got: %d, want: 200", res.Code)
	}

	var rb mockGenerateFullResponse
	b, err := ioutil.ReadAll(res.Body)

	if err = json.Unmarshal(b, &rb); err != nil {
		t.Errorf("error in test: %v", err)
	}

	expectedVale := "value"
	if rb.Data.Secret != expectedVale {
		t.Errorf("response incorrect, got: %s, want: %s", rb.Data.Secret, expectedVale)
	}
}

func TestGenerateSecretError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/generate", nil)
	res := httptest.NewRecorder()
	SetupServer("gen-fail").router.ServeHTTP(res, req)

	expectedCode := 500
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}
}

func TestGetSecretSuccess(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/secrets/1234", nil)
	res := httptest.NewRecorder()
	SetupServer("db-get-success").router.ServeHTTP(res, req)

	expectedCode := 200
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}

	var rb mockDBFullResponse
	b, err := ioutil.ReadAll(res.Body)

	if err = json.Unmarshal(b, &rb); err != nil {
		t.Errorf("error in test: %v", err)
	}

	expectedID := "1234"
	if rb.Data.ID != expectedID {
		t.Errorf("response incorrect, got: %s, want: %s", rb.Data.ID, expectedID)
	}

	expectedValue := "value"
	if rb.Data.Secret != expectedValue {
		t.Errorf("response incorrect, got: %s, want: %s", rb.Data.Secret, expectedValue)
	}
}

func TestGetSecretError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/secrets/1234", nil)
	res := httptest.NewRecorder()
	SetupServer("db-get-fail").router.ServeHTTP(res, req)

	expectedCode := 500
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}
}

func TestCreateSecretSuccess(t *testing.T) {
	jsonStr := []byte(`{"secret":"value"}`)
	req, _ := http.NewRequest("POST", "/api/secrets", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()
	SetupServer("db-create-success").router.ServeHTTP(res, req)

	expectedCode := 201
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}

	var rb mockDBFullResponse
	b, err := ioutil.ReadAll(res.Body)

	if err = json.Unmarshal(b, &rb); err != nil {
		t.Errorf("error in test: %v", err)
	}

	expectedID := "4321"
	if rb.Data.ID != expectedID {
		t.Errorf("response incorrect, got: %s, want: %s", rb.Data.ID, expectedID)
	}

	expectedValue := "value"
	if rb.Data.Secret != expectedValue {
		t.Errorf("response incorrect, got: %s, want: %s", rb.Data.Secret, expectedValue)
	}
}

func TestCreateSecretError(t *testing.T) {
	jsonStr := []byte(`{"secret}`)
	req, _ := http.NewRequest("POST", "/api/secrets", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()
	SetupServer("db-create-fail").router.ServeHTTP(res, req)

	expectedCode := 500
	if res.Code != expectedCode {
		t.Errorf("Status code incorrect, got: %d, want: %d", res.Code, expectedCode)
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
