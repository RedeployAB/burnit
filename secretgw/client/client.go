package client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	GET  string = "GET"
	POST string = "POST"
	PUT  string = "PUT"
)

// APIClient reqpresents a HTTP/HTTPS client
// to be used against other services.
type APIClient struct {
	BaseURL string
	Path    string
}

// Request performs a request against the target URL.
func (c APIClient) Request(o RequestOptions) (ResponseBody, error) {
	// Get the Base URL and Path from the struct.
	url := c.BaseURL + c.Path
	if o.Params["id"] != "" {
		url += "/" + o.Params["id"]
	}

	var r ResponseBody
	if (o.Method == POST || o.Method == PUT) && o.Body == nil {
		return r, &RequestError{code: 400, err: "bad request"}
	}

	client := &http.Client{}
	req, _ := http.NewRequest(o.Method, url, o.Body)
	if o.Header != nil && len(o.Header) > 0 {
		for k, v := range o.Header {
			req.Header.Add(k, v[0])
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return r, err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusNoContent {
		return r, &RequestError{code: res.StatusCode, err: res.Status}
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(b, &r)
	return r, nil
}

// RequestOptions is options for an HTTP/HTTPs request.
// Method, header, params (URL params) and body.
type RequestOptions struct {
	Method string
	Header http.Header
	Params map[string]string
	Body   io.Reader
}

// ResponseBody is the response returned
// to the client.
type ResponseBody struct {
	Data interface{} `json:"data"`
}

// RequestError implements error interface
// and represents errors encountered with Client
// requests.
type RequestError struct {
	err  string
	code int
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("code %d: %s", e.code, e.err)
}
