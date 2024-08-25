package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// NewSecureHTTPClient creates a new HTTP client with secure settings.
func NewSecureHTTPClient() *http.Client {
	// Load system CA certificates
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		rootCAs = x509.NewCertPool()
	}

	// Create a custom TLS configuration
	tlsConfig := &tls.Config{
		RootCAs:          rootCAs,
		MinVersion:       tls.VersionTLS12, // Use TLS 1.2 or higher
		CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}

	// Create an HTTP transport with the custom TLS configuration
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		// Enable HTTP/2
		ForceAttemptHTTP2: true,
		// Set other transport settings
		MaxIdleConns:       100,
		IdleConnTimeout:    90 * time.Second,
		DisableCompression: false,
	}

	// Create and return the HTTP client with the custom transport
	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second, // Set a timeout for the client
	}
}

// MakeSecureHTTPRequest is a utility function for making secure HTTP requests.
// It takes the HTTP method, URL, and request body as parameters and returns the HTTP response.
func MakeSecureHTTPRequest(method, url string, body io.Reader) (*http.Response, error) {
	// Create a new secure HTTP client
	client := NewSecureHTTPClient()

	// Create a new HTTP request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Printf("Error creating new http request: %s", err)
		return nil, err
	}

	// Set common headers (if any)
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making http request in util pkg %s", err)
		return nil, err
	}

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Received a non 200 status code from server: %s", resp.Status)
		return nil, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	return resp, nil
}
