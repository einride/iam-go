package iamreflect

import (
	"testing"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/genproto/googleapis/iam/admin/v1"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"gotest.tools/v3/assert"
)

func TestValidateRoles(t *testing.T) {
	for _, tt := range []struct {
		name     string
		roles    *iamv1.Roles
		expected *errdetails.BadRequest
	}{
		{
			name: "valid",
			roles: &iamv1.Roles{
				Role: []*admin.Role{
					{
						Name:                "roles/foo.barBaz",
						Title:               "Foo Bar Baz",
						Description:         "Longer description",
						IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
					},
					{
						Name:                "roles/baz.fooBar",
						Title:               "Baz Foo Bar",
						Description:         "Longer description",
						IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
					},
				},
			},
		},

		{
			name: "invalid name format",
			roles: &iamv1.Roles{
				Role: []*admin.Role{
					{
						Name:                "foobarbaz",
						Title:               "Foo Bar Baz",
						Description:         "Longer description",
						IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
					},
					{
						Name:                "roles/baz.fooBar",
						Title:               "Baz Foo Bar",
						Description:         "Longer description",
						IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
					},
				},
			},
			expected: &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "role[0].name",
						Description: "must have format `roles/{service}.{role}`",
					},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRoles(tt.roles)
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
