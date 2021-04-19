package authorization

import (
	"context"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/grpc/metadata"
)

type PermissionTester struct {
	permission      string
	iamPolicyServer iam.IAMPolicyServer
}

func NewPermissionTester(permission string, iamPolicyServer iam.IAMPolicyServer) *PermissionTester {
	return &PermissionTester{
		permission:      permission,
		iamPolicyServer: iamPolicyServer,
	}
}

func (t *PermissionTester) Declarations() []*exprpb.Decl {
	return []*exprpb.Decl{
		decls.NewFunction(
			testFunction,
			decls.NewOverload(
				testFunctionOverload,
				[]*exprpb.Type{
					decls.NewObjectType(string((&iamv1.Caller{}).ProtoReflect().Descriptor().FullName())),
					decls.String, // resource
				},
				decls.Bool,
			),
		),
	}
}

func (t *PermissionTester) Overloads() []*functions.Overload {
	return []*functions.Overload{
		{
			Operator: testFunctionOverload,
			Function: t.test,
		},
	}
}

func (t *PermissionTester) test(args ...ref.Val) ref.Val {
	if len(args) != 2 {
		return types.NewErr("no such overload")
	}
	caller, ok := args[0].Value().(*iamv1.Caller)
	if !ok {
		return types.NewErr("unexpected type of arg 1")
	}
	resource, ok := args[1].Value().(string)
	if !ok {
		return types.NewErr("unexpected type of arg 2")
	}
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("authorization", "Bearer: "+caller.IdToken),
	)
	response, err := t.iamPolicyServer.TestIamPermissions(ctx, &iam.TestIamPermissionsRequest{
		Resource:    resource,
		Permissions: []string{t.permission},
	})
	if err != nil {
		return types.NewErr("test permissions: %v", err)
	}
	if len(response.Permissions) == 1 && response.Permissions[0] == t.permission {
		return types.True
	}
	return types.False
}
