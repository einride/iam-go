package iamgoogle

import (
	"go.einride.tech/iam/iamjwt"
)

// Issuer is the issuer of Google ID tokens.
const Issuer = "https://accounts.google.com"

// IsGoogleIDToken returns true if the JWT payload is from a Google ID token.
// See: https://developers.google.com/identity/protocols/oauth2/openid-connect
func IsGoogleIDToken(p iamjwt.Payload) bool {
	return p.Issuer == Issuer
}
