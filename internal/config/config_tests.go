package config

import (
	"os"
	"testing"
)

// TestLoadConfigValue tests the LoadConfig function when the environment variable is set.
func TestLoadConfigValue(t *testing.T) {
	// Set the environment variable
	expectedAPIKey := "test-api-key"
	if err := os.Setenv("NEWSFETCHER_GOOGLE_API_KEY", expectedAPIKey); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("NEWSFETCHER_GOOGLE_API_KEY"); err != nil {
			t.Fatalf("Failed to unset environment variable: %v", err)
		}
	}()

	// Call the LoadConfig function
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if the GoogleAPIKey is set correctly
	if config.GoogleAPIKey != expectedAPIKey {
		t.Errorf("Expected GoogleAPIKey to be %s, got %s", expectedAPIKey, config.GoogleAPIKey)
	}
}

// TestLoadConfigIsMissing tests the LoadConfig function when the environment variable is not set.
func TestLoadConfigIsMissing(t *testing.T) {
	// Unset the environment variable
	defer func() {
		if err := os.Unsetenv("NEWSFETCHER_GOOGLE_API_KEY"); err != nil {
			t.Fatalf("Failed to unset environment variable: %v", err)
		}
	}()

	// Call the LoadConfig function
	_, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
