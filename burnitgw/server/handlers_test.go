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
		req, _ := http.NewRequest("GET", "/generate", nil)
		res := httptest.NewRecorder()
		SetupServer(test.mode).router.ServeHTTP(res, req)

		if res.Code != test.wantCode {
			t.Errorf("status code incorrect, got: %d, want: %d", res.Code, test.wantCode)
		}

		if test.mode == "gen-success" {
			var rb mockGenerateFullResponse
			b, err := ioutil.ReadAll(res.Body)

			if err = json.Unmarshal(b, &rb); err != nil {
				t.Fatalf("error in test: %v", err)
			}

			if rb.Data.Secret != test.wantSecret {
				t.Errorf("response incorrect, got: %s, want: %s", rb.Data.Secret, test.wantSecret)
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
			var rb mockDBFullResponse
			b, err := ioutil.ReadAll(res.Body)
			if err = json.Unmarshal(b, &rb); err != nil {
				t.Fatalf("error in test: %v", err)
			}

			if rb.Data.ID != test.param {
				t.Errorf("response incorrect, got: %s, want: %s", rb.Data.ID, test.param)
			}
			if rb.Data.Secret != test.wantSecret {
				t.Errorf("response incorrect, got: %s, want: %s", rb.Data.Secret, test.wantSecret)
			}
		}
	}
}

func TestCreateSecret(t *testing.T) {
	jsonStr := []byte(`{"secret":"value"}`)
	malformedJSONstr := []byte(`{"secret)`)

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
			var rb mockDBFullResponse
			b, err := ioutil.ReadAll(res.Body)
			if err = json.Unmarshal(b, &rb); err != nil {
				t.Fatalf("error in test: %v", err)
			}

			if rb.Data.ID != test.wantID {
				t.Errorf("response incorrect, got: %s, want: %s", rb.Data.ID, test.wantID)
			}
			if rb.Data.Secret != test.wantSecret {
				t.Errorf("response incorrect, got: %s, want: %s", rb.Data.Secret, test.wantSecret)
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
