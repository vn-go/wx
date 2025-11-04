package wx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// httpClientType represents an HTTP client with a base URI
type httpClientType struct {
	baseUri string
	client  *http.Client // Optional: allow custom client
}

// buildURL constructs the full URL from base URI and API path
func (c *httpClientType) buildURL(apiUri string) (string, error) {
	base, err := url.Parse(c.baseUri)
	if err != nil {
		return "", err
	}

	// Ensure apiUri is treated as a path
	apiPath, err := url.Parse(apiUri)
	if err != nil {
		return "", err
	}

	// Join paths properly
	fullURL := base.ResolveReference(apiPath)
	return fullURL.String(), nil
}

// PostJson sends a POST request with JSON body to the specified API endpoint
func (c *httpClientType) PostJson(apiUri string, data any) (any, error) {
	// Marshal the data into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Construct full URL by joining baseUri and apiUri
	fullURL, err := c.buildURL(apiUri)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Use custom client if set, otherwise use default
	client := c.client
	if client == nil {
		client = http.DefaultClient
	}

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode response
	var result any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return result, nil
}
func NewHttpClient(baseUri string) httpClientType {
	return httpClientType{
		baseUri: baseUri,
	}
}
