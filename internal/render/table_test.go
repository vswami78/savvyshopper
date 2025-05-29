package render

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"savvyshopper/domain"
)

func TestTable_Golden(t *testing.T) {
	offers := []domain.Offer{
		{Title: "Short Title", Price: 10.99, Retailer: domain.Amazon, URL: "https://example.com/1"},
		{Title: "Very Long Title That Should Be Truncated Because It Exceeds Sixty Characters", Price: 20.99, Retailer: domain.Walmart, URL: "https://example.com/2"},
	}

	var buf bytes.Buffer
	if err := Table(&buf, offers); err != nil {
		t.Fatalf("Table() error = %v", err)
	}

	got := strings.TrimSpace(buf.String())
	goldenPath := filepath.Join("testdata", "table.golden")
	golden, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("Failed to read golden file: %v", err)
	}
	expected := strings.TrimSpace(string(golden))

	if got != expected {
		t.Errorf("Table() output mismatch:\nGot:\n%s\nExpected:\n%s", got, expected)
	}
}
