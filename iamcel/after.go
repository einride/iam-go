package iamcel

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
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
		cel.Declarations(
			decls.NewVar("caller", decls.NewObjectType(string(caller.FullName()))),
			decls.NewVar("request", decls.NewObjectType(string(method.Input().FullName()))),
			decls.NewVar("response", decls.NewObjectType(string(method.Output().FullName()))),
			NewTestFunctionDeclaration(),
			NewTestAllFunctionDeclaration(),
			NewTestAnyFunctionDeclaration(),
			NewAncestorFunctionDeclaration(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("new IAM CEL `after` env: %w", err)
	}
	return env, nil
}
