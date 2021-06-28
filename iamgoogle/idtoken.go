package iamgoogle

import (
	"strings"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// IsGoogleIdentityToken returns true if the JWT payload is from a Google ID token.
// See: https://developers.google.com/identity/protocols/oauth2/openid-connect
func IsGoogleIdentityToken(token *iamv1.IdentityToken) bool {
	return strings.TrimPrefix(token.Iss, "https://") == "accounts.google.com"
}
