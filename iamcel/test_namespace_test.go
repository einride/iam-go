package iamcel

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/protobuf/testing/protocmp"
	"gotest.tools/v3/assert"
)

func TestTestNamespaceFunction(t *testing.T) {
	caller := (&iamv1.Caller{}).ProtoReflect().Descriptor()
	dependencies, err := collectDependencies(caller)
	assert.NilError(t, err)
	env, err := cel.NewEnv(
		cel.TypeDescs(dependencies),
		cel.Variable("caller", cel.ObjectType(string(caller.FullName()))),
		cel.Declarations(NewTestNamespaceFunctionDeclaration()),
	)
	assert.NilError(t, err)

	const (
		resource       = "resource/1234"
		otherResource  = "resources/2345"
		namespace      = "namespace/5678"
		otherNamespace = "namespace/6789"
		permission     = "resource.permission"
	)

	authzOptions := &iamv1.MethodAuthorizationOptions{
		Permissions: &iamv1.MethodAuthorizationOptions_Permission{
			Permission: permission,
		},
	}

	t.Run("allowed namespace and resource", func(t *testing.T) {
		tester := &mockTester{
			t: t,
			expectResources: map[string]string{
				namespace + "/" + resource: permission,
			},
			returnMap: map[string]bool{
				namespace + "/" + resource: true,
			},
		}
		ast, issues := env.Compile(
			fmt.Sprintf(
				`test_namespace(caller, '%s', '%s')`,
				namespace,
				resource,
			),
		)
		assert.NilError(t, issues.Err())
		//nolint: staticcheck // TODO: migrate to new top-level API
		program, err := env.Program(ast, cel.Functions(NewTestNamespaceFunctionImplementation(authzOptions, tester)))
		assert.NilError(t, err)
		result, _, err := program.Eval(map[string]interface{}{"caller": &iamv1.Caller{}})
		assert.NilError(t, err)
		assert.Equal(t, true, result.Value().(bool))
	})

	t.Run("disallowed namespace", func(t *testing.T) {
		tester := &mockTester{
			t: t,
			expectResources: map[string]string{
				otherNamespace + "/" + resource: permission,
			},
			returnMap: map[string]bool{},
		}
		ast, issues := env.Compile(
			fmt.Sprintf(
				`test_namespace(caller, '%s', '%s')`,
				otherNamespace,
				resource,
			),
		)
		assert.NilError(t, issues.Err())
		//nolint: staticcheck // TODO: migrate to new top-level API
		program, err := env.Program(ast, cel.Functions(NewTestNamespaceFunctionImplementation(authzOptions, tester)))
		assert.NilError(t, err)
		result, _, err := program.Eval(map[string]interface{}{"caller": &iamv1.Caller{}})
		assert.NilError(t, err)
		assert.Equal(t, false, result.Value().(bool))
	})

	t.Run("disallowed resource", func(t *testing.T) {
		tester := &mockTester{
			t: t,
			expectResources: map[string]string{
				namespace + "/" + otherResource: permission,
			},
			returnMap: map[string]bool{},
		}
		ast, issues := env.Compile(
			fmt.Sprintf(
				`test_namespace(caller, '%s', '%s')`,
				namespace,
				otherResource,
			),
		)
		assert.NilError(t, issues.Err())
		//nolint: staticcheck // TODO: migrate to new top-level API
		program, err := env.Program(ast, cel.Functions(NewTestNamespaceFunctionImplementation(authzOptions, tester)))
		assert.NilError(t, err)
		result, _, err := program.Eval(map[string]interface{}{"caller": &iamv1.Caller{}})
		assert.NilError(t, err)
		assert.Equal(t, false, result.Value().(bool))
	})
}

type mockTester struct {
	t               *testing.T
	expectResources map[string]string
	returnMap       map[string]bool
	returnErr       error
}

func (m *mockTester) TestPermissions(
	_ context.Context,
	caller *iamv1.Caller,
	resources map[string]string,
) (map[string]bool, error) {
	assert.DeepEqual(m.t, caller, &iamv1.Caller{}, protocmp.Transform())
	assert.DeepEqual(m.t, resources, m.expectResources)
	return m.returnMap, m.returnErr
}

func TestNamespaceResource(t *testing.T) {
	assert.Equal(t, namespaceResource("namespace", "resource"), "namespace/resource")
	assert.Equal(t, namespaceResource("namespace/1", "resource"), "namespace/1/resource")
	assert.Equal(t, namespaceResource("namespace/1/", "resource"), "namespace/1/resource")
	assert.Equal(t, namespaceResource("namespace/1//", "resource"), "namespace/1/resource")
	assert.Equal(t, namespaceResource("namespace/1////", "resource"), "namespace/1/resource")
}
