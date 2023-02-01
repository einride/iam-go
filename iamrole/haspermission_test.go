package iamrole

import (
	"testing"

	"cloud.google.com/go/iam/admin/apiv1/adminpb"
	"gotest.tools/v3/assert"
)

func TestHasPermission(t *testing.T) {
	for _, tt := range []struct {
		name       string
		role       *adminpb.Role
		permission string
		expected   bool
	}{
		{
			name: "via service wildcard",
			role: &adminpb.Role{
				IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
			},
			permission: "pubsub.subscriptions.consume",
			expected:   true,
		},

		{
			name: "via resource wildcard",
			role: &adminpb.Role{
				IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
			},
			permission: "storage.buckets.create",
			expected:   true,
		},

		{
			name: "via exact match",
			role: &adminpb.Role{
				IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
			},
			permission: "kms.keys.create",
			expected:   true,
		},

		{
			name: "no match",
			role: &adminpb.Role{
				IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
			},
			permission: "kms.keys.get",
			expected:   false,
		},

		{
			name: "no match with wildcard",
			role: &adminpb.Role{
				IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
			},
			permission: "storage.objects.get",
			expected:   false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, HasPermission(tt.role, tt.permission))
		})
	}
}
