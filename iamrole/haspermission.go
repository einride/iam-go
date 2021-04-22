package iamrole

import (
	"go.einride.tech/iam/iampermission"
	"google.golang.org/genproto/googleapis/iam/admin/v1"
)

// HasPermission reports whether the provided role has the provided permission.
// Always returns false for wildcard permissions.
func HasPermission(role *admin.Role, permission string) bool {
	for _, includedPermission := range role.IncludedPermissions {
		if iampermission.Match(includedPermission, permission) {
			return true
		}
	}
	return false
}
