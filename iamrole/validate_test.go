package iamrole

import (
	"strings"
	"testing"

	"google.golang.org/genproto/googleapis/iam/admin/v1"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/testing/protocmp"
	"gotest.tools/v3/assert"
)

func TestValidate(t *testing.T) {
	for _, tt := range []struct {
		name     string
		role     *admin.Role
		expected *errdetails.BadRequest
	}{
		{
			name: "valid",
			role: &admin.Role{
				Name:                "roles/foo.barBaz",
				Title:               "Foo Bar Baz",
				Description:         "Longer description",
				IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
			},
		},

		{
			name: "invalid name format",
			role: &admin.Role{
				Name:                "foobarbaz",
				Title:               "Foo Bar Baz",
				Description:         "Longer description",
				IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
			},
			expected: &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "name",
						Description: "must have format `roles/{role}`",
					},
				},
			},
		},

		{
			name: "too long title",
			role: &admin.Role{
				Name:                "roles/foo.barBaz",
				Title:               strings.Repeat("a", 101),
				Description:         "Longer description",
				IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
			},
			expected: &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "title",
						Description: "must be non-empty and <= 100 characters",
					},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.DeepEqual(t, tt.expected, Validate(tt.role), protocmp.Transform())
		})
	}
}
