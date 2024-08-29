package config

import (
	"fmt"
	"os"
)

// Config holds the configuration settings for the application.
type Config struct {
	GoogleAPIKey string // The API key for accessing the Google News API
}

// LoadConfig loads the configuration settings from environment variables.
func LoadConfig() (*Config, error) {
	// Read the Google API key from the environment variable
	googleAPIKey := os.Getenv("NEWSFETCHER_GOOGLE_API_KEY")

	// Check if the required environment variable is set
	if googleAPIKey == "" {
		return nil, fmt.Errorf("environment variable NEWSFETCHER_GOOGLE_API_KEY is required but not set")
	}

	// Create a Config struct and populate it with values from environment variables
	config := &Config{
		GoogleAPIKey: googleAPIKey,
	}

	return config, nil
}
