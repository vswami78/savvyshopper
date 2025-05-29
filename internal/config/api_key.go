package config

import (
	"os"

	"github.com/vswami78/savvyshopper/domain"
)

// APIKey returns the Zinc API key from the environment, or ErrAuth if missing.
// The key is used for HTTP Basic Auth with Zinc API.
func APIKey() (string, error) {
	key := os.Getenv("ZINC_API_KEY")
	if key == "" {
		return "", domain.ErrAuth
	}
	return key, nil
}