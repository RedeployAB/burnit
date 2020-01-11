package request

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	// GET represents string GET.
	GET string = "GET"
	// POST represents string POST.
	POST string = "POST"
	// PUT represents string PUT.
	PUT string = "PUT"
)

// Client is an interface for HTTP requests.
type Client interface {
	Request(o Options) (ResponseBody, error)
}

// HTTPClient reqpresents a HTTP/HTTPS client
// to be used against other services.
type HTTPClient struct {
	BaseURL string
	Path    string
}

// NewHTTPClient creates and returns a new
// HTTPClient with the provided BaseURL and Path.
func NewHTTPClient(baseURL, path string) *HTTPClient {
	return &HTTPClient{
		BaseURL: baseURL,
		Path:    path,
	}
}

// Request performs a request against the target URL.
func (c HTTPClient) Request(o Options) (ResponseBody, error) {
	// Get the Base URL and Path from the struct.
	url := c.BaseURL + c.Path
	if o.Params["id"] != "" {
		url += "/" + o.Params["id"]
	}
	if o.Query != "" {
		url += "?" + o.Query
	}

	var r ResponseBody
	if (o.Method == POST || o.Method == PUT) && o.Body == nil {
		return r, &Error{code: 400, err: "bad request"}
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
		return r, &Error{code: res.StatusCode, err: res.Status}
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(b, &r)
	return r, nil
}

// Options is options for an HTTP/HTTPs request.
// Method, header, params (URL params) and body.
type Options struct {
	Method string
	Header http.Header
	Params map[string]string
	Query  string
	Body   io.Reader
}

// ResponseBody is the response returned
// to the client.
type ResponseBody struct {
	Data interface{} `json:"data"`
}

// Error implements error interface
// and represents errors encountered with Client
// requests.
type Error struct {
	err  string
	code int
}

func (e *Error) Error() string {
	return fmt.Sprintf("code %d: %s", e.code, e.err)
}
