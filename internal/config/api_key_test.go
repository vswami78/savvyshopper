package config

import (
	"os"
	"testing"

	"savvyshopper/domain"
)

func TestAPIKey(t *testing.T) {
	// Save original env var and restore after test
	originalKey := os.Getenv("ZINC_API_KEY")
	defer os.Setenv("ZINC_API_KEY", originalKey)

	tests := []struct {
		name    string
		envKey  string
		want    string
		wantErr error
	}{
		{
			name:    "key not set",
			envKey:  "",
			want:    "",
			wantErr: domain.ErrAuth,
		},
		{
			name:    "key set",
			envKey:  "test-key-123",
			want:    "test-key-123",
			wantErr: nil,
		},
		{
			name:    "mock key set",
			envKey:  "mock-api-key-for-testing",
			want:    "mock-data",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up test environment
			os.Setenv("ZINC_API_KEY", tt.envKey)

			// Run test
			got, err := APIKey()

			// Check results
			if got != tt.want {
				t.Errorf("APIKey() = %v, want %v", got, tt.want)
			}
			if err != tt.wantErr {
				t.Errorf("APIKey() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
