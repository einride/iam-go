package iamauthz

import (
	"context"

	"go.einride.tech/iam/iammember"
	"go.einride.tech/iam/iampermission"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BeforeLongRunningOperationMethodAuthorization struct {
	operationsPermissions []*iamv1.LongRunningOperationPermissions
	permissionTester      PermissionTester
	memberResolver        iammember.Resolver
}

func NewBeforeLongRunningOperationMethodAuthorization(
	operationsPermissions []*iamv1.LongRunningOperationPermissions,
	permissionTester PermissionTester,
	memberResolver iammember.Resolver,
) (*BeforeLongRunningOperationMethodAuthorization, error) {
	return &BeforeLongRunningOperationMethodAuthorization{
		operationsPermissions: operationsPermissions,
		permissionTester:      permissionTester,
		memberResolver:        memberResolver,
	}, nil
}

func (a *BeforeLongRunningOperationMethodAuthorization) AuthorizeRequest(
	ctx context.Context,
	request iampermission.LongRunningOperationRequest,
) (context.Context, error) {
	Authorize(ctx)
	permission, ok := iampermission.ResolveLongRunningOperationPermission(a.operationsPermissions, request)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "no permission configured for long-running operation request")
	}
	memberResolveResult, err := a.memberResolver.ResolveIAMMembers(ctx)
	if err != nil {
		return nil, err
	}
	result, err := a.permissionTester.TestResourcePermission(
		ctx, memberResolveResult.Members, request.GetName(), permission,
	)
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, status.Errorf(codes.PermissionDenied, "operation requires permission %s", permission)
	}
	return ctx, nil
}
