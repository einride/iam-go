package iamgoogle

import (
	"strings"

	"google.golang.org/api/idtoken"
)

// IsEmailVerified returns true if the payload has the `email_verified` claim.
func IsEmailVerified(payload *idtoken.Payload) bool {
	result, ok := payload.Claims["email_verified"].(bool)
	return result && ok
}

// Email returns the payload's `email` claim.
func Email(payload *idtoken.Payload) (string, bool) {
	result, ok := payload.Claims["email"].(string)
	return result, ok
}

// HostedDomain returns the payload's `hd` claim.
func HostedDomain(payload *idtoken.Payload) (string, bool) {
	result, ok := payload.Claims["hd"].(string)
	return result, ok
}

// IsGoogleCloudServiceAccountEmail returns true if the payload has a verified email belonging to a Google Cloud
// service account.
func IsGoogleCloudServiceAccountEmail(payload *idtoken.Payload) bool {
	if !IsEmailVerified(payload) {
		return false
	}
	email, ok := Email(payload)
	return ok && strings.HasSuffix(email, ".gserviceaccount.com")
}
