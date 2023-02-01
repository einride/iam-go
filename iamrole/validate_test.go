package iamrole

import (
	"strings"
	"testing"

	"cloud.google.com/go/iam/admin/apiv1/adminpb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"gotest.tools/v3/assert"
)

func TestValidate(t *testing.T) {
	for _, tt := range []struct {
		name     string
		role     *adminpb.Role
		expected *errdetails.BadRequest
	}{
		{
			name: "valid",
			role: &adminpb.Role{
				Name:                "roles/foo.barBaz",
				Title:               "Foo Bar Baz",
				Description:         "Longer description",
				IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
			},
		},

		{
			name: "invalid name format",
			role: &adminpb.Role{
				Name:                "foobarbaz",
				Title:               "Foo Bar Baz",
				Description:         "Longer description",
				IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
			},
			expected: &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "name",
						Description: "role name 'foobarbaz' is not on the format `roles/{service}.{role}`",
					},
				},
			},
		},

		{
			name: "too long title",
			role: &adminpb.Role{
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
			err := Validate(tt.role)
			if tt.expected == nil {
				assert.NilError(t, err)
			} else {
				actual, ok := status.Convert(err).Details()[0].(*errdetails.BadRequest)
				assert.Assert(t, ok)
				assert.DeepEqual(t, tt.expected, actual, protocmp.Transform())
			}
		})
	}
}
