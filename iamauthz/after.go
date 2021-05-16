package iamauthz

import (
	"context"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"go.einride.tech/iam/iammember"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type AfterMethodAuthorization struct {
	methodAuthorizationOptions *iamv1.MethodAuthorizationOptions
	memberResolver             iammember.Resolver
	program                    cel.Program
}

func NewAfterMethodAuthorization(
	method protoreflect.MethodDescriptor,
	permissionTester PermissionTester,
	memberResolver iammember.Resolver,
) (*AfterMethodAuthorization, error) {
	methodAuthorizationOptions := proto.GetExtension(
		method.Options(), iamv1.E_MethodAuthorization,
	).(*iamv1.MethodAuthorizationOptions)
	if methodAuthorizationOptions == nil {
		return nil, fmt.Errorf("missing method_authorization annotation")
	}
	afterStrategy, ok := methodAuthorizationOptions.Strategy.(*iamv1.MethodAuthorizationOptions_After)
	if !ok {
		return nil, fmt.Errorf("strategy must be 'after'")
	}
	caller := (&iamv1.Caller{}).ProtoReflect().Descriptor()
	fns := NewFunctions(methodAuthorizationOptions, permissionTester)
	env, err := cel.NewEnv(
		cel.TypeDescs(collectTypeDescs(caller, method.Input(), method.Output())),
		cel.Declarations(
			decls.NewVar("caller", decls.NewObjectType(string(caller.FullName()))),
			decls.NewVar("request", decls.NewObjectType(string(method.Input().FullName()))),
			decls.NewVar("response", decls.NewObjectType(string(method.Output().FullName()))),
		),
		cel.Declarations(fns.Declarations()...),
	)
	if err != nil {
		return nil, err
	}
	ast, issues := env.Compile(afterStrategy.After.Expression)
	if issues.Err() != nil {
		return nil, issues.Err()
	}
	program, err := env.Program(ast, cel.Functions(fns.Functions()...))
	if err != nil {
		return nil, err
	}
	return &AfterMethodAuthorization{
		methodAuthorizationOptions: methodAuthorizationOptions,
		memberResolver:             memberResolver,
		program:                    program,
	}, nil
}

func (a *AfterMethodAuthorization) AuthorizeRequestAndResponse(
	ctx context.Context,
	request proto.Message,
	response proto.Message,
) (context.Context, error) {
	Authorize(ctx)
	ctx, members, err := a.memberResolver.ResolveIAMMembers(ctx)
	if err != nil {
		return nil, err
	}
	val, _, err := a.program.Eval(map[string]interface{}{
		"caller":   &iamv1.Caller{Members: members},
		"request":  request,
		"response": response,
	})
	if err != nil {
		return nil, err
	}
	boolVal, ok := val.Value().(bool)
	if !ok {
		return nil, status.Error(codes.Internal, "authorization policy returned non-bool result")
	}
	if !boolVal {
		return nil, status.Error(codes.PermissionDenied, a.methodAuthorizationOptions.GetAfter().GetDescription())
	}
	return ctx, nil
}
