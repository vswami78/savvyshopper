package config

import (
	"testing"
)

func TestMockData(t *testing.T) {
	offers, err := MockData()
	if err != nil {
		t.Fatalf("MockData() error = %v", err)
	}

	if len(offers) != 2 {
		t.Errorf("MockData() returned %d offers, want 2", len(offers))
	}

	// Check first offer
	if offers[0].Title != "Mock Product 1" {
		t.Errorf("First offer title = %v, want Mock Product 1", offers[0].Title)
	}
	if offers[0].Price != 19.99 {
		t.Errorf("First offer price = %v, want 19.99", offers[0].Price)
	}
	if offers[0].Retailer != "Amazon" {
		t.Errorf("First offer retailer = %v, want Amazon", offers[0].Retailer)
	}

	// Check second offer
	if offers[1].Title != "Mock Product 2" {
		t.Errorf("Second offer title = %v, want Mock Product 2", offers[1].Title)
	}
	if offers[1].Price != 29.99 {
		t.Errorf("Second offer price = %v, want 29.99", offers[1].Price)
	}
	if offers[1].Retailer != "Walmart" {
		t.Errorf("Second offer retailer = %v, want Walmart", offers[1].Retailer)
	}
}
