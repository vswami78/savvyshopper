package price

import (
	"context"
	"errors"
	"sort"
	"testing"

	"savvyshopper/domain"
)

type mockSearcher struct {
	results []domain.Offer
	err     error
}

func (m *mockSearcher) Search(ctx context.Context, query string) ([]domain.Offer, error) {
	return m.results, m.err
}

func TestSearchPrices_MergeSortAndTruncate(t *testing.T) {
	// Mock searchers for Amazon and Walmart
	amazonResults := []domain.Offer{
		{Title: "A", Price: 10, Retailer: domain.Amazon},
		{Title: "B", Price: 5, Retailer: domain.Amazon},
		{Title: "C", Price: 20, Retailer: domain.Amazon},
	}
	walmartResults := []domain.Offer{
		{Title: "D", Price: 7, Retailer: domain.Walmart},
		{Title: "E", Price: 3, Retailer: domain.Walmart},
		{Title: "F", Price: 15, Retailer: domain.Walmart},
	}

	searchers := map[domain.Retailer]Searcher{
		domain.Amazon:  &mockSearcher{results: amazonResults},
		domain.Walmart: &mockSearcher{results: walmartResults},
	}

	results, err := SearchPrices(context.Background(), "test", searchers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 6 {
		t.Errorf("expected 6 results, got %d", len(results))
	}
	// Should be sorted by price
	if !sort.SliceIsSorted(results, func(i, j int) bool { return results[i].Price < results[j].Price }) {
		t.Errorf("results are not sorted by price")
	}
	// All prices should be >= 0
	for _, offer := range results {
		if offer.Price < 0 {
			t.Errorf("found negative price: %v", offer)
		}
	}
}

func TestSearchPrices_NoResults(t *testing.T) {
	searchers := map[domain.Retailer]Searcher{
		domain.Amazon:  &mockSearcher{results: nil},
		domain.Walmart: &mockSearcher{results: nil},
	}

	_, err := SearchPrices(context.Background(), "test", searchers)
	if !errors.Is(err, domain.ErrNoResults) {
		t.Errorf("expected ErrNoResults, got %v", err)
	}
}
