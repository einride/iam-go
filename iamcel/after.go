package iamcel

import (
	"fmt"

	"github.com/google/cel-go/cel"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// NewAfterEnv creates a new CEL environment for authorization checks that run after the request has been handled.
func NewAfterEnv(method protoreflect.MethodDescriptor) (*cel.Env, error) {
	caller := (&iamv1.Caller{}).ProtoReflect().Descriptor()
	dependencies, err := collectDependencies(caller, method.Input(), method.Output())
	if err != nil {
		return nil, fmt.Errorf("new IAM CEL `after` env: %w", err)
	}
	env, err := cel.NewEnv(
		cel.TypeDescs(dependencies),
		cel.Variable("caller", cel.ObjectType(string(caller.FullName()))),
		cel.Variable("request", cel.ObjectType(string(method.Input().FullName()))),
		cel.Variable("response", cel.ObjectType(string(method.Output().FullName()))),
		cel.Declarations(
			// TODO: Migrate declarations to new top-level API.
			NewTestFunctionDeclaration(),
			NewTestAllFunctionDeclaration(),
			NewTestAnyFunctionDeclaration(),
			NewAncestorFunctionDeclaration(),
			NewMemberFunctionDeclaration(),
			NewJoinFunctionDeclaration(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("new IAM CEL `after` env: %w", err)
	}
	return env, nil
}
