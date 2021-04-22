package authorization

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const (
	testFunction            = "test"
	testFunctionOverload    = "test_caller_string_bool"
	testAllFunction         = "test_all"
	testAllFunctionOverload = "test_all_caller_strings_bool"
	testAnyFunction         = "test_any"
	testAnyFunctionOverload = "test_any_caller_strings_bool"
)

func NewEnv(
	input protoreflect.MessageDescriptor,
	output protoreflect.MessageDescriptor,
	authorization *iamv1.Authorization,
) (*cel.Env, error) {
	caller := &iamv1.Caller{}
	variables := []*expr.Decl{
		decls.NewVar("caller", decls.NewObjectType(string(caller.ProtoReflect().Descriptor().FullName()))),
		decls.NewVar("request", decls.NewObjectType(string(input.FullName()))),
	}
	switch authorization.Strategy.(type) {
	case *iamv1.Authorization_After:
		variables = append(variables, decls.NewVar("response", decls.NewObjectType(string(output.FullName()))))
	}
	reg := collectFileDescriptorSet(input.ParentFile(), output.ParentFile())
	return cel.NewEnv(
		cel.TypeDescs(reg),
		cel.Declarations(variables...),
		cel.Declarations(
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
			decls.NewFunction(
				testAllFunction,
				decls.NewOverload(
					testAllFunctionOverload,
					[]*exprpb.Type{
						decls.NewObjectType(string((&iamv1.Caller{}).ProtoReflect().Descriptor().FullName())),
						decls.NewListType(decls.String), // [resource]
					},
					decls.Bool,
				),
			),
			decls.NewFunction(
				testAnyFunction,
				decls.NewOverload(
					testAnyFunctionOverload,
					[]*exprpb.Type{
						decls.NewObjectType(string((&iamv1.Caller{}).ProtoReflect().Descriptor().FullName())),
						decls.NewListType(decls.String), // [resource]
					},
					decls.Bool,
				),
			),
		),
	)
}

func collectFileDescriptorSet(files ...protoreflect.FileDescriptor) *protoregistry.Files {
	fdMap := map[string]protoreflect.FileDescriptor{}
	for _, parentFile := range files {
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

func NewProgram(
	input protoreflect.MessageDescriptor,
	output protoreflect.MessageDescriptor,
	authorization *iamv1.Authorization,
	iamPolicyServer iam.IAMPolicyServer,
) (cel.Program, error) {
	//permissionTester := NewPermissionTester(authorization.GetPermission(), iamPolicyServer)
	//env, err := NewEnv(input, output, authorization)
	//if err != nil {
	//	return nil, err
	//}
	//ast, issues := env.Compile(policy.GetExpression())
	//if err := issues.Err(); err != nil {
	//	return nil, err
	//}
	//if !proto.Equal(ast.ResultType(), decls.Bool) {
	//	return nil, fmt.Errorf("non-bool result type: %v", ast.ResultType())
	//}
	//program, err := env.Program(ast, cel.Functions(permissionTester.Overloads()...))
	//if err != nil {
	//	return nil, err
	//}
	return nil, nil
}
