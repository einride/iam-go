package iamcel

import (
	"context"
	"fmt"
	"strings"

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

// TestNamespaceFunction is the name of the test_namespace function.
const TestNamespaceFunction = "test_namespace"

// permission test function overloads.
const testNamespaceFunctionOverload = "test_namespace_caller_string_string_bool"

// NewTestNamespaceFunctionDeclaration creates a new declaration for the test_namespace function.
func NewTestNamespaceFunctionDeclaration() *expr.Decl {
	return decls.NewFunction(
		TestNamespaceFunction,
		decls.NewOverload(
			testNamespaceFunctionOverload,
			[]*expr.Type{
				decls.NewObjectType(string((&iamv1.Caller{}).ProtoReflect().Descriptor().FullName())),
				decls.String,
				decls.String,
			},
			decls.Bool,
		),
	)
}

// NewTestNamespaceFunctionImplementation creates a new implementation for the test_namespace function.
func NewTestNamespaceFunctionImplementation(
	options *iamv1.MethodAuthorizationOptions,
	tester PermissionTester,
) *functions.Overload {
	return &functions.Overload{
		Operator: testNamespaceFunctionOverload,
		Function: func(values ...ref.Val) ref.Val {
			if len(values) != 3 {
				return types.NewErr("test: unexpected number of arguments, expected 3 but got %d", len(values))
			}
			caller, ok := values[0].Value().(*iamv1.Caller)
			if !ok {
				return types.NewErr("test: unexpected type of arg 1, expected %T but got %T", &iamv1.Caller{}, values[0].Value())
			}
			namespace, ok := values[1].Value().(string)
			if !ok {
				return types.NewErr("test: unexpected type of arg 2, expected string but got %T", values[1].Value())
			}
			resource, ok := values[2].Value().(string)
			if !ok {
				return types.NewErr("test: unexpected type of arg 3, expected string but got %T", values[2].Value())
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

			namespacedResource := namespaceResource(namespace, resource)
			result, err := tester.TestPermissions(ctx, caller, map[string]string{namespacedResource: permission})
			switch {
			case err != nil:
				if s, ok := status.FromError(err); ok {
					return types.NewErr("%s: %s", s.Code(), s.Message())
				}
				return types.NewErr("test: error testing permission '%s': %v", permission, err)
			case !result[namespacedResource]:
				return types.False
			default:
				return types.True
			}
		},
	}
}

func namespaceResource(namespace, resource string) string {
	for len(namespace) > 1 && strings.HasSuffix(namespace, "/") {
		namespace = strings.TrimSuffix(namespace, "/")
	}
	return fmt.Sprintf("%s/%s", namespace, resource)
}
