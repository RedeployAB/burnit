package app

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RedeployAB/burnit/common/auth"
	"github.com/RedeployAB/burnit/common/security"
	"github.com/RedeployAB/burnit/secretdb/config"
	"github.com/RedeployAB/burnit/secretdb/db"
	"github.com/RedeployAB/burnit/secretdb/internal/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var correctPassphrase = security.Hash("passphrase")
var id1 = "507f1f77bcf86cd799439011"
var oid1, _ = primitive.ObjectIDFromHex(id1)

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
		model = &models.Secret{ID: oid1, Secret: "values"}
		err = nil
	case "find-invalid-oid":
		model = nil
		err = &db.QueryError{Message: "invalid oid", Code: -1}
	case "find-not-found":
		model = nil
		err = &db.QueryError{Message: "not found", Code: 0}
	case "find-success-passphrase":
		model = &models.Secret{ID: oid1, Secret: "values", Passphrase: correctPassphrase}
		err = nil
	case "find-error":
		model = nil
		err = errors.New("db error")
	case "find-delete-error":
		model = &models.Secret{ID: oid1, Secret: "values"}
		err = nil
	}
	return model, err
}

func (r *mockHandlerRepository) Insert(s *models.Secret) (*models.Secret, error) {
	return &models.Secret{}, nil
}

func (r *mockHandlerRepository) Delete(id string) (int64, error) {
	var result int64
	var err error

	if r.action == "find" && r.mode == "find-delete-error" {
		result = 0
		err = errors.New("delete error")
	}
	return result, err
}

// The different methods on the handler will require
// states. When creating
func SetupServer(action, mode string) Server {
	conf := config.Configuration{
		Server: config.Server{
			Port:       "5000",
			DBAPIKey:   "ABCDEF",
			Passphrase: "testphrase",
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

	connection := &mockConnection{}
	repo := &mockHandlerRepository{action: action, mode: mode}
	ts := auth.NewMemoryTokenStore()
	ts.Set(conf.Server.DBAPIKey, "app")

	r := mux.NewRouter()
	srv := &Server{
		router:     r,
		connection: connection,
		repository: repo,
		tokenStore: ts,
	}
	srv.routes(ts)
	return *srv
}

func TestGetSecretSuccess(t *testing.T) {
	apiKey := "ABCDEF"
	req, _ := http.NewRequest("GET", "/api/v0/secrets/"+id1, nil)
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("find", "find-success").router.ServeHTTP(res, req)

	if res.Code != 200 {
		t.Errorf("status code incorrect, got: %d, want 200", res.Code)
	}
}

func TestGetSecretInvalidOID(t *testing.T) {
	apiKey := "ABCDEF"
	req, _ := http.NewRequest("GET", "/api/v0/secrets/1234", nil)
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("find", "find-invalid-oid").router.ServeHTTP(res, req)

	if res.Code != 404 {
		t.Errorf("status code incorrect, got: %d, want 404", res.Code)
	}
}

func TestGetSecretNotFound(t *testing.T) {
	apiKey := "ABCDEF"
	req, _ := http.NewRequest("GET", "/api/v0/secrets/"+id1, nil)
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("find", "find-not-found").router.ServeHTTP(res, req)

	if res.Code != 404 {
		t.Errorf("status code incorrect, got: %d, want 404", res.Code)
	}
}

func TestGetSecretDBError(t *testing.T) {
	apiKey := "ABCDEF"
	req, _ := http.NewRequest("GET", "/api/v0/secrets/"+id1, nil)
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("find", "find-error").router.ServeHTTP(res, req)

	if res.Code != 500 {
		t.Errorf("status code incorrect, got: %d, want 404", res.Code)
	}
}

func TestGetSecretDeleteError(t *testing.T) {
	apiKey := "ABCDEF"
	req, _ := http.NewRequest("GET", "/api/v0/secrets/"+id1, nil)
	req.Header.Add("x-api-key", apiKey)
	res := httptest.NewRecorder()
	SetupServer("find", "find-delete-error").router.ServeHTTP(res, req)

	if res.Code != 500 {
		t.Errorf("status code incorrect, got: %d, want 404", res.Code)
	}
}

func TestGetSecretWithPassphraseSuccess(t *testing.T) {
	apiKey := "ABCDEF"
	req, _ := http.NewRequest("GET", "/api/v0/secrets/"+id1, nil)
	req.Header.Add("x-api-key", apiKey)
	req.Header.Add("x-passphrase", "passphrase")
	res := httptest.NewRecorder()
	SetupServer("find", "find-success-passphrase").router.ServeHTTP(res, req)

	if res.Code != 200 {
		t.Errorf("status code incorrect, got: %d, want 200", res.Code)
	}
}

func TestGetSecretWithInvalidPassphrase(t *testing.T) {
	apiKey := "ABCDEF"
	req, _ := http.NewRequest("GET", "/api/v0/secrets/"+id1, nil)
	req.Header.Add("x-api-key", apiKey)
	req.Header.Add("x-passphrase", "notpassphrase")
	res := httptest.NewRecorder()
	SetupServer("find", "find-success-passphrase").router.ServeHTTP(res, req)

	if res.Code != 401 {
		t.Errorf("status code incorrect, got: %d, want 401", res.Code)
	}
}
