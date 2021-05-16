package iamauthz

import (
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

func ResolvePermissionForResource(options *iamv1.MethodAuthorizationOptions, resource string) (string, error) {
	switch permissions := options.Permissions.(type) {
	case *iamv1.MethodAuthorizationOptions_Permission:
		return permissions.Permission, nil
	case *iamv1.MethodAuthorizationOptions_ResourcePermissions:
		return "", nil
	default:
		return "", nil
	}
}
