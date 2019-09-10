package internal

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

// GeneratorResponseBody represents responses from secret generator.
type GeneratorResponseBody struct {
	Data generatorData `json:"data"`
}

type generatorData struct {
	Secret string `json:"secret"`
}

// Fetch fetches data from generator.
func (c GeneratorClient) Fetch() (GeneratorResponseBody, error) {
	url := c.BaseURL + c.Path

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	var r GeneratorResponseBody
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
