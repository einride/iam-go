package iamgoogle

import (
	"testing"

	"google.golang.org/api/idtoken"
	"gotest.tools/v3/assert"
)

func TestContextMemberResolver_IsGoogleServiceAccount(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{
			email: "default-compute@developer.gserviceaccount.com",
			valid: true,
		},
		{
			email: "user-managed@einride.iam.gserviceaccount.com",
			valid: true,
		},
		{
			email: "google-managed@@cloudservices.gserviceaccount.com",
			valid: true,
		},
		{
			email: "missing-dot@gserviceaccount.com",
			valid: false,
		},
		{
			email: "any@example.com",
			valid: false,
		},
	}

	for _, tt := range tests {
		claims := map[string]interface{}{
			"email_verified": true,
			"email":          tt.email,
		}

		assert.Equal(t, IsGoogleCloudServiceAccountEmail(&idtoken.Payload{
			Claims: claims,
		}), tt.valid)
	}
}
