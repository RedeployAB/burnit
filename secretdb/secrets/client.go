package secrets

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// GeneratorClient is used to fetch secrets from generator.
type GeneratorClient struct {
	BaseURL string
	Path    string
}

// ResponseBody represents responses from secret generator.
type ResponseBody struct {
	Secret string `json:"secret"`
}

// Fetch fetches data from generator.
func (c GeneratorClient) Fetch() (ResponseBody, error) {
	url := c.BaseURL + c.Path

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	var r ResponseBody
	res, err := client.Do(req)
	if err != nil {
		return r, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(b, &r)

	return r, err
}
