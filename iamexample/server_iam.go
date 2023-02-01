package iamexample

import (
	"context"

	"cloud.google.com/go/iam/admin/apiv1/adminpb"
	"cloud.google.com/go/iam/apiv1/iampb"
)

// SetIamPolicy implements iampb.IAMPolicyServer.
func (s *Server) SetIamPolicy(
	ctx context.Context,
	request *iampb.SetIamPolicyRequest,
) (*iampb.Policy, error) {
	return s.IAM.SetIamPolicy(ctx, request)
}

// GetIamPolicy implements iampb.IAMPolicyServer.
func (s *Server) GetIamPolicy(
	ctx context.Context,
	request *iampb.GetIamPolicyRequest,
) (*iampb.Policy, error) {
	return s.IAM.GetIamPolicy(ctx, request)
}

// TestIamPermissions implements iampb.IAMPolicyServer.
func (s *Server) TestIamPermissions(
	ctx context.Context,
	request *iampb.TestIamPermissionsRequest,
) (*iampb.TestIamPermissionsResponse, error) {
	return s.IAM.TestIamPermissions(ctx, request)
}

// ListRoles implements adminpb.IAMServer.
func (s *Server) ListRoles(
	ctx context.Context,
	request *adminpb.ListRolesRequest,
) (*adminpb.ListRolesResponse, error) {
	return s.IAM.ListRoles(ctx, request)
}

// GetRole implements adminpb.IAMServer.
func (s *Server) GetRole(
	ctx context.Context,
	request *adminpb.GetRoleRequest,
) (*adminpb.Role, error) {
	return s.IAM.GetRole(ctx, request)
}
