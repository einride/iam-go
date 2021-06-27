package iamcel

import (
	"context"
	"reflect"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"go.einride.tech/iam/iampermission"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestAllFunction is the name of the function for testing that all resources have a specified permission.
const TestAllFunction = "test_all"

const testAllFunctionOverload = "test_all_caller_strings_bool"

// NewTestAllFunctionDeclaration creates a new declaration for the test_all function.
func NewTestAllFunctionDeclaration() *expr.Decl {
	return decls.NewFunction(
		TestAllFunction,
		decls.NewOverload(
			testAllFunctionOverload,
			[]*expr.Type{
				decls.NewObjectType(string((&iamv1.Caller{}).ProtoReflect().Descriptor().FullName())),
				decls.NewListType(decls.String),
			},
			decls.Bool,
		),
	)
}

// NewTestAllFunctionImplementation creates a new implementation for the test_all function.
func NewTestAllFunctionImplementation(
	options *iamv1.MethodAuthorizationOptions,
	tester PermissionTester,
) *functions.Overload {
	return &functions.Overload{
		Operator: testAllFunctionOverload,
		Binary: func(callerVal ref.Val, resourcesVal ref.Val) ref.Val {
			caller, ok := callerVal.Value().(*iamv1.Caller)
			if !ok {
				return types.NewErr("test_all: unexpected type of arg 1, expected %T but got %T", &iamv1.Caller{}, caller)
			}
			convertedResources, err := resourcesVal.ConvertToNative(reflect.TypeOf([]string(nil)))
			if err != nil {
				return types.NewErr("test_all: unexpected type of arg 2, expected []string but got %T", resourcesVal)
			}
			resources, ok := convertedResources.([]string)
			if !ok {
				return types.NewErr("test_all: unexpected type of arg 2, expected []string but got %T", resourcesVal)
			}
			if len(resources) == 0 {
				return types.False
			}
			resourcePermissions := make(map[string]string, len(resources))
			for _, resource := range resources {
				permission, ok := iampermission.ResolveMethodPermission(options, resource)
				if !ok {
					return types.NewErr(
						"%s: no permission configured for resource '%s'", codes.InvalidArgument, resource,
					)
				}
				resourcePermissions[resource] = permission
			}
			// TODO: When cel-go supports async functions, use the caller context here.
			ctx := context.Background()
			if caller.GetContext().GetDeadline() != nil {
				var cancel context.CancelFunc
				ctx, cancel = context.WithDeadline(ctx, caller.GetContext().GetDeadline().AsTime())
				defer cancel()
			}
			if result, err := tester.TestPermissions(ctx, caller, resourcePermissions); err != nil {
				if s, ok := status.FromError(err); ok {
					return types.NewErr("%s: %s", s.Code(), s.Message())
				}
				return types.NewErr("test: error testing permission: %v", err)
			} else {
				if len(result) != len(resources) {
					return types.False
				}
				for _, ok := range result {
					if !ok {
						return types.False
					}
				}
				return types.True
			}
		},
	}
}
