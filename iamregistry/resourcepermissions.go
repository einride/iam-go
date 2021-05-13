package iamregistry

import (
	"fmt"

	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iamresource"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/genproto/googleapis/api/annotations"
)

// ResourcePermissions contain a mapping from resource types to permissions.
type ResourcePermissions struct {
	patternPermissions map[string]string
	rootPermission     string
}

// NewResourcePermissions creates a mapping of resource types to permissions.
func NewResourcePermissions(
	resourcePermissions *iamv1.ResourcePermissions,
	resources []*annotations.ResourceDescriptor,
) (*ResourcePermissions, error) {
	result := ResourcePermissions{
		patternPermissions: make(map[string]string, len(resources)),
	}
ResourcePermissionLoop:
	for _, resourcePermission := range resourcePermissions.Resource {
		if resourcePermission.Type == iamresource.Root {
			result.rootPermission = resourcePermission.Permission
			continue
		}
		for _, resource := range resources {
			if resource.Type == resourcePermission.Type {
				for _, pattern := range resource.Pattern {
					if existingPermission, ok := result.patternPermissions[pattern]; ok {
						return nil, fmt.Errorf(
							"register resource %s: pattern %s already mapped to permission %s",
							resourcePermission.Type,
							pattern,
							existingPermission,
						)
					}
					result.patternPermissions[pattern] = resourcePermission.Permission
				}
				continue ResourcePermissionLoop
			}
		}
		return nil, fmt.Errorf("found no resource with type %s", resourcePermission.Type)
	}
	return &result, nil
}

// FindPermissionByResourceName looks up a permission by resource name.
func (r *ResourcePermissions) FindPermissionByResourceName(name string) (string, bool) {
	if r.rootPermission != "" && name == iamresource.Root {
		return r.rootPermission, true
	}
	for pattern, permission := range r.patternPermissions {
		if resourcename.Match(pattern, name) {
			return permission, true
		}
	}
	return "", false
}
