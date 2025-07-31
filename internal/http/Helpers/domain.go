package Helpers

import (
	"net/http"
	"strings"
)

// GetDomain extracts the domain from the HTTP request.
func GetDomain(r *http.Request) string {
	if r.Header.Get("X-Forwarded-Host") != "" {
		return r.Header.Get("X-Forwarded-Host")
	}

	if r.Host != "" {
		return r.Host
	}

	return ""
}

func GetSubdomain(r *http.Request) string {
	fullDomain := GetDomain(r)
	if fullDomain == "" {
		return ""
	}

	// Split the domain by '.' and return the first part as the subdomain
	parts := strings.Split(fullDomain, ".")
	if len(parts) > 2 {
		return parts[0] // Return the first part as subdomain
	}

	return "" // No subdomain found
}
