package iamfirebase

import (
	"strings"

	"go.einride.tech/iam/iamjwt"
)

// ProjectID returns the token payload's Firebase project ID.
func ProjectID(token iamjwt.Token) string {
	return strings.TrimPrefix(token.Issuer, Issuer+"/")
}
