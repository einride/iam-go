package iamcel

import (
	"testing"

	"github.com/google/cel-go/cel"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"gotest.tools/v3/assert"
)

func TestMemberFunction(t *testing.T) {
	caller := (&iamv1.Caller{}).ProtoReflect().Descriptor()
	dependencies, err := collectDependencies(caller)
	assert.NilError(t, err)
	env, err := cel.NewEnv(
		cel.TypeDescs(dependencies),
		cel.Variable("caller", cel.ObjectType(string(caller.FullName()))),
		cel.Declarations(NewMemberFunctionDeclaration()),
	)
	assert.NilError(t, err)

	t.Run("single kind", func(t *testing.T) {
		ast, issues := env.Compile(`caller.member('kind1')`)
		assert.NilError(t, issues.Err())
		//nolint: staticcheck // TODO: migrate to new top-level API
		program, err := env.Program(ast, cel.Functions(NewMemberFunctionImplementation()))
		assert.NilError(t, err)
		result, _, err := program.Eval(
			map[string]interface{}{
				"caller": &iamv1.Caller{
					Members: []string{
						"kind1:value1",
						"kind2:value2",
						"kind2:value3",
					},
				},
			},
		)
		assert.NilError(t, err)
		assert.Equal(t, "value1", result.Value().(string))
	})
	t.Run("pick first", func(t *testing.T) {
		ast, issues := env.Compile(`caller.member('kind2')`)
		assert.NilError(t, issues.Err())
		//nolint: staticcheck // TODO: migrate to new top-level API
		program, err := env.Program(ast, cel.Functions(NewMemberFunctionImplementation()))
		assert.NilError(t, err)
		result, _, err := program.Eval(
			map[string]interface{}{
				"caller": &iamv1.Caller{
					Members: []string{
						"kind1:value1",
						"kind2:value2",
						"kind2:value3",
					},
				},
			},
		)
		assert.NilError(t, err)
		assert.Equal(t, "value2", result.Value().(string))
	})
	t.Run("no such kind", func(t *testing.T) {
		ast, issues := env.Compile(`caller.member('kind3')`)
		assert.NilError(t, issues.Err())
		//nolint: staticcheck // TODO: migrate to new top-level API
		program, err := env.Program(ast, cel.Functions(NewMemberFunctionImplementation()))
		assert.NilError(t, err)
		_, _, err = program.Eval(
			map[string]interface{}{
				"caller": &iamv1.Caller{
					Members: []string{
						"kind1:value1",
						"kind2:value2",
						"kind2:value3",
					},
				},
			},
		)
		assert.Error(t, err, "member: no kind 'kind3' found in caller")
	})
}
