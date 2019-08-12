package internal

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// DBClient is used to fetch secrets from db.
type DBClient struct {
	BaseURL string
	Path    string
}

// DBResponseBody represents response from db.
type DBResponseBody struct {
	Data dbData `json:"data"`
}

type dbData struct {
	ID        string `json:"id,omitempty"`
	Secret    string `json:"secret,omitempty"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
}

// DBRequestBody represents a request body for the
// DB.
type DBRequestBody struct {
	Secret     string
	Passphrase string
}

// Do performs requests to secretdb service.
func (c DBClient) Do(method, id string, header http.Header, body io.Reader) (DBResponseBody, error) {

	url := c.BaseURL + c.Path
	if id != "" {
		url += "/" + id
	}

	var r DBResponseBody

	if method == "POST" && body == nil {
		return r, &DBRequestError{code: 400, err: "bad request"}
	}

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, body)
	if header != nil && len(header) > 0 {
		for k, v := range header {
			req.Header.Add(k, v[0])
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return r, err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusNoContent {
		return r, &DBRequestError{code: res.StatusCode, err: res.Status}
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(b, &r)

	return r, err
}
