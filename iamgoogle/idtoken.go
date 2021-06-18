package iamgoogle

import (
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// Issuer is the issuer of Google ID tokens.
const Issuer = "https://accounts.google.com"

// IsGoogleIdentityToken returns true if the JWT payload is from a Google ID token.
// See: https://developers.google.com/identity/protocols/oauth2/openid-connect
func IsGoogleIdentityToken(token *iamv1.IdentityToken) bool {
	return token.Iss == Issuer
}
