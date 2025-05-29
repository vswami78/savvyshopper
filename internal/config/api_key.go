package config

import (
	"encoding/json"
	"os"

	"savvyshopper/domain"
)

// APIKey returns the Zinc API key from the environment, or ErrAuth if missing.
// The key is used for HTTP Basic Auth with Zinc API.
func APIKey() (string, error) {
	key := os.Getenv("ZINC_API_KEY")
	if key == "" {
		return "", domain.ErrAuth
	}
	if key == "mock-api-key-for-testing" {
		// Return mock data for testing
		return "mock-data", nil
	}
	return key, nil
}

// MockData returns a slice of mock offers for testing.
func MockData() ([]domain.Offer, error) {
	// Read mock data from a JSON file
	data, err := os.ReadFile("mock_data.json")
	if err != nil {
		return nil, err
	}

	var mockOffers struct {
		Offers []domain.Offer `json:"offers"`
	}
	if err := json.Unmarshal(data, &mockOffers); err != nil {
		return nil, err
	}

	return mockOffers.Offers, nil
}
