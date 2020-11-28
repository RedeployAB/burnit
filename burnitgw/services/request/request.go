package request

import (
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
	Request(o Options) (*Response, error)
}

// client reqpresents a HTTP/HTTPS client
// to be used against other services.
type client struct {
	Address string
	Path    string
}

// NewClient creates and returns a new
// HTTP client with the provided Address and Path.
func NewClient(address, path string) Client {
	return &client{
		Address: address,
		Path:    path,
	}
}

// Request performs a request against the target URL.
func (c client) Request(o Options) (*Response, error) {
	url := c.Address + c.Path
	if o.Params["id"] != "" {
		url += "/" + o.Params["id"]
	}
	if o.Query != "" {
		url += "?" + o.Query
	}

	if (o.Method == POST || o.Method == PUT) && o.Body == nil {
		return nil, &Error{code: 400, err: "400 Bad Request"}
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
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusNoContent {
		return nil, &Error{code: res.StatusCode, err: res.Status}
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &Response{StatusCode: res.StatusCode, Body: b}, nil
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

// Response is the result from the request
// including status code and body read into
// []byte.
type Response struct {
	StatusCode int
	Body       []byte
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
