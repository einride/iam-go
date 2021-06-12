package iamcel

import (
	"testing"

	"github.com/google/cel-go/cel"
	"gotest.tools/v3/assert"
)

func TestResourceNameFunctions(t *testing.T) {
	env, err := cel.NewEnv(cel.Declarations(NewAncestorFunctionDeclaration()))
	assert.NilError(t, err)
	t.Run("ok", func(t *testing.T) {
		ast, issues := env.Compile(`ancestor('foo/1/bar/2', 'foo/{foo}')`)
		assert.NilError(t, issues.Err())
		program, err := env.Program(ast, cel.Functions(NewAncestorFunctionImplementation()))
		assert.NilError(t, err)
		result, _, err := program.Eval(map[string]interface{}(nil))
		assert.NilError(t, err)
		assert.Equal(t, "foo/1", result.Value().(string))
	})
	t.Run("no match", func(t *testing.T) {
		ast, issues := env.Compile(`ancestor('baz/1/bar/2', 'foo/{foo}')`)
		assert.NilError(t, issues.Err())
		program, err := env.Program(ast, cel.Functions(NewAncestorFunctionImplementation()))
		assert.NilError(t, err)
		result, _, err := program.Eval(map[string]interface{}(nil))
		assert.NilError(t, err)
		assert.Equal(t, "", result.Value().(string))
	})
}
