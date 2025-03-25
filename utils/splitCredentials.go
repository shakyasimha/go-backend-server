package utils

import "strings"

// Helper function to split username:password
func SplitCredentials(creds string) (string, string) {
	parts := strings.SplitN(creds, ":", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
