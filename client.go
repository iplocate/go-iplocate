// Package iplocate provides a Go client for the IPLocate.io API.
// IPLocate.io provides comprehensive IP geolocation and threat intelligence data.
package iplocate

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the default base URL for the IPLocate API
	DefaultBaseURL = "https://iplocate.io/api"
	// DefaultTimeout is the default HTTP request timeout
	DefaultTimeout = 30 * time.Second
)

// Client represents an IPLocate API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new IPLocate client with the given HTTP client.
// If httpClient is nil, a default client with 30 second timeout is used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: DefaultTimeout,
		}
	}
	return &Client{
		baseURL:    DefaultBaseURL,
		httpClient: httpClient,
	}
}

// WithAPIKey sets the API key for authentication
func (c *Client) WithAPIKey(apiKey string) *Client {
	c.apiKey = apiKey
	return c
}

// WithTimeout sets a custom timeout for HTTP requests
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.httpClient.Timeout = timeout
	return c
}

// WithBaseURL sets a custom base URL for the API
func (c *Client) WithBaseURL(baseURL string) *Client {
	c.baseURL = strings.TrimSuffix(baseURL, "/")
	return c
}

// LookupResponse represents the complete response from the IPLocate API
type LookupResponse struct {
	IP           string   `json:"ip"`
	Country      *string  `json:"country"`
	CountryCode  *string  `json:"country_code"`
	IsEU         bool     `json:"is_eu"`
	City         *string  `json:"city"`
	Continent    *string  `json:"continent"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	TimeZone     *string  `json:"time_zone"`
	PostalCode   *string  `json:"postal_code"`
	Subdivision  *string  `json:"subdivision"`
	CurrencyCode *string  `json:"currency_code"`
	CallingCode  *string  `json:"calling_code"`
	Network      *string  `json:"network"`
	ASN          *ASN     `json:"asn"`
	Privacy      Privacy  `json:"privacy"`
	Company      *Company `json:"company"`
	Hosting      *Hosting `json:"hosting"`
	Abuse        *Abuse   `json:"abuse"`
}

// ASN represents Autonomous System Number information
type ASN struct {
	ASN         string `json:"asn"`
	Route       string `json:"route"`
	Netname     string `json:"netname"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
	Domain      string `json:"domain"`
	Type        string `json:"type"`
	RIR         string `json:"rir"`
}

// Privacy represents privacy and threat detection information
type Privacy struct {
	IsAbuser      bool `json:"is_abuser"`
	IsAnonymous   bool `json:"is_anonymous"`
	IsBogon       bool `json:"is_bogon"`
	IsHosting     bool `json:"is_hosting"`
	IsIcloudRelay bool `json:"is_icloud_relay"`
	IsProxy       bool `json:"is_proxy"`
	IsTor         bool `json:"is_tor"`
	IsVPN         bool `json:"is_vpn"`
}

// Company represents company information associated with the IP
type Company struct {
	Name        string `json:"name"`
	Domain      string `json:"domain"`
	CountryCode string `json:"country_code"`
	Type        string `json:"type"`
}

// Hosting represents hosting provider information
type Hosting struct {
	Provider *string `json:"provider"`
	Domain   *string `json:"domain"`
	Network  *string `json:"network"`
	Region   *string `json:"region"`
	Service  *string `json:"service"`
}

// Abuse represents abuse contact information
type Abuse struct {
	Address     *string `json:"address"`
	CountryCode *string `json:"country_code"`
	Email       *string `json:"email"`
	Name        *string `json:"name"`
	Network     *string `json:"network"`
	Phone       *string `json:"phone"`
}

// APIError represents an error response from the IPLocate API
type APIError struct {
	Message    string `json:"error"`
	StatusCode int    `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("IPLocate API error (%d): %s", e.StatusCode, e.Message)
}

// Lookup returns geolocation and threat intelligence data for the specified IP address
func (c *Client) Lookup(ip string) (*LookupResponse, error) {
	// Validate IP address format
	if parsedIP := net.ParseIP(ip); parsedIP == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ip)
	}

	endpoint := fmt.Sprintf("%s/lookup/%s", c.baseURL, url.PathEscape(ip))
	return c.doRequest(endpoint)
}

// LookupSelf returns geolocation and threat intelligence data for the client's current IP address
func (c *Client) LookupSelf() (*LookupResponse, error) {
	endpoint := fmt.Sprintf("%s/lookup/", c.baseURL)
	return c.doRequest(endpoint)
}

// doRequest performs the HTTP request to the IPLocate API
func (c *Client) doRequest(endpoint string) (*LookupResponse, error) {
	// Parse the endpoint URL to add query parameters
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoint URL: %w", err)
	}

	// Add API key as query parameter if provided
	if c.apiKey != "" {
		query := parsedURL.Query()
		query.Set("apikey", c.apiKey)
		parsedURL.RawQuery = query.Encode()
	}

	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "go-iplocate/1.0.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-200 status codes
	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			// If we can't parse the error response, return the raw body
			return nil, fmt.Errorf("API request failed (%d): %s", resp.StatusCode, string(body))
		}
		apiErr.StatusCode = resp.StatusCode
		return nil, &apiErr
	}

	var result LookupResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}
