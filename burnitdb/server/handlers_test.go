package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RedeployAB/burnit/burnitdb/config"
	"github.com/RedeployAB/burnit/burnitdb/internal/models"
	"github.com/RedeployAB/burnit/common/auth"
	"github.com/RedeployAB/burnit/common/security"
	"github.com/gorilla/mux"
)

var correctPassphrase = security.Bcrypt("passphrase")
var id1 = "507f1f77bcf86cd799439011"
var apiKey = "ABCDEF"

// Mock to handle repository answers in handler tests.
type mockHandlerRepository struct {
	action string
	mode   string
}

func (r *mockHandlerRepository) Find(id string) (*models.Secret, error) {
	// Return different results based on underlying structs
	// state.
	var model *models.Secret
	var err error

	switch r.mode {
	case "find-success":
		model = &models.Secret{ID: id1, Secret: "values"}
		err = nil
	case "find-not-found":
		model = nil
		err = nil
	case "find-success-passphrase":
		model = &models.Secret{ID: id1, Secret: "values", Passphrase: correctPassphrase}
		err = nil
	case "find-error":
		model = nil
		err = errors.New("find error")
	case "find-delete-error":
		model = &models.Secret{ID: id1, Secret: "values"}
		err = nil
	}
	return model, err
}

func (r *mockHandlerRepository) Insert(s *models.Secret) (*models.Secret, error) {
	var model *models.Secret
	var err error

	switch r.mode {
	case "insert-success":
		model = &models.Secret{ID: id1, Secret: "value"}
		err = nil
	case "insert-error":
		model = nil
		err = errors.New("insert error")
	}
	return model, err
}

func (r *mockHandlerRepository) Delete(id string) (int64, error) {
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

func (r *mockHandlerRepository) DeleteExpired() (int64, error) {
	return 0, nil
}

func (r *mockHandlerRepository) GetDriver() string {
	return "mongo"
}

// The different methods on the handler will require
// states. When creating
func SetupServer(action, mode string) Server {
	conf := config.Configuration{
		Server: config.Server{
			Port: "5000",
			Security: config.Security{
				APIKey: "ABCDEF",
				Encryption: config.Encryption{
					Key: "testphrase",
				},
				HashMethod: "bcrypt",
			},
		},
		Database: config.Database{
			Address:  "mongo://db",
			Database: "secrets",
			Username: "user",
			Password: "password",
			SSL:      false,
			URI:      "",
		},
	}

	client := &mockClient{}
	repo := &mockHandlerRepository{action: action, mode: mode}
	ts := auth.NewMemoryTokenStore()
	ts.Set(conf.Server.Security.APIKey, "server")

	r := mux.NewRouter()
	srv := &Server{
		router:      r,
		dbClient:    client,
		repository:  repo,
		tokenStore:  ts,
		compareHash: security.CompareBcryptHash,
	}
	srv.routes(ts)
	return *srv
}

func TestGetSecretSuccess(t *testing.T) {
	req, _ := http.NewRequest("GET", "/secrets/"+id1, nil)
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("find", "find-success").router.ServeHTTP(res, req)

	if res.Code != 200 {
		t.Errorf("status code incorrect, got: %d, want 200", res.Code)
	}
}

func TestGetSecretInvalidOID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/secrets/1234", nil)
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("find", "find-invalid-oid").router.ServeHTTP(res, req)

	expectedCode := 404
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}
}

func TestGetSecretNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/secrets/"+id1, nil)
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("find", "find-not-found").router.ServeHTTP(res, req)

	expectedCode := 404
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d ", res.Code, expectedCode)
	}
}

func TestGetSecretDBError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/secrets/"+id1, nil)
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("find", "find-error").router.ServeHTTP(res, req)

	expectedCode := 500
	if res.Code != 500 {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}
}

func TestGetSecretDeleteError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/secrets/"+id1, nil)
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("find", "find-delete-error").router.ServeHTTP(res, req)

	expectedCode := 500
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}
}

func TestGetSecretWithPassphraseSuccess(t *testing.T) {
	req, _ := http.NewRequest("GET", "/secrets/"+id1, nil)
	req.Header.Add("x-api-key", apiKey)
	req.Header.Add("x-passphrase", "passphrase")
	res := httptest.NewRecorder()
	SetupServer("find", "find-success-passphrase").router.ServeHTTP(res, req)

	expectedCode := 200
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}
}

func TestGetSecretWithInvalidPassphrase(t *testing.T) {
	req, _ := http.NewRequest("GET", "/secrets/"+id1, nil)
	req.Header.Add("x-api-key", apiKey)
	req.Header.Add("x-passphrase", "notpassphrase")
	res := httptest.NewRecorder()
	SetupServer("find", "find-success-passphrase").router.ServeHTTP(res, req)

	expectedCode := 401
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}
}

func TestCreateSecretSuccess(t *testing.T) {
	jsonStr := []byte(`{"secret":"value"}`)
	req, _ := http.NewRequest("POST", "/secrets", bytes.NewBuffer(jsonStr))
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("create", "insert-success").router.ServeHTTP(res, req)

	expectedCode := 201
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}
}

func TestCreateSecretParseBodyError(t *testing.T) {
	// Creating faulty JSON.
	jsonStr := []byte(`{"secret":"value}`)
	req, _ := http.NewRequest("POST", "/secrets", bytes.NewBuffer(jsonStr))
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("insert", "insert-success").router.ServeHTTP(res, req)

	expectedCode := 400
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}
}

func TestCreateSecretCreateError(t *testing.T) {
	jsonStr := []byte(`{"secret":"value"}`)
	req, _ := http.NewRequest("POST", "/secrets", bytes.NewBuffer(jsonStr))
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("insert", "insert-error").router.ServeHTTP(res, req)

	expectedCode := 500
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}
}

func TestUpdateSecret(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/secrets/"+id1, nil)
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("", "").router.ServeHTTP(res, req)

	expectedCode := 501
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}
}

func TestDeleteSecretSuccess(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/secrets/"+id1, nil)
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("delete", "delete-success").router.ServeHTTP(res, req)

	expectedCode := 200
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}
}

func TestDeleteSecretNotFound(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/secrets/"+id1, nil)
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("delete", "delete-not-found").router.ServeHTTP(res, req)

	expectedCode := 404
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}
}

func TestDeleteSecretError(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/secrets/"+id1, nil)
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("delete", "delete-error").router.ServeHTTP(res, req)

	expectedCode := 500
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}
}

func TestNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("", "").router.ServeHTTP(res, req)

	if res.Code != 404 {
		t.Errorf("Status code incorrect, got: %d, want: 404", res.Code)
	}
}
