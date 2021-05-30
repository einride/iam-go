package iamfirebase

import (
	"strings"

	"go.einride.tech/iam/iamjwt"
)

// ProjectID returns the token payload's Firebase project ID.
func ProjectID(p iamjwt.Payload) string {
	return strings.TrimPrefix(p.Issuer, Issuer+"/")
}
