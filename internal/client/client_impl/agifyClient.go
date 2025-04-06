package client_impl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type AgifyClient struct {
	BaseURL string
}

type AgifyResponse struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Count int    `json:"count"`
}

func NewAgifyClient() *AgifyClient {
	return &AgifyClient{
		BaseURL: "https://api.agify.io/",
	}
}

func (c *AgifyClient) GetAgeByName(ctx context.Context, name string) (*int, error) {
	url := fmt.Sprintf("%s?name=%s", c.BaseURL, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to agify.io: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request agify.io: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned status: %d", c.BaseURL, resp.StatusCode)
	}

	var result AgifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse %s response: %w", c.BaseURL, err)
	}

	return &result.Age, nil
}
