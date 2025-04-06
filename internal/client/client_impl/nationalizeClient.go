package client_impl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type NationalizeClient struct {
	BaseURL string
}

type NationalizeResponse struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

func NewNationalityClient() *NationalizeClient {
	return &NationalizeClient{
		BaseURL: "https://api.nationalize.io/",
	}
}

func (c *NationalizeClient) GetNationalityByName(ctx context.Context, name string) (*string, error) {

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

	var result NationalizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse %s response: %w", c.BaseURL, err)
	}

	var country *string
	curProb := 0.0

	for _, v := range result.Country {
		if curProb <= v.Probability {
			country = &v.CountryID
		}
	}

	return country, nil
}
