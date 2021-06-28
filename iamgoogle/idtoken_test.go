package iamgoogle

import (
	"testing"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"gotest.tools/v3/assert"
)

func TestIsGoogleIdentityToken(t *testing.T) {
	assert.Assert(t, IsGoogleIdentityToken(&iamv1.IdentityToken{Iss: "https://accounts.google.com"}))
	assert.Assert(t, IsGoogleIdentityToken(&iamv1.IdentityToken{Iss: "accounts.google.com"}))
	assert.Assert(t, !IsGoogleIdentityToken(&iamv1.IdentityToken{Iss: "foo.com"}))
	assert.Assert(t, !IsGoogleIdentityToken(&iamv1.IdentityToken{Iss: "http://accounts.google.com"}))
}
