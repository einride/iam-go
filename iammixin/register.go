package iammixin

import (
	"context"

	"cloud.google.com/go/iam/admin/apiv1/adminpb"
	"cloud.google.com/go/iam/apiv1/iampb"
	"google.golang.org/grpc"
)

// Server is an interface for servers that implement the essential IAM mixins.
type Server interface {
	iampb.IAMPolicyServer
	ListRoles(context.Context, *adminpb.ListRolesRequest) (*adminpb.ListRolesResponse, error)
	GetRole(context.Context, *adminpb.GetRoleRequest) (*adminpb.Role, error)
}

// Register the IAM mixin server with the provided gRPC server.
func Register(server *grpc.Server, serverImpl Server) {
	iampb.RegisterIAMPolicyServer(server, serverImpl)
	adminpb.RegisterIAMServer(server, &adminAdapter{server: serverImpl})
}

// adminAdapter provides unimplemented methods for the non-essential IAM adminpb mixins.
type adminAdapter struct {
	adminpb.UnimplementedIAMServer
	server Server
}

// ListRoles implements adminpb.IAMServer.
func (a *adminAdapter) ListRoles(
	ctx context.Context,
	request *adminpb.ListRolesRequest,
) (*adminpb.ListRolesResponse, error) {
	return a.server.ListRoles(ctx, request)
}

// GetRole implements adminpb.IAMServer.
func (a *adminAdapter) GetRole(
	ctx context.Context,
	request *adminpb.GetRoleRequest,
) (*adminpb.Role, error) {
	return a.server.GetRole(ctx, request)
}

// SetIamPolicy implements adminpb.IAMServer.
func (a *adminAdapter) SetIamPolicy(
	ctx context.Context,
	request *iampb.SetIamPolicyRequest,
) (*iampb.Policy, error) {
	return a.server.SetIamPolicy(ctx, request)
}

// GetIamPolicy implements adminpb.IAMServer.
func (a *adminAdapter) GetIamPolicy(
	ctx context.Context,
	request *iampb.GetIamPolicyRequest,
) (*iampb.Policy, error) {
	return a.server.GetIamPolicy(ctx, request)
}

// TestIamPermissions implements adminpb.IAMServer.
func (a *adminAdapter) TestIamPermissions(
	ctx context.Context,
	request *iampb.TestIamPermissionsRequest,
) (*iampb.TestIamPermissionsResponse, error) {
	return a.server.TestIamPermissions(ctx, request)
}
