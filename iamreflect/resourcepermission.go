package iamreflect

import (
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iamresource"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// ResolveMethodPermission resolves a permission for a method, given the requested resource.
func ResolveMethodPermission(
	options *iamv1.MethodAuthorizationOptions,
	resourceName string,
) (string, bool) {
	switch permissions := options.Permissions.(type) {
	case *iamv1.MethodAuthorizationOptions_Permission:
		return permissions.Permission, true
	case *iamv1.MethodAuthorizationOptions_ResourcePermissions:
		return ResolveResourcePermission(permissions.ResourcePermissions.GetResourcePermission(), resourceName)
	default:
		return "", false
	}
}

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
