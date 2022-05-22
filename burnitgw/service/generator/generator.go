package generator

import (
	"encoding/json"
	"net/http"

	"github.com/RedeployAB/burnit/burnitgw/service/request"
)

// Secret represents the secret response
// from secret generator service.
type Secret struct {
	Value string `json:"value"`
}

// Service provides handling operations for secret
// generator service.
type Service interface {
	Generate(r *http.Request) (*Secret, error)
}

// service implements Service interface and
// provides a concrete representation
// of the secret generator service handler.
type service struct {
	client request.Client
}

// NewService creates a new service for handling secret
// generator actions.
func NewService(client request.Client) Service {
	return &service{client: client}
}

// Generate a new secret with a request to the secret
// generator service.
func (s service) Generate(r *http.Request) (*Secret, error) {
	res, err := s.client.Request(request.Options{
		Method: http.MethodGet,
		Query:  r.URL.RawQuery,
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
