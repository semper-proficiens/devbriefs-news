package utils

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestNewSecureHTTPClient tests the NewSecureHTTPClient function.
func TestNewSecureHTTPClient(t *testing.T) {
	client := NewSecureHTTPClient()

	// Check if the client is not nil
	if client == nil {
		t.Fatal("Expected non-nil HTTP client")
	}

	// Check if the transport is of type *http.Transport
	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("Expected transport to be of type *http.Transport")
	}

	// Check if the TLSClientConfig is set correctly
	if transport.TLSClientConfig == nil {
		t.Fatal("Expected non-nil TLSClientConfig")
	}

	// Check if the MinVersion is set to TLS 1.2
	if transport.TLSClientConfig.MinVersion != tls.VersionTLS12 {
		t.Fatalf("Expected MinVersion to be TLS 1.2, got %v", transport.TLSClientConfig.MinVersion)
	}

	// Check if the Timeout is set correctly
	if client.Timeout != 10*time.Second {
		t.Fatalf("Expected Timeout to be 10 seconds, got %v", client.Timeout)
	}
}

// TestMakeSecureGetHTTPRequest tests the MakeSecureGetHTTPRequest function.
func TestMakeSecureGetHTTPRequest(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer ts.Close()

	// Make a secure GET request to the test server
	resp, err := MakeSecureGetHTTPRequest(ts.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	// Check if the status code is 200
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}
}

// TestMakeSecureGetHTTPRequest_Non200 tests the MakeSecureGetHTTPRequest function for non-200 status codes.
func TestMakeSecureGetHTTPRequest_Non200(t *testing.T) {
	// Create a test server that returns a 404 status code
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	// Make a secure GET request to the test server
	resp, err := MakeSecureGetHTTPRequest(ts.URL)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check if the response is nil
	if resp != nil {
		t.Fatalf("Expected nil response, got %v", resp)
	}
}
