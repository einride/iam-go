package iamexample

import (
	"context"

	"google.golang.org/genproto/googleapis/iam/admin/v1"
	"google.golang.org/genproto/googleapis/iam/v1"
)

// SetIamPolicy implements iam.IAMPolicyServer.
func (s *Server) SetIamPolicy(
	ctx context.Context,
	request *iam.SetIamPolicyRequest,
) (*iam.Policy, error) {
	return s.IAM.SetIamPolicy(ctx, request)
}

// GetIamPolicy implements iam.IAMPolicyServer.
func (s *Server) GetIamPolicy(
	ctx context.Context,
	request *iam.GetIamPolicyRequest,
) (*iam.Policy, error) {
	return s.IAM.GetIamPolicy(ctx, request)
}

// TestIamPermissions implements iam.IAMPolicyServer.
func (s *Server) TestIamPermissions(
	ctx context.Context,
	request *iam.TestIamPermissionsRequest,
) (*iam.TestIamPermissionsResponse, error) {
	return s.IAM.TestIamPermissions(ctx, request)
}

// ListRoles implements admin.IAMServer.
func (s *Server) ListRoles(
	ctx context.Context,
	request *admin.ListRolesRequest,
) (*admin.ListRolesResponse, error) {
	return s.IAM.ListRoles(ctx, request)
}

// GetRole implements admin.IAMServer.
func (s *Server) GetRole(
	ctx context.Context,
	request *admin.GetRoleRequest,
) (*admin.Role, error) {
	return s.IAM.GetRole(ctx, request)
}
