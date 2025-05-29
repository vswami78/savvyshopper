package e2e

import (
	"context"
	"strings"
	"testing"

	"savvyshopper/domain"
	"savvyshopper/internal/price"
	"savvyshopper/runner"
)

// TestRunnerEndToEnd verifies the runner works with mock searchers.
func TestRunnerEndToEnd(t *testing.T) {
	// Create mock searchers
	mockSearchers := map[domain.Retailer]price.Searcher{
		domain.Amazon:  &mockSearcher{retailer: domain.Amazon},
		domain.Walmart: &mockSearcher{retailer: domain.Walmart},
	}

	// Run the search
	var buf strings.Builder
	err := runner.Run(context.Background(), []string{"test query"}, &buf, mockSearchers)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Verify output
	output := buf.String()
	if !strings.Contains(output, "$") {
		t.Errorf("Output does not contain $ sign")
	}
	if !strings.Contains(output, "Amazon") || !strings.Contains(output, "Walmart") {
		t.Errorf("Output does not contain both retailer names")
	}
}

// mockSearcher implements price.Searcher for testing.
type mockSearcher struct {
	retailer domain.Retailer
}

func (m *mockSearcher) Search(ctx context.Context, query string) ([]domain.Offer, error) {
	// Return 3 mock offers
	return []domain.Offer{
		{
			Title:    "Test Product 1",
			Price:    19.99,
			URL:      "https://example.com/1",
			Retailer: m.retailer,
		},
		{
			Title:    "Test Product 2",
			Price:    29.99,
			URL:      "https://example.com/2",
			Retailer: m.retailer,
		},
		{
			Title:    "Test Product 3",
			Price:    39.99,
			URL:      "https://example.com/3",
			Retailer: m.retailer,
		},
	}, nil
}
