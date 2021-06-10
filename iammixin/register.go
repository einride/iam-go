package iammixin

import (
	"context"

	"google.golang.org/genproto/googleapis/iam/admin/v1"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/grpc"
)

// Server is an interface for servers that implement the essential IAM mixins.
type Server interface {
	iam.IAMPolicyServer
	ListRoles(context.Context, *admin.ListRolesRequest) (*admin.ListRolesResponse, error)
	GetRole(context.Context, *admin.GetRoleRequest) (*admin.Role, error)
}

// Register the IAM mixin server with the provided gRPC server.
func Register(server *grpc.Server, serverImpl Server) {
	iam.RegisterIAMPolicyServer(server, serverImpl)
	admin.RegisterIAMServer(server, &adminAdapter{server: serverImpl})
}

// adminAdapter provides unimplemented methods for the non-essential IAM admin mixins.
type adminAdapter struct {
	admin.UnimplementedIAMServer
	server Server
}

// ListRoles implements admin.IAMServer.
func (a *adminAdapter) ListRoles(
	ctx context.Context,
	request *admin.ListRolesRequest,
) (*admin.ListRolesResponse, error) {
	return a.server.ListRoles(ctx, request)
}

// GetRole implements admin.IAMServer.
func (a *adminAdapter) GetRole(
	ctx context.Context,
	request *admin.GetRoleRequest,
) (*admin.Role, error) {
	return a.server.GetRole(ctx, request)
}

// SetIamPolicy implements admin.IAMServer.
func (a *adminAdapter) SetIamPolicy(
	ctx context.Context,
	request *iam.SetIamPolicyRequest,
) (*iam.Policy, error) {
	return a.server.SetIamPolicy(ctx, request)
}

// GetIamPolicy implements admin.IAMServer.
func (a *adminAdapter) GetIamPolicy(
	ctx context.Context,
	request *iam.GetIamPolicyRequest,
) (*iam.Policy, error) {
	return a.server.GetIamPolicy(ctx, request)
}

// TestIamPermissions implements admin.IAMServer.
func (a *adminAdapter) TestIamPermissions(
	ctx context.Context,
	request *iam.TestIamPermissionsRequest,
) (*iam.TestIamPermissionsResponse, error) {
	return a.server.TestIamPermissions(ctx, request)
}
