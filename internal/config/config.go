package config

import (
	"os"
)

// GetBaseURL returns the backend URL from environment or default
func GetBaseURL() string {
	if url := os.Getenv("CRM_BACKEND_URL"); url != "" {
		return url
	}
	return "http://localhost:8082/api"
}

// GetAdminAPIKey returns the admin API key from environment or default
func GetAdminToken() string {
	if key := os.Getenv("CRM_ADMIN_API_KEY"); key != "" {
		return key
	}
	// Default key for development - NEVER use this in production!
	return ""
}
