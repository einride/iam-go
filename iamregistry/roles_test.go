package iamregistry

import (
	"testing"

	"cloud.google.com/go/iam/admin/apiv1/adminpb"
	"gotest.tools/v3/assert"
)

func TestRoles_RangeRolesByPermission(t *testing.T) {
	t.Run("wildcard", func(t *testing.T) {
		roles, err := NewRoles(
			&adminpb.Role{
				Name:                "roles/test.admin",
				Title:               "Test admin",
				Description:         "Test description",
				IncludedPermissions: []string{"test.*"},
			},
		)
		assert.NilError(t, err)
		var found bool
		roles.RangeRolesByPermission("test.foo.bar", func(role *adminpb.Role) bool {
			assert.Equal(t, "roles/test.admin", role.GetName())
			found = true
			return true
		})
		assert.Assert(t, found)
	})
}
