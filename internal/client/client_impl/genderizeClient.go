package client_impl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GenderizeClient struct {
	BaseURL string
}

type GenderizeResponse struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

func NewGenderizeClient() *GenderizeClient {
	return &GenderizeClient{
		BaseURL: "https://api.genderize.io/",
	}
}

func (c *GenderizeClient) GetGenderByName(ctx context.Context, name string) (*string, error) {
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

	var result GenderizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse %s response: %w", c.BaseURL, err)
	}

	return &result.Gender, nil
}
