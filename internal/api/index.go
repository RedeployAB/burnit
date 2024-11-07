package api

import "context"

// Index represents the index of the service.
type Index struct {
	Name          string   `json:"name"`
	Version       string   `json:"version"`
	Documentation string   `json:"documentation,omitempty"`
	Endpoints     []string `json:"endpoints,omitempty"`
}

// Valid validates the Index.
func (i Index) Valid(ctx context.Context) map[string]string {
	return nil
}
