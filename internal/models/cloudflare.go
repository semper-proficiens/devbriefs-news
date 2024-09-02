package models

// CloudflareIPRanges represents the structure of the response from Cloudflare's IP ranges API
type CloudflareIPRanges struct {
	Result struct {
		IPv4CIDRs []string `json:"ipv4_cidrs"`
		IPv6CIDRs []string `json:"ipv6_cidrs"`
	} `json:"result"`
}
