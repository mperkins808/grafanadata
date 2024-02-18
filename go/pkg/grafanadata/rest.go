package grafanadata

import (
	"fmt"
	"io"
	"net/http"
)

// Calls the http client Do method
func (c *grafanaClient) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	return resp, err
}

// NewRequest creates a new HTTP request with the API token included in the headers.
func (c *grafanaClient) NewRequest(method, endpoint string, body io.Reader) (*http.Request, error) {

	// Create a new HTTP request
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header with the Bearer token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.token))
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
