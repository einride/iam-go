package iamreflect

import (
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iamresource"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// ResolveResourcePermission resolves a permission for a resource name, given a set of resource permissions.
func ResolveResourcePermission(
	resourcePermissions []*iamv1.ResourcePermission,
	resourceName string,
) (string, bool) {
	for _, resourcePermission := range resourcePermissions {
		if resourcePermission.GetResource().GetType() == iamresource.Root && resourceName == iamresource.Root {
			return resourcePermission.GetPermission(), true
		}
		for _, pattern := range resourcePermission.GetResource().GetPattern() {
			if resourcename.Match(pattern, resourceName) {
				return resourcePermission.GetPermission(), true
			}
		}
	}
	return "", false
}
