package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// CloudflareIPRanges represents the structure of the response from Cloudflare's IP ranges API
type CloudflareIPRanges struct {
	Result struct {
		IPv4CIDRs []string `json:"ipv4_cidrs"`
		IPv6CIDRs []string `json:"ipv6_cidrs"`
	} `json:"result"`
}

// FetchCloudflareIPv4Ranges fetches the Cloudflare IP ranges from the Cloudflare API
func FetchCloudflareIPv4Ranges(url string) ([]string, error) {
	if url == "" {
		return nil, fmt.Errorf("url is empty")
	}
	resp, err := MakeSecureGetHTTPRequest(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Cloudflare IP ranges: %v", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			// Handle the error if needed, for example, log it
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var ipRanges CloudflareIPRanges
	if err = json.Unmarshal(body, &ipRanges); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return ipRanges.Result.IPv4CIDRs, nil
}
