package iamcel

import (
	"context"

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

// TestFunction is the name of the test permission function.
const TestFunction = "test"

// permission test function overloads.
const testFunctionOverload = "test_caller_string_bool"

// NewTestFunctionDeclaration creates a new declaration for the test permission function.
func NewTestFunctionDeclaration() *expr.Decl {
	return decls.NewFunction(
		TestFunction,
		decls.NewOverload(
			testFunctionOverload,
			[]*expr.Type{
				decls.NewObjectType(string((&iamv1.Caller{}).ProtoReflect().Descriptor().FullName())),
				decls.String,
			},
			decls.Bool,
		),
	)
}

// NewTestFunctionImplementation creates a new implementation for the test permission function.
func NewTestFunctionImplementation(
	options *iamv1.MethodAuthorizationOptions,
	tester PermissionTester,
) *functions.Overload {
	return &functions.Overload{
		Operator: testFunctionOverload,
		Binary: func(callerVal ref.Val, resourceVal ref.Val) ref.Val {
			caller, ok := callerVal.Value().(*iamv1.Caller)
			if !ok {
				return types.NewErr("test: unexpected type of arg 1, expected %T but got %T", &iamv1.Caller{}, caller)
			}
			resource, ok := resourceVal.Value().(string)
			if !ok {
				return types.NewErr("test: unexpected type of arg 2, expected string but got %T", resource)
			}
			permission, ok := iampermission.ResolveMethodPermission(options, resource)
			if !ok {
				return types.NewErr("%s: no permission configured for resource '%s'", codes.PermissionDenied, resource)
			}
			// TODO: When cel-go supports async functions, use the caller context here.
			ctx := context.Background()
			if caller.GetContext().GetDeadline() != nil {
				var cancel context.CancelFunc
				ctx, cancel = context.WithDeadline(ctx, caller.GetContext().GetDeadline().AsTime())
				defer cancel()
			}
			if result, err := tester.TestPermissions(ctx, caller, map[string]string{resource: permission}); err != nil {
				if s, ok := status.FromError(err); ok {
					return types.NewErr("%s: %s", s.Code(), s.Message())
				}
				return types.NewErr("test: error testing permission '%s': %v", permission, err)
			} else if !result[resource] {
				return types.False
			} else {
				return types.True
			}
		},
	}
}
