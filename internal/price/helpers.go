package price

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
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

// retryWithBackoff retries the given function with exponential back-off.
// It retries up to maxRetries times, with a delay of baseDelay * 2^attempt.
func retryWithBackoff(ctx context.Context, maxRetries int, baseDelay time.Duration, fn func() error) error {
	var err error
	for attempt := 0; attempt < maxRetries; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}
		delay := time.Duration(math.Pow(2, float64(attempt))) * baseDelay
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return err
}

// makeRequest sends a POST request to the Zinc API and returns the parsed response.
func makeRequest(ctx context.Context, endpoint string, payload []byte, retailer domain.Retailer) ([]domain.Offer, error) {
	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %v", domain.ErrNetwork, err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	// TODO: Add API key from config.APIKey()

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Send request with retry
	var resp *http.Response
	err = retryWithBackoff(ctx, 3, 100*time.Millisecond, func() error {
		var err error
		resp, err = client.Do(req)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("%w: failed to send request: %v", domain.ErrNetwork, err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("%w: unauthorized request", domain.ErrAuth)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: unexpected status code: %d", domain.ErrNetwork, resp.StatusCode)
	}

	// Parse response
	var zincResp zincResponse
	if err := json.NewDecoder(resp.Body).Decode(&zincResp); err != nil {
		return nil, fmt.Errorf("%w: failed to parse response: %v", domain.ErrNetwork, err)
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
