package iamfirebase

import (
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// Issuer is the issuer of Firebase ID tokens.
const Issuer = "https://securetoken.google.com"

// IsFirebaseIdentityToken returns true if the JWT payload is from a Firebase ID token.
// See: https://firebase.google.com/docs/rules/rules-and-auth
func IsFirebaseIdentityToken(token *iamv1.IdentityToken) bool {
	return token.GetIss() == Issuer+"/"+token.GetAud()
}
