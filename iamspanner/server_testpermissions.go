package iamspanner

import (
	"context"

	"cloud.google.com/go/iam/admin/apiv1/adminpb"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iamresource"
	"go.einride.tech/iam/iamrole"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// TestPermissions implements iamcel.PermissionTester.
func (s *IAMServer) TestPermissions(
	ctx context.Context,
	caller *iamv1.Caller,
	resourcePermissions map[string]string,
) (map[string]bool, error) {
	result := make(map[string]bool, len(resourcePermissions))
	tx := s.client.Single()
	defer tx.Close()
	resources := make([]string, 0, len(resourcePermissions))
	for resource := range resourcePermissions {
		resources = append(resources, resource)
	}
	if err := s.ReadBindingsByResourcesAndMembersInTransaction(
		ctx,
		tx,
		resources,
		caller.GetMembers(),
		func(_ context.Context, boundResource string, role *adminpb.Role, _ string) error {
			for resource, permission := range resourcePermissions {
				result[resource] = result[resource] ||
					(boundResource == iamresource.Root ||
						resource == boundResource ||
						resourcename.HasParent(resource, boundResource) &&
							iamrole.HasPermission(role, permission))
			}
			return nil
		},
	); err != nil {
		return nil, s.handleStorageError(ctx, err)
	}
	for resource := range resourcePermissions {
		if _, ok := result[resource]; !ok {
			result[resource] = false
		}
	}
	return result, nil
}
