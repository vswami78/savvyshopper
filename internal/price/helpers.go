package price

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"savvyshopper/domain"
)

// zincPayload represents the JSON payload for Zinc API requests.
type zincPayload struct {
	SearchTerm string `json:"search_term"`
	Retailer   string `json:"retailer"`
	MaxResults int    `json:"max_results"`
}

// buildPayload creates a Zinc API payload for the given search term and retailer.
func buildPayload(query string, retailer domain.Retailer) ([]byte, error) {
	payload := zincPayload{
		SearchTerm: query,
		Retailer:   string(retailer),
		MaxResults: 3, // We only need top 3 results per retailer
	}
	return json.Marshal(payload)
}

// zincResponse represents the JSON response from Zinc API.
type zincResponse struct {
	Results []struct {
		Title string  `json:"title"`
		Price float64 `json:"price"`
		URL   string  `json:"url"`
	} `json:"results"`
}

// makeRequest sends a POST request to the Zinc API and returns the parsed response.
func makeRequest(ctx context.Context, endpoint string, payload []byte, retailer domain.Retailer) ([]domain.Offer, error) {
	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	// TODO: Add API key from config.APIKey()

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var zincResp zincResponse
	if err := json.NewDecoder(resp.Body).Decode(&zincResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to domain.Offer slice
	offers := make([]domain.Offer, len(zincResp.Results))
	for i, result := range zincResp.Results {
		offers[i] = domain.Offer{
			Title:    result.Title,
			Price:    result.Price,
			URL:      result.URL,
			Retailer: retailer,
		}
	}

	return offers, nil
}
