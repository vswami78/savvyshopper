package price

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"savvyshopper/domain"
)

func TestAmazonSearcher(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		// Verify content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Parse request body
		var payload zincPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("failed to parse request body: %v", err)
		}

		// Verify payload
		if payload.SearchTerm != "test" {
			t.Errorf("expected search term 'test', got %s", payload.SearchTerm)
		}
		if payload.Retailer != string(domain.Amazon) {
			t.Errorf("expected retailer 'Amazon', got %s", payload.Retailer)
		}
		if payload.MaxResults != 3 {
			t.Errorf("expected max results 3, got %d", payload.MaxResults)
		}

		// Send mock response
		response := zincResponse{
			Results: []struct {
				Title string  `json:"title"`
				Price float64 `json:"price"`
				URL   string  `json:"url"`
			}{
				{
					Title: "Test Product",
					Price: 19.99,
					URL:   "https://example.com/product",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create searcher with test server URL
	searcher := NewAmazonSearcher(server.URL)

	// Test search
	offers, err := searcher.Search(context.Background(), "test")
	if err != nil {
		t.Errorf("AmazonSearcher.Search() error = %v", err)
	}

	// Verify results
	if len(offers) != 1 {
		t.Errorf("expected 1 offer, got %d", len(offers))
	}
	if offers[0].Title != "Test Product" {
		t.Errorf("expected title 'Test Product', got %s", offers[0].Title)
	}
	if offers[0].Price != 19.99 {
		t.Errorf("expected price 19.99, got %f", offers[0].Price)
	}
	if offers[0].URL != "https://example.com/product" {
		t.Errorf("expected URL 'https://example.com/product', got %s", offers[0].URL)
	}
	if offers[0].Retailer != domain.Amazon {
		t.Errorf("expected retailer 'Amazon', got %s", offers[0].Retailer)
	}
}

func TestWalmartSearcher(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		// Verify content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Parse request body
		var payload zincPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("failed to parse request body: %v", err)
		}

		// Verify payload
		if payload.SearchTerm != "test" {
			t.Errorf("expected search term 'test', got %s", payload.SearchTerm)
		}
		if payload.Retailer != string(domain.Walmart) {
			t.Errorf("expected retailer 'Walmart', got %s", payload.Retailer)
		}
		if payload.MaxResults != 3 {
			t.Errorf("expected max results 3, got %d", payload.MaxResults)
		}

		// Send mock response
		response := zincResponse{
			Results: []struct {
				Title string  `json:"title"`
				Price float64 `json:"price"`
				URL   string  `json:"url"`
			}{
				{
					Title: "Test Product",
					Price: 29.99,
					URL:   "https://example.com/product",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create searcher with test server URL
	searcher := NewWalmartSearcher(server.URL)

	// Test search
	offers, err := searcher.Search(context.Background(), "test")
	if err != nil {
		t.Errorf("WalmartSearcher.Search() error = %v", err)
	}

	// Verify results
	if len(offers) != 1 {
		t.Errorf("expected 1 offer, got %d", len(offers))
	}
	if offers[0].Title != "Test Product" {
		t.Errorf("expected title 'Test Product', got %s", offers[0].Title)
	}
	if offers[0].Price != 29.99 {
		t.Errorf("expected price 29.99, got %f", offers[0].Price)
	}
	if offers[0].URL != "https://example.com/product" {
		t.Errorf("expected URL 'https://example.com/product', got %s", offers[0].URL)
	}
	if offers[0].Retailer != domain.Walmart {
		t.Errorf("expected retailer 'Walmart', got %s", offers[0].Retailer)
	}
}
