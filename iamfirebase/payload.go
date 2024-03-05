package iamfirebase

import (
	"strings"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// ProjectID returns the token payload's Firebase project ID.
func ProjectID(token *iamv1.IdentityToken) string {
	return strings.TrimPrefix(token.GetIss(), Issuer+"/")
}
