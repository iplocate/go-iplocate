# IPLocate geolocation client for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/iplocate/go-iplocate.svg)](https://pkg.go.dev/github.com/iplocate/go-iplocate)
[![Go Report Card](https://goreportcard.com/badge/github.com/iplocate/go-iplocate)](https://goreportcard.com/report/github.com/iplocate/go-iplocate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go client for the [IPLocate.io](https://iplocate.io) geolocation API. Look up detailed geolocation and threat intelligence data for any IP address:

- **IP geolocation**: IP to country, IP to city, IP to region/state, coordinates, timezone, postal code
- **ASN information**: Internet service provider, network details, routing information  
- **Privacy & threat detection**: VPN, proxy, Tor, hosting provider detection
- **Company information**: Business details associated with IP addresses - company name, domain, type (ISP/hosting/education/government/business)
- **Abuse contact**: Network abuse reporting information
- **Hosting detection**: Cloud provider and hosting service detection using our proprietary hosting detection engine

See what information we can provide for [your IP address](https://iplocate.io/what-is-my-ip).

## Getting started

You can make 1,000 free requests per day with a [free account](https://iplocate.io/signup). For higher plans, check out [API pricing](https://www.iplocate.io/pricing).

### Installation

```bash
go get github.com/iplocate/go-iplocate
```

### Quick start

```go
package main

import (
    "fmt"
    "log"

    "github.com/iplocate/go-iplocate"
)

func main() {
    // Create a new client with default HTTP client
    // Get your free API key from https://iplocate.io/signup
    client := iplocate.NewClient(nil).WithAPIKey("your-api-key")

    // Look up an IP address
    result, err := client.Lookup("8.8.8.8")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("IP: %s\n", result.IP)
    if result.Country != nil {
        fmt.Printf("Country: %s\n", *result.Country)
    }
    if result.City != nil {
        fmt.Printf("City: %s\n", *result.City)
    }
    
    // Check privacy flags
    fmt.Printf("Is VPN: %v\n", result.Privacy.IsVPN)
    fmt.Printf("Is Proxy: %v\n", result.Privacy.IsProxy)
}
```

### Get the country for an IP address

```go
fmt.Printf("Country: %s (%s)\n", *result.Country, *result.CountryCode)
```

### Get the currency code for a country by IP address

```go
fmt.Printf("Currency: %s\n", *result.CurrencyCode)
```

### Get the calling code for a country by IP address

```go
fmt.Printf("Calling code: %s\n", *result.CallingCode)
```

## Authentication

Get your free API key from [IPLocate.io](https://iplocate.io/signup), and pass it to `.WithAPIKey()`:

```go
client := iplocate.NewClient(nil).WithAPIKey("your-api-key")
```

## Examples

### IP address geolocation lookup

```go
client := iplocate.NewClient(nil).WithAPIKey("your-api-key")
result, err := client.Lookup("203.0.113.1")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Country: %s (%s)\n", *result.Country, *result.CountryCode)
fmt.Printf("Coordinates: %.4f, %.4f\n", *result.Latitude, *result.Longitude)
```

### Get your own IP address information

```go
result, err := client.LookupSelf().WithAPIKey("your-api-key")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Your IP: %s\n", result.IP)
```

### Check for VPN/Proxy Detection

```go
result, err := client.Lookup("192.0.2.1")
if err != nil {
    log.Fatal(err)
}

if result.Privacy.IsVPN {
    fmt.Println("This IP is using a VPN")
}

if result.Privacy.IsProxy {
    fmt.Println("This IP is using a proxy")
}

if result.Privacy.IsTor {
    fmt.Println("This IP is using Tor")
}
```

### ASN and network information

```go
result, err := client.Lookup("8.8.8.8")
if err != nil {
    log.Fatal(err)
}

if result.ASN != nil {
    fmt.Printf("ASN: %s\n", result.ASN.ASN)
    fmt.Printf("ISP: %s\n", result.ASN.Name)
    fmt.Printf("Network: %s\n", result.ASN.Route)
}
```

### Custom configuration

```go
// With custom HTTP client
customHTTPClient := &http.Client{Timeout: 60 * time.Second}
client := iplocate.NewClient(customHTTPClient).
    WithAPIKey("your-api-key").
    WithBaseURL("https://custom-endpoint.com")

// Or with default client and custom timeout
client := iplocate.NewClient(nil).
    WithAPIKey("your-api-key").
    WithTimeout(60 * time.Second)
```

## Response structure

The `LookupResponse` struct contains all available data:

```go
type LookupResponse struct {
    IP           string    `json:"ip"`
    Country      *string   `json:"country"`
    CountryCode  *string   `json:"country_code"`
    IsEU         bool      `json:"is_eu"`
    City         *string   `json:"city"`
    Continent    *string   `json:"continent"`
    Latitude     *float64  `json:"latitude"`
    Longitude    *float64  `json:"longitude"`
    TimeZone     *string   `json:"time_zone"`
    PostalCode   *string   `json:"postal_code"`
    Subdivision  *string   `json:"subdivision"`
    CurrencyCode *string   `json:"currency_code"`
    CallingCode  *string   `json:"calling_code"`
    Network      *string   `json:"network"`
    ASN          *ASN      `json:"asn"`
    Privacy      Privacy   `json:"privacy"`
    Company      *Company  `json:"company"`
    Hosting      *Hosting  `json:"hosting"`
    Abuse        *Abuse    `json:"abuse"`
}
```

Note: Fields marked with `*` are pointers and may be `nil` if data is not available.

## Error handling

```go
result, err := client.Lookup("invalid-ip")
if err != nil {
    if apiErr, ok := err.(*iplocate.APIError); ok {
        fmt.Printf("API Error (%d): %s\n", apiErr.StatusCode, apiErr.Error)
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
}
```

Common API errors:

- `400 Bad Request`: Invalid IP address format
- `403 Forbidden`: Invalid API key
- `404 Not Found`: IP address not found
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

## API reference

For complete API documentation, visit [iplocate.io/docs](https://iplocate.io/docs).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Testing

Run the test suite:

```bash
go test -v
```

Run tests with coverage:

```bash
go test -v -cover
```

## About IPLocate.io

Since 2017, IPLocate has set out to provide the most reliable and accurate IP address data.

We process 50TB+ of data to produce our comprehensive IP geolocation, IP to company, proxy and VPN detection, hosting detection, ASN, and WHOIS data sets. Our API handles over 15 billion requests a month for thousands of businesses and developers.

- Email: [support@iplocate.io](mailto:support@iplocate.io)
- Website: [iplocate.io](https://iplocate.io)
- Documentation: [iplocate.io/docs](https://iplocate.io/docs)
- Sign up for a free API Key: [iplocate.io/signup](https://iplocate.io/signup)
