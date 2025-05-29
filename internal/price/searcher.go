package price

import (
	"context"

	"savvyshopper/domain"
)

// searcher defines the interface for searching offers.
type searcher interface {
	Search(ctx context.Context, query string) ([]domain.Offer, error)
}

// amazonSearcher implements the searcher interface for Amazon.
type amazonSearcher struct {
	endpoint string
}

// walmartSearcher implements the searcher interface for Walmart.
type walmartSearcher struct {
	endpoint string
}

// NewAmazonSearcher creates a new amazonSearcher instance.
func NewAmazonSearcher(endpoint string) searcher {
	return &amazonSearcher{endpoint: endpoint}
}

// NewWalmartSearcher creates a new walmartSearcher instance.
func NewWalmartSearcher(endpoint string) searcher {
	return &walmartSearcher{endpoint: endpoint}
}

// Search for amazonSearcher.
func (s *amazonSearcher) Search(ctx context.Context, query string) ([]domain.Offer, error) {
	// Build payload
	payload, err := buildPayload(query, domain.Amazon)
	if err != nil {
		return nil, err
	}

	// Make request
	return makeRequest(ctx, s.endpoint, payload, domain.Amazon)
}

// Search for walmartSearcher.
func (s *walmartSearcher) Search(ctx context.Context, query string) ([]domain.Offer, error) {
	// Build payload
	payload, err := buildPayload(query, domain.Walmart)
	if err != nil {
		return nil, err
	}

	// Make request
	return makeRequest(ctx, s.endpoint, payload, domain.Walmart)
}
