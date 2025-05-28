package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/iplocate/go-iplocate"
)

func printSection(title string) {
	fmt.Printf("\n=== %s ===\n", title)
}

func printStringField(label string, value *string) {
	if value != nil {
		fmt.Printf("%s: %s\n", label, *value)
	} else {
		fmt.Printf("%s: <not available>\n", label)
	}
}

func printFloatField(label string, value *float64) {
	if value != nil {
		fmt.Printf("%s: %.4f\n", label, *value)
	} else {
		fmt.Printf("%s: <not available>\n", label)
	}
}

func main() {
	// Create a new client
	client := iplocate.NewClient(nil)

	// Or add an API key for higher rate limits
	if apiKey := os.Getenv("IPLOCATE_API_KEY"); apiKey != "" {
		client = client.WithAPIKey(apiKey)
		fmt.Println("Using API key for enhanced rate limits")
	} else {
		fmt.Println("No API key provided - using free tier (50 requests/day)")
		fmt.Println("Set IPLOCATE_API_KEY environment variable for 1,000 requests/day")
	}

	// Look up an IP address
	fmt.Println("\n🔍 Looking up 8.8.8.8 (Google DNS)...")
	result, err := client.Lookup("8.8.8.8")
	if err != nil {
		log.Fatalf("Error looking up IP: %v", err)
	}

	// Print all available data in organized sections
	printSection("Basic Information")
	fmt.Printf("IP Address: %s\n", result.IP)
	printStringField("Country", result.Country)
	printStringField("Country Code", result.CountryCode)
	fmt.Printf("Is EU: %v\n", result.IsEU)
	printStringField("City", result.City)
	printStringField("Continent", result.Continent)
	printStringField("Subdivision", result.Subdivision)
	printStringField("Postal Code", result.PostalCode)
	printStringField("Time Zone", result.TimeZone)
	printStringField("Currency Code", result.CurrencyCode)
	printStringField("Calling Code", result.CallingCode)

	printSection("Geographic Coordinates")
	printFloatField("Latitude", result.Latitude)
	printFloatField("Longitude", result.Longitude)
	if result.Latitude != nil && result.Longitude != nil {
		fmt.Printf("Google Maps: https://maps.google.com/?q=%.4f,%.4f\n", *result.Latitude, *result.Longitude)
	}

	printSection("Network Information")
	printStringField("Network", result.Network)
	if result.ASN != nil {
		fmt.Printf("ASN: %s\n", result.ASN.ASN)
		fmt.Printf("ASN Name: %s\n", result.ASN.Name)
		fmt.Printf("Route: %s\n", result.ASN.Route)
		fmt.Printf("Netname: %s\n", result.ASN.Netname)
		fmt.Printf("Domain: %s\n", result.ASN.Domain)
		fmt.Printf("Type: %s\n", result.ASN.Type)
		fmt.Printf("RIR: %s\n", result.ASN.RIR)
		fmt.Printf("ASN Country: %s\n", result.ASN.CountryCode)
	} else {
		fmt.Println("ASN information: <not available>")
	}

	printSection("Privacy & Threat Detection")
	fmt.Printf("Is Abuser: %v\n", result.Privacy.IsAbuser)
	fmt.Printf("Is Anonymous: %v\n", result.Privacy.IsAnonymous)
	fmt.Printf("Is Bogon: %v\n", result.Privacy.IsBogon)
	fmt.Printf("Is Hosting: %v\n", result.Privacy.IsHosting)
	fmt.Printf("Is iCloud Relay: %v\n", result.Privacy.IsIcloudRelay)
	fmt.Printf("Is Proxy: %v\n", result.Privacy.IsProxy)
	fmt.Printf("Is Tor: %v\n", result.Privacy.IsTor)
	fmt.Printf("Is VPN: %v\n", result.Privacy.IsVPN)

	printSection("Company Information")
	if result.Company != nil {
		fmt.Printf("Company Name: %s\n", result.Company.Name)
		fmt.Printf("Company Domain: %s\n", result.Company.Domain)
		fmt.Printf("Company Country: %s\n", result.Company.CountryCode)
		fmt.Printf("Company Type: %s\n", result.Company.Type)
	} else {
		fmt.Println("Company information: <not available>")
	}

	printSection("Hosting Information")
	if result.Hosting != nil {
		printStringField("Provider", result.Hosting.Provider)
		printStringField("Domain", result.Hosting.Domain)
		printStringField("Network", result.Hosting.Network)
		printStringField("Region", result.Hosting.Region)
		printStringField("Service", result.Hosting.Service)
	} else {
		fmt.Println("Hosting information: <not available>")
	}

	printSection("Abuse Contact Information")
	if result.Abuse != nil {
		printStringField("Name", result.Abuse.Name)
		printStringField("Email", result.Abuse.Email)
		printStringField("Phone", result.Abuse.Phone)
		printStringField("Address", result.Abuse.Address)
		printStringField("Network", result.Abuse.Network)
		printStringField("Country", result.Abuse.CountryCode)
	} else {
		fmt.Println("Abuse contact information: <not available>")
	}

	printSection("Raw JSON Response")
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
	} else {
		fmt.Println(string(jsonData))
	}

	// Look up your own IP
	printSection("Your Current IP Information")
	selfResult, err := client.LookupSelf()
	if err != nil {
		log.Fatalf("Error looking up own IP: %v", err)
	}

	fmt.Printf("Your IP: %s\n", selfResult.IP)
	if selfResult.Country != nil {
		fmt.Printf("Your Country: %s", *selfResult.Country)
		if selfResult.CountryCode != nil {
			fmt.Printf(" (%s)", *selfResult.CountryCode)
		}
		fmt.Println()
	}
	if selfResult.City != nil {
		fmt.Printf("Your City: %s\n", *selfResult.City)
	}

	// Privacy flags for your IP
	fmt.Println("\nYour IP Privacy Status:")
	fmt.Printf("  🔒 VPN: %v\n", selfResult.Privacy.IsVPN)
	fmt.Printf("  🌐 Proxy: %v\n", selfResult.Privacy.IsProxy)
	fmt.Printf("  🧅 Tor: %v\n", selfResult.Privacy.IsTor)
	fmt.Printf("  🏢 Hosting: %v\n", selfResult.Privacy.IsHosting)
	if selfResult.Privacy.IsVPN || selfResult.Privacy.IsProxy || selfResult.Privacy.IsTor {
		fmt.Println("  ⚠️  Privacy tools detected!")
	} else {
		fmt.Println("  ✅ Direct connection detected")
	}
}
