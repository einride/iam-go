package iamrole

import (
	"cloud.google.com/go/iam/admin/apiv1/adminpb"
	"go.einride.tech/iam/iampermission"
)

// HasPermission reports whether the provided role has the provided permission.
// Always returns false for wildcard permissions.
func HasPermission(role *adminpb.Role, permission string) bool {
	for _, includedPermission := range role.IncludedPermissions {
		if iampermission.Match(includedPermission, permission) {
			return true
		}
	}
	return false
}
