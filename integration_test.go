// +build integration

package iplocate

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRealAPI tests against the actual IPLocate API
// Run with: go test -tags=integration
func TestRealAPI(t *testing.T) {
	client := NewClient(nil)
	
	// Test with a well-known public IP (Google DNS)
	result, err := client.Lookup("8.8.8.8")
	require.NoError(t, err)
	
	assert.Equal(t, "8.8.8.8", result.IP)
	assert.NotNil(t, result.Country)
	assert.Equal(t, "United States", *result.Country)
	assert.NotNil(t, result.CountryCode)
	assert.Equal(t, "US", *result.CountryCode)
}

// TestRealAPIWithKey tests with an API key if available
func TestRealAPIWithKey(t *testing.T) {
	apiKey := os.Getenv("IPLOCATE_API_KEY")
	if apiKey == "" {
		t.Skip("IPLOCATE_API_KEY environment variable not set")
	}
	
	client := NewClient(nil).WithAPIKey(apiKey)
	
	result, err := client.Lookup("8.8.8.8")
	require.NoError(t, err)
	
	assert.Equal(t, "8.8.8.8", result.IP)
	assert.NotNil(t, result.Country)
}

// TestSelfLookup tests the self lookup functionality
func TestSelfLookup(t *testing.T) {
	client := NewClient(nil)
	
	result, err := client.LookupSelf()
	require.NoError(t, err)
	
	// Should return some IP address
	assert.NotEmpty(t, result.IP)
}