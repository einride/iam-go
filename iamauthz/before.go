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
	"google.golang.org/protobuf/reflect/protoregistry"
)

type BeforeMethodAuthorization struct {
	methodAuthorizationOptions *iamv1.MethodAuthorizationOptions
	memberResolver             iammember.Resolver
	program                    cel.Program
}

func NewBeforeMethodAuthorization(
	method protoreflect.MethodDescriptor,
	permissionTester PermissionTester,
	memberResolver iammember.Resolver,
) (*BeforeMethodAuthorization, error) {
	methodAuthorizationOptions := proto.GetExtension(
		method.Options(), iamv1.E_MethodAuthorization,
	).(*iamv1.MethodAuthorizationOptions)
	if methodAuthorizationOptions == nil {
		return nil, fmt.Errorf("missing method_authorization annotation")
	}
	beforeStrategy, ok := methodAuthorizationOptions.Strategy.(*iamv1.MethodAuthorizationOptions_Before)
	if !ok {
		return nil, fmt.Errorf("strategy must be 'before'")
	}
	fns := NewFunctions(methodAuthorizationOptions, permissionTester)
	caller := (&iamv1.Caller{}).ProtoReflect().Descriptor()
	env, err := cel.NewEnv(
		cel.TypeDescs(collectTypeDescs(caller, method.Input())),
		cel.Declarations(
			decls.NewVar("caller", decls.NewObjectType(string(caller.FullName()))),
			decls.NewVar("request", decls.NewObjectType(string(method.Input().FullName()))),
		),
		cel.Declarations(fns.Declarations()...),
	)
	if err != nil {
		return nil, err
	}
	ast, issues := env.Compile(beforeStrategy.Before.Expression)
	if issues.Err() != nil {
		return nil, issues.Err()
	}
	program, err := env.Program(ast, cel.Functions(fns.Functions()...))
	if err != nil {
		return nil, err
	}
	return &BeforeMethodAuthorization{
		methodAuthorizationOptions: methodAuthorizationOptions,
		memberResolver:             memberResolver,
		program:                    program,
	}, nil
}

func (a *BeforeMethodAuthorization) AuthorizeRequest(
	ctx context.Context,
	request proto.Message,
) (context.Context, error) {
	Authorize(ctx)
	ctx, members, err := a.memberResolver.ResolveIAMMembers(ctx)
	if err != nil {
		return nil, err
	}
	val, _, err := a.program.Eval(map[string]interface{}{
		"caller":  &iamv1.Caller{Members: members},
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
		return nil, status.Error(codes.PermissionDenied, a.methodAuthorizationOptions.GetBefore().GetDescription())
	}
	return ctx, nil
}

func collectTypeDescs(messages ...protoreflect.MessageDescriptor) *protoregistry.Files {
	fdMap := map[string]protoreflect.FileDescriptor{}
	for _, message := range messages {
		parentFile := message.ParentFile()
		fdMap[parentFile.Path()] = parentFile
		// Initialize list of dependencies
		deps := make([]protoreflect.FileImport, parentFile.Imports().Len())
		for i := 0; i < parentFile.Imports().Len(); i++ {
			deps[i] = parentFile.Imports().Get(i)
		}
		// Expand list for new dependencies
		for i := 0; i < len(deps); i++ {
			dep := deps[i]
			if _, found := fdMap[dep.Path()]; found {
				continue
			}
			fdMap[dep.Path()] = dep.FileDescriptor
			for j := 0; j < dep.FileDescriptor.Imports().Len(); j++ {
				deps = append(deps, dep.FileDescriptor.Imports().Get(j))
			}
		}
	}
	var registry protoregistry.Files
	for _, fd := range fdMap {
		if err := registry.RegisterFile(fd); err != nil {
			panic(err)
		}
	}
	return &registry
}
