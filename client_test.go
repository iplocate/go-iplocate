package iplocate

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client := NewClient(nil)
	assert.Equal(t, DefaultBaseURL, client.baseURL)
	assert.Equal(t, DefaultTimeout, client.httpClient.Timeout)
	assert.Empty(t, client.apiKey)
}

func TestNewClientWithCustomHTTPClient(t *testing.T) {
	customClient := &http.Client{Timeout: 60 * time.Second}
	client := NewClient(customClient)
	assert.Equal(t, DefaultBaseURL, client.baseURL)
	assert.Equal(t, customClient, client.httpClient)
	assert.Equal(t, 60*time.Second, client.httpClient.Timeout)
	assert.Empty(t, client.apiKey)
}

func TestWithAPIKey(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(nil).WithAPIKey(apiKey)
	assert.Equal(t, DefaultBaseURL, client.baseURL)
	assert.Equal(t, apiKey, client.apiKey)
}

func TestWithTimeout(t *testing.T) {
	client := NewClient(nil)
	customTimeout := 60 * time.Second
	client.WithTimeout(customTimeout)
	assert.Equal(t, customTimeout, client.httpClient.Timeout)
}

func TestWithBaseURL(t *testing.T) {
	client := NewClient(nil)
	customURL := "https://api.custom.com"
	client.WithBaseURL(customURL)
	assert.Equal(t, customURL, client.baseURL)

	// Test with trailing slash
	client.WithBaseURL("https://api.custom.com/")
	assert.Equal(t, "https://api.custom.com", client.baseURL)
}

func TestChainedConfiguration(t *testing.T) {
	apiKey := "test-key"
	timeout := 45 * time.Second
	baseURL := "https://custom.api.com"

	client := NewClient(nil).
		WithAPIKey(apiKey).
		WithTimeout(timeout).
		WithBaseURL(baseURL)

	assert.Equal(t, apiKey, client.apiKey)
	assert.Equal(t, timeout, client.httpClient.Timeout)
	assert.Equal(t, baseURL, client.baseURL)
}

func TestLookup_Success(t *testing.T) {
	// Mock response data
	mockResponse := LookupResponse{
		IP:          "8.8.8.8",
		Country:     stringPtr("United States"),
		CountryCode: stringPtr("US"),
		IsEU:        false,
		City:        stringPtr("Mountain View"),
		Continent:   stringPtr("North America"),
		Latitude:    float64Ptr(37.386),
		Longitude:   float64Ptr(-122.0838),
		TimeZone:    stringPtr("America/Los_Angeles"),
		PostalCode:  stringPtr("94035"),
		Subdivision: stringPtr("California"),
		Network:     stringPtr("8.8.8.0/24"),
		ASN: &ASN{
			ASN:         "AS15169",
			Route:       "8.8.8.0/24",
			Netname:     "GOOGLE",
			Name:        "Google LLC",
			CountryCode: "US",
			Domain:      "google.com",
			Type:        "hosting",
			RIR:         "ARIN",
		},
		Privacy: Privacy{
			IsAbuser:      false,
			IsAnonymous:   false,
			IsBogon:       false,
			IsHosting:     true,
			IsIcloudRelay: false,
			IsProxy:       false,
			IsTor:         false,
			IsVPN:         false,
		},
		Company: &Company{
			Name:        "Google LLC",
			Domain:      "google.com",
			CountryCode: "US",
			Type:        "hosting",
		},
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/lookup/8.8.8.8", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "go-iplocate/1.0.0", r.Header.Get("User-Agent"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient(nil).WithBaseURL(server.URL)
	result, err := client.Lookup("8.8.8.8")

	require.NoError(t, err)
	assert.Equal(t, mockResponse.IP, result.IP)
	assert.Equal(t, mockResponse.Country, result.Country)
	assert.Equal(t, mockResponse.CountryCode, result.CountryCode)
	assert.Equal(t, mockResponse.IsEU, result.IsEU)
	assert.Equal(t, mockResponse.ASN.ASN, result.ASN.ASN)
	assert.Equal(t, mockResponse.Privacy.IsHosting, result.Privacy.IsHosting)
}

func TestLookup_WithAPIKey(t *testing.T) {
	apiKey := "test-api-key"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, apiKey, r.URL.Query().Get("apikey"))

		mockResponse := LookupResponse{
			IP:          "8.8.8.8",
			Country:     stringPtr("United States"),
			CountryCode: stringPtr("US"),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient(nil).WithAPIKey(apiKey).WithBaseURL(server.URL)
	result, err := client.Lookup("8.8.8.8")

	require.NoError(t, err)
	assert.Equal(t, "8.8.8.8", result.IP)
}

func TestLookup_InvalidIP(t *testing.T) {
	client := NewClient(nil)
	_, err := client.Lookup("invalid-ip")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid IP address")
}

func TestLookup_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid IP address",
		})
	}))
	defer server.Close()

	client := NewClient(nil).WithBaseURL(server.URL)
	_, err := client.Lookup("127.0.0.1")

	require.Error(t, err)
	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	assert.Contains(t, apiErr.Message, "Invalid IP address")
}

func TestLookupSelf_Success(t *testing.T) {
	mockResponse := LookupResponse{
		IP:          "203.0.113.1",
		Country:     stringPtr("Example Country"),
		CountryCode: stringPtr("EX"),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/lookup/", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient(nil).WithBaseURL(server.URL)
	result, err := client.LookupSelf()

	require.NoError(t, err)
	assert.Equal(t, mockResponse.IP, result.IP)
	assert.Equal(t, mockResponse.Country, result.Country)
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		Message:    "Rate limit exceeded",
		StatusCode: 429,
	}

	expected := "IPLocate API error (429): Rate limit exceeded"
	assert.Equal(t, expected, err.Error())
}

// Helper functions for test data
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
