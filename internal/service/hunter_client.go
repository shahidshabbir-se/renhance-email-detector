package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shahidshabbir-se/renhance-email-detector/internal/types"
)

const hunterBaseURL = "https://api.hunter.io/v2/domain-search"

type HunterClient struct {
	APIKey string
	Client *http.Client
}

func NewHunterClient(apiKey string) *HunterClient {
	return &HunterClient{
		APIKey: apiKey,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (hc *HunterClient) SearchDomain(domain string) (*types.HunterResponse, error) {
	url := fmt.Sprintf("%s?domain=%s&api_key=%s", hunterBaseURL, domain, hc.APIKey)

	resp, err := hc.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.HunterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
