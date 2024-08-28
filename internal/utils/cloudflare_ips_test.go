package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock response for the Cloudflare IP ranges API
var mockCloudflareResponse = CloudflareIPRanges{
	Result: struct {
		IPv4CIDRs []string `json:"ipv4_cidrs"`
		IPv6CIDRs []string `json:"ipv6_cidrs"`
	}{
		IPv4CIDRs: []string{"192.0.2.0/24", "198.51.100.0/24"},
		IPv6CIDRs: []string{"2001:db8::/32"},
	},
}

// TestFetchCloudflareIPv4Ranges tests the FetchCloudflareIPv4Ranges function.
func TestFetchCloudflareIPv4Ranges(t *testing.T) {
	// Create a test server that returns the mock response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockCloudflareResponse)
	}))
	defer ts.Close()

	// Call the FetchCloudflareIPv4Ranges function with the test server URL
	ipv4Ranges, err := FetchCloudflareIPv4Ranges(ts.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if the returned IPv4 ranges match the mock response
	expectedIPv4Ranges := mockCloudflareResponse.Result.IPv4CIDRs
	if len(ipv4Ranges) != len(expectedIPv4Ranges) {
		t.Fatalf("Expected %d IPv4 ranges, got %d", len(expectedIPv4Ranges), len(ipv4Ranges))
	}
	for i, cidr := range ipv4Ranges {
		if cidr != expectedIPv4Ranges[i] {
			t.Fatalf("Expected IPv4 range %s, got %s", expectedIPv4Ranges[i], cidr)
		}
	}
}

// TestFetchCloudflareIPv4Ranges_Non200 tests the FetchCloudflareIPv4Ranges function for non-200 status codes.
func TestFetchCloudflareIPv4Ranges_Non200(t *testing.T) {
	// Create a test server that returns a 500 status code
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	// Call the FetchCloudflareIPv4Ranges function with the test server URL
	_, err := FetchCloudflareIPv4Ranges(ts.URL)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// TestFetchCloudflareIPv4Ranges_InvalidJSON tests the FetchCloudflareIPv4Ranges function for invalid JSON response.
func TestFetchCloudflareIPv4Ranges_InvalidJSON(t *testing.T) {
	// Create a test server that returns invalid JSON
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer ts.Close()

	// Call the FetchCloudflareIPv4Ranges function with the test server URL
	_, err := FetchCloudflareIPv4Ranges(ts.URL)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
