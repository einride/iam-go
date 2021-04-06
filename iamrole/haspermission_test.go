package iamrole

import (
	"testing"

	"google.golang.org/genproto/googleapis/iam/admin/v1"
	"gotest.tools/v3/assert"
)

func TestHasPermission(t *testing.T) {
	for _, tt := range []struct {
		name       string
		role       *admin.Role
		permission string
		expected   bool
	}{
		{
			name: "via service wildcard",
			role: &admin.Role{
				IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
			},
			permission: "pubsub.subscriptions.consume",
			expected:   true,
		},

		{
			name: "via resource wildcard",
			role: &admin.Role{
				IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
			},
			permission: "storage.buckets.create",
			expected:   true,
		},

		{
			name: "via exact match",
			role: &admin.Role{
				IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
			},
			permission: "kms.keys.create",
			expected:   true,
		},

		{
			name: "no match",
			role: &admin.Role{
				IncludedPermissions: []string{"pubsub.*", "storage.buckets.*", "kms.keys.create"},
			},
			permission: "kms.keys.get",
			expected:   false,
		},

		{
			name: "no match with wildcard",
			role: &admin.Role{
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
