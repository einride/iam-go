package iamauthz

import (
	"context"
	"fmt"

	"github.com/google/cel-go/cel"
	"go.einride.tech/iam/iamcel"
	"go.einride.tech/iam/iammember"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type BeforeMethodAuthorization struct {
	options        *iamv1.MethodAuthorizationOptions
	memberResolver iammember.Resolver
	program        cel.Program
}

func NewBeforeMethodAuthorization(
	method protoreflect.MethodDescriptor,
	options *iamv1.MethodAuthorizationOptions,
	permissionTester iamcel.PermissionTester,
	memberResolver iammember.Resolver,
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
		cel.Functions(
			iamcel.NewTestFunctionImplementation(options, permissionTester),
			iamcel.NewTestAllFunctionImplementation(options, permissionTester),
			iamcel.NewTestAnyFunctionImplementation(options, permissionTester),
			iamcel.NewAncestorFunctionImplementation(),
		),
	)
	if err != nil {
		return nil, err
	}
	return &BeforeMethodAuthorization{
		options:        options,
		memberResolver: memberResolver,
		program:        program,
	}, nil
}

func (a *BeforeMethodAuthorization) AuthorizeRequest(
	ctx context.Context,
	request proto.Message,
) (context.Context, error) {
	Authorize(ctx)
	memberResolveResult, err := a.memberResolver.ResolveIAMMembers(ctx)
	if err != nil {
		return nil, err
	}
	val, _, err := a.program.Eval(map[string]interface{}{
		"caller":  &iamv1.Caller{Members: memberResolveResult.Members()},
		"request": request,
	})
	if err != nil {
		return nil, err
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
