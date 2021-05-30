package iamfirebase

import (
	"go.einride.tech/iam/iamjwt"
)

// Issuer is the issuer of Firebase ID tokens.
const Issuer = "https://securetoken.google.com"

// IsFirebaseIDToken returns true if the JWT payload is from a Firebase ID token.
// See: https://firebase.google.com/docs/rules/rules-and-auth
func IsFirebaseIDToken(p iamjwt.Payload) bool {
	return p.Issuer == Issuer+"/"+p.Audience
}
