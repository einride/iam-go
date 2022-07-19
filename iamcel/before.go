package iamcel

import (
	"fmt"

	"github.com/google/cel-go/cel"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// NewBeforeEnv creates a new CEL environment for authorization checks that run before the request has been handled.
func NewBeforeEnv(method protoreflect.MethodDescriptor) (*cel.Env, error) {
	caller := (&iamv1.Caller{}).ProtoReflect().Descriptor()
	descriptors, err := collectDependencies(caller, method.Input())
	if err != nil {
		return nil, fmt.Errorf("new IAM CEL `before` env: %w", err)
	}
	env, err := cel.NewEnv(
		cel.TypeDescs(descriptors),
		cel.Variable("caller", cel.ObjectType(string(caller.FullName()))),
		cel.Variable("request", cel.ObjectType(string(method.Input().FullName()))),
		cel.Declarations(
			// TODO: Migrate declarations to new top-level API.
			NewTestFunctionDeclaration(),
			NewTestAllFunctionDeclaration(),
			NewTestAnyFunctionDeclaration(),
			NewAncestorFunctionDeclaration(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("new IAM CEL `before` env: %w", err)
	}
	return env, nil
}
