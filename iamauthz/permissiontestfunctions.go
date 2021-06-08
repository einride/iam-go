package iamauthz

import (
	"context"
	"reflect"
	"time"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

const (
	testFunction            = "test"
	testFunctionOverload    = "test_caller_string_bool"
	testAllFunction         = "test_all"
	testAllFunctionOverload = "test_all_caller_strings_bool"
	testAnyFunction         = "test_any"
	testAnyFunctionOverload = "test_any_caller_strings_bool"
)

type PermissionTester interface {
	TestResourcePermission(
		ctx context.Context, members []string, resource string, permission string,
	) (bool, error)
	TestResourcePermissions(
		ctx context.Context, members []string, resourcePermissions map[string]string,
	) (map[string]bool, error)
}

type PermissionTestFunctions struct {
	options *iamv1.MethodAuthorizationOptions
	tester  PermissionTester
}

func NewPermissionTestFunctions(
	options *iamv1.MethodAuthorizationOptions,
	tester PermissionTester,
) *PermissionTestFunctions {
	return &PermissionTestFunctions{
		options: options,
		tester:  tester,
	}
}

func (f *PermissionTestFunctions) Declarations() []*expr.Decl {
	caller := (&iamv1.Caller{}).ProtoReflect().Descriptor()
	return []*expr.Decl{
		decls.NewFunction(
			testFunction,
			decls.NewOverload(
				testFunctionOverload,
				[]*expr.Type{
					decls.NewObjectType(string(caller.FullName())),
					decls.String,
				},
				decls.Bool,
			),
		),
		decls.NewFunction(
			testAllFunction,
			decls.NewOverload(
				testAllFunctionOverload,
				[]*expr.Type{
					decls.NewObjectType(string(caller.FullName())),
					decls.NewListType(decls.String),
				},
				decls.Bool,
			),
		),
		decls.NewFunction(
			testAnyFunction,
			decls.NewOverload(
				testAnyFunctionOverload,
				[]*expr.Type{
					decls.NewObjectType(string(caller.FullName())),
					decls.NewListType(decls.String),
				},
				decls.Bool,
			),
		),
	}
}

func (f *PermissionTestFunctions) Functions() []*functions.Overload {
	return []*functions.Overload{
		{Operator: testFunctionOverload, Binary: f.testCallerResource},
		{Operator: testAllFunctionOverload, Binary: f.testAllCallerResources},
		{Operator: testAnyFunctionOverload, Binary: f.testAnyCallerResources},
	}
}

func (f *PermissionTestFunctions) testCallerResource(callerVal ref.Val, resourceVal ref.Val) ref.Val {
	caller, ok := callerVal.Value().(*iamv1.Caller)
	if !ok {
		return types.NewErr("test: unexpected type of arg 1, expected %T but got %T", &iamv1.Caller{}, caller)
	}
	resource, ok := resourceVal.Value().(string)
	if !ok {
		return types.NewErr("test: unexpected type of arg 2, expected string but got %T", resource)
	}
	permission, err := ResolvePermissionForResource(f.options, resource)
	if err != nil {
		return types.NewErr(err.Error())
	}
	// TODO: When cel-go supports async functions, use the caller context here.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if result, err := f.tester.TestResourcePermission(
		ctx,
		caller.Members,
		resource,
		permission,
	); err != nil {
		return types.NewErr("test: error testing permission: %v", err)
	} else if !result {
		return types.False
	} else {
		return types.True
	}
}

func (f *PermissionTestFunctions) testAllCallerResources(callerVal, resourcesVal ref.Val) ref.Val {
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
		permission, err := ResolvePermissionForResource(f.options, resource)
		if err != nil {
			return types.NewErr(err.Error())
		}
		resourcePermissions[resource] = permission
	}
	// TODO: When cel-go supports async functions, use the caller context here.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if result, err := f.tester.TestResourcePermissions(
		ctx,
		caller.Members,
		resourcePermissions,
	); err != nil {
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
}

func (f *PermissionTestFunctions) testAnyCallerResources(callerVal, resourcesVal ref.Val) ref.Val {
	caller, ok := callerVal.Value().(*iamv1.Caller)
	if !ok {
		return types.NewErr("test_any: unexpected type of arg 1, expected %T but got %T", &iamv1.Caller{}, caller)
	}
	convertedResources, err := resourcesVal.ConvertToNative(reflect.TypeOf([]string(nil)))
	if err != nil {
		return types.NewErr("test_any: unexpected type of arg 2, expected []string but got %T", resourcesVal)
	}
	resources, ok := convertedResources.([]string)
	if !ok {
		return types.NewErr("test_any: unexpected type of arg 2, expected []string but got %T", resourcesVal)
	}
	if len(resources) == 0 {
		return types.False
	}
	resourcePermissions := make(map[string]string, len(resources))
	for _, resource := range resources {
		permission, err := ResolvePermissionForResource(f.options, resource)
		if err != nil {
			return types.NewErr(err.Error())
		}
		resourcePermissions[resource] = permission
	}
	// TODO: When cel-go supports async functions, use the caller context here.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if result, err := f.tester.TestResourcePermissions(
		ctx,
		caller.Members,
		resourcePermissions,
	); err != nil {
		return types.NewErr("test: error testing permission: %v", err)
	} else {
		if len(result) != len(resources) {
			return types.False
		}
		for _, ok := range result {
			if ok {
				return types.True
			}
		}
		return types.False
	}
}
