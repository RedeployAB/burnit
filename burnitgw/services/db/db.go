package db

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RedeployAB/burnit/burnitgw/services/request"
)

// Secret represents the secret response
// from secret db service.
type Secret struct {
	ID        string    `json:"id,omitempty"`
	Value     string    `json:"value,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}

// Service provides handling operations for secret
// db service.
type Service interface {
	Get(r *http.Request, params map[string]string) (*Secret, error)
	Create(r *http.Request) (*Secret, error)
}

// service implements Service interface and
// provides a concrete representation
// of the secret db service handler.
type service struct {
	client request.Client
}

// NewService creates a new service for handling secret
// db actions.
func NewService(client request.Client) Service {
	return &service{client: client}
}

// Get a secret with a request to the secret db service.
// If no params (string map) is provided, params
// will be extracted from the request (id).
func (s service) Get(r *http.Request, params map[string]string) (*Secret, error) {
	res, err := s.client.Request(request.Options{
		Method: request.GET,
		Header: r.Header,
		Params: params,
	})
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 202 {
		var reqErr request.Error
		if err := json.Unmarshal(res.Body, &reqErr); err != nil {
			return nil, err
		}
		return nil, &reqErr
	}

	var secret Secret
	if err := json.Unmarshal(res.Body, &secret); err != nil {
		return nil, err
	}
	return &secret, nil
}

// Create a secret with a request to the sercret db service.
func (s service) Create(r *http.Request) (*Secret, error) {
	res, err := s.client.Request(request.Options{
		Method: request.POST,
		Header: r.Header,
		Body:   r.Body,
	})
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 202 {
		var resErr request.Error
		if err := json.Unmarshal(res.Body, &resErr); err != nil {
			return nil, err
		}
		return nil, &resErr
	}

	var secret Secret
	if err := json.Unmarshal(res.Body, &secret); err != nil {
		return nil, err
	}
	return &secret, nil
}
