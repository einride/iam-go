package iamspanner

import (
	"context"
	"fmt"

	"cloud.google.com/go/iam/admin/apiv1/adminpb"
	"cloud.google.com/go/iam/apiv1/iampb"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iampermission"
	"go.einride.tech/iam/iamresource"
)

// TestIamPermissions implements iampb.IAMPolicyServer.
func (s *IAMServer) TestIamPermissions(
	ctx context.Context,
	request *iampb.TestIamPermissionsRequest,
) (*iampb.TestIamPermissionsResponse, error) {
	if err := validateTestIamPermissionsRequest(request); err != nil {
		return nil, err
	}
	caller, err := s.callerResolver.ResolveCaller(ctx)
	if err != nil {
		return nil, err
	}
	permissions := make(map[string]struct{}, len(request.GetPermissions()))
	tx := s.client.Single()
	defer tx.Close()
	if err := s.ReadBindingsByResourcesAndMembersInTransaction(
		ctx,
		tx,
		[]string{request.GetResource()},
		caller.GetMembers(),
		func(_ context.Context, _ string, role *adminpb.Role, _ string) error {
			for _, permission := range request.GetPermissions() {
				if s.roles.RoleHasPermission(role.GetName(), permission) {
					permissions[permission] = struct{}{}
				}
			}
			return nil
		},
	); err != nil {
		return nil, s.handleStorageError(ctx, err)
	}
	response := &iampb.TestIamPermissionsResponse{
		Permissions: make([]string, 0, len(permissions)),
	}
	for _, permission := range request.GetPermissions() {
		if _, ok := permissions[permission]; ok {
			response.Permissions = append(response.Permissions, permission)
		}
	}
	return response, nil
}

func validateTestIamPermissionsRequest(request *iampb.TestIamPermissionsRequest) error {
	var result validation.MessageValidator
	switch request.GetResource() {
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
	for i, permission := range request.GetPermissions() {
		if err := iampermission.Validate(permission); err != nil {
			result.AddFieldError(fmt.Sprintf("permissions[%d]", i), err)
		}
	}
	return result.Err()
}
