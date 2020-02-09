package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RedeployAB/burnit/secretgw/internal/request"
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
		resBody.Data = mockGenerateResponse{Secret: "secretvalue"}
		err = nil
	case "gen-fail":
		resBody.Data = mockGenerateResponse{Secret: "fail"}
		err = errors.New("call to api failed")
	case "db-get-success":
		resBody.Data = mockDBResponse{
			ID:     "1234",
			Secret: "secretvalue",
		}
		err = nil
	case "db-get-fail":
		resBody.Data = mockDBResponse{}
		err = errors.New("call to api failed")
	case "db-create-success":
		resBody.Data = mockDBResponse{
			ID:     "4321",
			Secret: "secretvalue",
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

func TestCallSecretGenSuccess(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v0/generate", nil)
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

	if rb.Data.Secret != "secretvalue" {
		t.Errorf("response incorrect, got: %s, want: secretvalue", rb.Data.Secret)
	}
}

func TestCallSecretGenFail(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v0/generate", nil)
	res := httptest.NewRecorder()
	SetupServer("gen-fail").router.ServeHTTP(res, req)

	if res.Code != 500 {
		t.Errorf("status code incorrect, got: %d, want: 500", res.Code)
	}
}

func TestCallGetSecretSuccess(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v0/secrets/1234", nil)
	res := httptest.NewRecorder()
	SetupServer("db-get-success").router.ServeHTTP(res, req)

	if res.Code != 200 {
		t.Errorf("status code incorrect, got: %d, want: 200", res.Code)
	}

	var rb mockDBFullResponse
	b, err := ioutil.ReadAll(res.Body)

	if err = json.Unmarshal(b, &rb); err != nil {
		t.Errorf("Unmarshal failed")
	}

	if rb.Data.ID != "1234" {
		t.Errorf("Response incorrect, got: %s, want: 1234", rb.Data.ID)
	}

	if rb.Data.Secret != "secretvalue" {
		t.Errorf("Response incorrect, got: %s, want: secretvalue", rb.Data.Secret)
	}
}

func TestCallGetSecretFail(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v0/secrets/1234", nil)
	res := httptest.NewRecorder()
	SetupServer("db-get-fail").router.ServeHTTP(res, req)

	if res.Code != 500 {
		t.Errorf("Status code incorrect, got: %d, want: 500", res.Code)
	}
}

func TestCallCreateSecretSuccess(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/v0/secrets", nil)
	res := httptest.NewRecorder()
	SetupServer("db-create-success").router.ServeHTTP(res, req)

	if res.Code != 201 {
		t.Errorf("Status code incorrect, got: %d, want: 200", res.Code)
	}

	var rb mockDBFullResponse
	b, err := ioutil.ReadAll(res.Body)

	if err = json.Unmarshal(b, &rb); err != nil {
		t.Errorf("Unmarshal failed")
	}

	if rb.Data.ID != "4321" {
		t.Errorf("Response incorrect, got: %s, want: 4321", rb.Data.ID)
	}

	if rb.Data.Secret != "secretvalue" {
		t.Errorf("Response incorrect, got: %s, want: secretvalue", rb.Data.Secret)
	}
}

func TestCallCreateSecretFail(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/v0/secrets", nil)
	res := httptest.NewRecorder()
	SetupServer("db-create-fail").router.ServeHTTP(res, req)

	if res.Code != 500 {
		t.Errorf("Status code incorrect, got: %d, want: 500", res.Code)
	}
}

func TestCallNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	SetupServer("").router.ServeHTTP(res, req)

	if res.Code != 404 {
		t.Errorf("Status code incorrect, got: %d, want: 404", res.Code)
	}
}
