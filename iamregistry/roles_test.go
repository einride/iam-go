package iamregistry

import (
	"testing"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/genproto/googleapis/iam/admin/v1"
	"gotest.tools/v3/assert"
)

func TestRoles_RangeRolesByPermission(t *testing.T) {
	t.Run("wildcard", func(t *testing.T) {
		roles, err := NewRoles(&iamv1.Roles{
			Role: []*admin.Role{
				{
					Name:                "roles/test.admin",
					Title:               "Test admin",
					Description:         "Test description",
					IncludedPermissions: []string{"test.*"},
				},
			},
		})
		assert.NilError(t, err)
		var found bool
		roles.RangeRolesByPermission("test.foo.bar", func(role *admin.Role) bool {
			assert.Equal(t, "roles/test.admin", role.Name)
			found = true
			return true
		})
		assert.Assert(t, found)
	})
}
