package iamauthz

import (
	"context"

	"go.einride.tech/iam/iamcaller"
	"go.einride.tech/iam/iamcel"
	"go.einride.tech/iam/iampermission"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BeforeLongRunningOperationMethodAuthorization struct {
	options          *iamv1.LongRunningOperationsAuthorizationOptions
	permissionTester iamcel.PermissionTester
	callerResolver   iamcaller.Resolver
}

func NewBeforeLongRunningOperationMethodAuthorization(
	options *iamv1.LongRunningOperationsAuthorizationOptions,
	permissionTester iamcel.PermissionTester,
	callerResolver iamcaller.Resolver,
) (*BeforeLongRunningOperationMethodAuthorization, error) {
	return &BeforeLongRunningOperationMethodAuthorization{
		options:          options,
		permissionTester: permissionTester,
		callerResolver:   callerResolver,
	}, nil
}

func (a *BeforeLongRunningOperationMethodAuthorization) AuthorizeRequest(
	ctx context.Context,
	request iampermission.LongRunningOperationRequest,
) (context.Context, error) {
	Authorize(ctx)
	permission, ok := iampermission.ResolveLongRunningOperationPermission(a.options.GetOperationPermissions(), request)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "no permission configured for long-running operation request")
	}
	caller, err := a.callerResolver.ResolveCaller(ctx)
	if err != nil {
		return nil, err
	}
	result, err := a.permissionTester.TestPermissions(
		ctx, caller, map[string]string{request.GetName(): permission},
	)
	if err != nil {
		return nil, err
	}
	if !result[request.GetName()] {
		return nil, status.Errorf(codes.PermissionDenied, "operation requires permission %s", permission)
	}
	return ctx, nil
}
