package iamspanner

import (
	"context"
	"fmt"

	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iampermission"
	"go.einride.tech/iam/iamresource"
	"google.golang.org/genproto/googleapis/iam/admin/v1"
	"google.golang.org/genproto/googleapis/iam/v1"
)

// TestIamPermissions implements iam.IAMPolicyServer.
func (s *IAMServer) TestIamPermissions(
	ctx context.Context,
	request *iam.TestIamPermissionsRequest,
) (*iam.TestIamPermissionsResponse, error) {
	if err := validateTestIamPermissionsRequest(request); err != nil {
		return nil, err
	}
	memberResolveResult, err := s.memberResolver.ResolveIAMMembers(ctx)
	if err != nil {
		return nil, err
	}
	permissions := make(map[string]struct{}, len(request.Permissions))
	tx := s.client.Single()
	defer tx.Close()
	if err := s.ReadBindingsByResourcesAndMembersInTransaction(
		ctx,
		tx,
		[]string{request.Resource},
		memberResolveResult.Members,
		func(ctx context.Context, _ string, role *admin.Role, _ string) error {
			for _, permission := range request.Permissions {
				if s.roles.RoleHasPermission(role.Name, permission) {
					permissions[permission] = struct{}{}
				}
			}
			return nil
		},
	); err != nil {
		return nil, s.handleStorageError(ctx, err)
	}
	response := &iam.TestIamPermissionsResponse{
		Permissions: make([]string, 0, len(permissions)),
	}
	for _, permission := range request.Permissions {
		if _, ok := permissions[permission]; ok {
			response.Permissions = append(response.Permissions, permission)
		}
	}
	return response, nil
}

func validateTestIamPermissionsRequest(request *iam.TestIamPermissionsRequest) error {
	var result validation.MessageValidator
	switch request.Resource {
	case iamresource.Root: // OK
	case "":
		result.AddFieldViolation("resource", "missing required field")
	default:
		if err := resourcename.Validate(request.GetResource()); err != nil {
			result.AddFieldError("resource", err)
		}
		if resourcename.ContainsWildcard(request.GetResource()) {
			result.AddFieldViolation("resource", "must not contain wildcard")
		}
	}
	for i, permission := range request.Permissions {
		if err := iampermission.Validate(permission); err != nil {
			result.AddFieldError(fmt.Sprintf("permissions[%d]", i), err)
		}
	}
	return result.Err()
}
