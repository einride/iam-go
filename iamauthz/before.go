package iamauthz

import (
	"context"
	"fmt"

	"github.com/google/cel-go/cel"
	"go.einride.tech/iam/iamcaller"
	"go.einride.tech/iam/iamcel"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type BeforeMethodAuthorization struct {
	options        *iamv1.MethodAuthorizationOptions
	callerResolver iamcaller.Resolver
	program        cel.Program
}

func NewBeforeMethodAuthorization(
	method protoreflect.MethodDescriptor,
	options *iamv1.MethodAuthorizationOptions,
	permissionTester iamcel.PermissionTester,
	callerResolver iamcaller.Resolver,
) (*BeforeMethodAuthorization, error) {
	beforeStrategy, ok := options.Strategy.(*iamv1.MethodAuthorizationOptions_Before)
	if !ok {
		return nil, fmt.Errorf("strategy must be 'before'")
	}
	env, err := iamcel.NewBeforeEnv(method)
	if err != nil {
		return nil, err
	}
	ast, issues := env.Compile(beforeStrategy.Before.Expression)
	if issues.Err() != nil {
		return nil, issues.Err()
	}
	program, err := env.Program(
		ast,
		//nolint: staticcheck // TODO: migrate to new top-level CEL API
		cel.Functions(
			iamcel.NewTestFunctionImplementation(options, permissionTester),
			iamcel.NewTestAllFunctionImplementation(options, permissionTester),
			iamcel.NewTestAnyFunctionImplementation(options, permissionTester),
			iamcel.NewAncestorFunctionImplementation(),
			iamcel.NewMemberFunctionImplementation(),
			iamcel.NewJoinFunctionImplementation(),
		),
	)
	if err != nil {
		return nil, err
	}
	return &BeforeMethodAuthorization{
		options:        options,
		callerResolver: callerResolver,
		program:        program,
	}, nil
}

func (a *BeforeMethodAuthorization) AuthorizeRequest(
	ctx context.Context,
	request proto.Message,
) (context.Context, error) {
	Authorize(ctx)
	caller, err := a.callerResolver.ResolveCaller(ctx)
	if err != nil {
		return nil, err
	}
	val, _, err := a.program.ContextEval(ctx, map[string]interface{}{
		"caller":  caller,
		"request": request,
	})
	if err != nil {
		return nil, forwardErrorCodes(err)
	}
	boolVal, ok := val.Value().(bool)
	if !ok {
		return nil, status.Error(codes.Internal, "authorization policy returned non-bool result")
	}
	if !boolVal {
		return nil, status.Error(codes.PermissionDenied, a.options.GetBefore().GetDescription())
	}
	return ctx, nil
}
