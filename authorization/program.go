package authorization

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	authorizationv1 "go.einride.tech/authorization-aip/proto/gen/einride/authorization/v1"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func NewPolicyEnv(
	input protoreflect.MessageDescriptor,
	output protoreflect.MessageDescriptor,
	policy *authorizationv1.Policy,
) (*cel.Env, error) {
	caller := &authorizationv1.Caller{}
	variables := []*expr.Decl{
		decls.NewVar("caller", decls.NewObjectType(string(caller.ProtoReflect().Descriptor().FullName()))),
		decls.NewVar("request", decls.NewObjectType(string(input.FullName()))),
	}
	if policy.GetDecisionPoint() == authorizationv1.PolicyDecisionPoint_AFTER {
		variables = append(variables, decls.NewVar("response", decls.NewObjectType(string(output.FullName()))))
	}
	return cel.NewEnv(
		cel.TypeDescs(input.ParentFile(), output.ParentFile()),
		cel.Declarations(variables...),
		cel.Declarations(
			decls.NewFunction(
				testFunction,
				decls.NewOverload(
					testFunctionOverload,
					[]*exprpb.Type{
						decls.NewObjectType(string((&authorizationv1.Caller{}).ProtoReflect().Descriptor().FullName())),
						decls.String, // resource
					},
					decls.Bool,
				),
			),
		),
	)
}

func NewPolicyProgram(
	input protoreflect.MessageDescriptor,
	output protoreflect.MessageDescriptor,
	policy *authorizationv1.Policy,
	iamPolicyServer iam.IAMPolicyServer,
) (cel.Program, error) {
	permissionTester := NewPermissionTester(policy.GetPermission(), iamPolicyServer)
	env, err := NewPolicyEnv(input, output, policy)
	if err != nil {
		return nil, err
	}
	ast, issues := env.Compile(policy.GetExpression())
	if err := issues.Err(); err != nil {
		return nil, err
	}
	if !proto.Equal(ast.ResultType(), decls.Bool) {
		return nil, fmt.Errorf("non-bool result type: %v", ast.ResultType())
	}
	program, err := env.Program(ast, cel.Functions(permissionTester.Overloads()...))
	if err != nil {
		return nil, err
	}
	return program, nil
}
