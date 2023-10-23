package iamcel

import (
	"testing"

	"github.com/google/cel-go/cel"
	"gotest.tools/v3/assert"
)

func TestJoinFunction(t *testing.T) {
	env, err := cel.NewEnv(cel.Declarations(NewJoinFunctionDeclaration()))
	assert.NilError(t, err)
	t.Run("ok", func(t *testing.T) {
		ast, issues := env.Compile(`join('parent/1', 'child/2')`)
		assert.NilError(t, issues.Err())
		//nolint: staticcheck // TODO: migrate to new top-level API
		program, err := env.Program(ast, cel.Functions(NewJoinFunctionImplementation()))
		assert.NilError(t, err)
		result, _, err := program.Eval(map[string]interface{}(nil))
		assert.NilError(t, err)
		assert.Equal(t, "parent/1/child/2", result.Value().(string))
	})
	t.Run("root parent", func(t *testing.T) {
		ast, issues := env.Compile(`join('/', 'child/2')`)
		assert.NilError(t, issues.Err())
		//nolint: staticcheck // TODO: migrate to new top-level API
		program, err := env.Program(ast, cel.Functions(NewJoinFunctionImplementation()))
		assert.NilError(t, err)
		result, _, err := program.Eval(map[string]interface{}(nil))
		assert.NilError(t, err)
		assert.Equal(t, "child/2", result.Value().(string))
	})
	t.Run("root child", func(t *testing.T) {
		ast, issues := env.Compile(`join('parent/1', '/')`)
		assert.NilError(t, issues.Err())
		//nolint: staticcheck // TODO: migrate to new top-level API
		program, err := env.Program(ast, cel.Functions(NewJoinFunctionImplementation()))
		assert.NilError(t, err)
		result, _, err := program.Eval(map[string]interface{}(nil))
		assert.NilError(t, err)
		assert.Equal(t, "parent/1", result.Value().(string))
	})
	t.Run("root parent and child", func(t *testing.T) {
		ast, issues := env.Compile(`join('/', '/')`)
		assert.NilError(t, issues.Err())
		//nolint: staticcheck // TODO: migrate to new top-level API
		program, err := env.Program(ast, cel.Functions(NewJoinFunctionImplementation()))
		assert.NilError(t, err)
		result, _, err := program.Eval(map[string]interface{}(nil))
		assert.NilError(t, err)
		assert.Equal(t, "/", result.Value().(string))
	})
	t.Run("empty parent", func(t *testing.T) {
		ast, issues := env.Compile(`join('', 'child/2')`)
		assert.NilError(t, issues.Err())
		//nolint: staticcheck // TODO: migrate to new top-level API
		program, err := env.Program(ast, cel.Functions(NewJoinFunctionImplementation()))
		assert.NilError(t, err)
		result, _, err := program.Eval(map[string]interface{}(nil))
		assert.NilError(t, err)
		assert.Equal(t, "child/2", result.Value().(string))
	})
	t.Run("empty child", func(t *testing.T) {
		ast, issues := env.Compile(`join('parent/1', '')`)
		assert.NilError(t, issues.Err())
		//nolint: staticcheck // TODO: migrate to new top-level API
		program, err := env.Program(ast, cel.Functions(NewJoinFunctionImplementation()))
		assert.NilError(t, err)
		result, _, err := program.Eval(map[string]interface{}(nil))
		assert.NilError(t, err)
		assert.Equal(t, "parent/1", result.Value().(string))
	})
	t.Run("parent slash suffix", func(t *testing.T) {
		ast, issues := env.Compile(`join('parent/1/', 'child/2')`)
		assert.NilError(t, issues.Err())
		//nolint: staticcheck // TODO: migrate to new top-level API
		program, err := env.Program(ast, cel.Functions(NewJoinFunctionImplementation()))
		assert.NilError(t, err)
		result, _, err := program.Eval(map[string]interface{}(nil))
		assert.NilError(t, err)
		assert.Equal(t, "parent/1/child/2", result.Value().(string))
	})
	t.Run("child slash suffix", func(t *testing.T) {
		ast, issues := env.Compile(`join('parent/1', 'child/2/')`)
		assert.NilError(t, issues.Err())
		//nolint: staticcheck // TODO: migrate to new top-level API
		program, err := env.Program(ast, cel.Functions(NewJoinFunctionImplementation()))
		assert.NilError(t, err)
		result, _, err := program.Eval(map[string]interface{}(nil))
		assert.NilError(t, err)
		assert.Equal(t, "parent/1/child/2", result.Value().(string))
	})
	t.Run("parent slash prefix", func(t *testing.T) {
		ast, issues := env.Compile(`join('/parent/1', 'child/2')`)
		assert.NilError(t, issues.Err())
		//nolint: staticcheck // TODO: migrate to new top-level API
		program, err := env.Program(ast, cel.Functions(NewJoinFunctionImplementation()))
		assert.NilError(t, err)
		result, _, err := program.Eval(map[string]interface{}(nil))
		assert.NilError(t, err)
		assert.Equal(t, "parent/1/child/2", result.Value().(string))
	})
	t.Run("child slash prefix", func(t *testing.T) {
		ast, issues := env.Compile(`join('parent/1', '/child/2')`)
		assert.NilError(t, issues.Err())
		//nolint: staticcheck // TODO: migrate to new top-level API
		program, err := env.Program(ast, cel.Functions(NewJoinFunctionImplementation()))
		assert.NilError(t, err)
		result, _, err := program.Eval(map[string]interface{}(nil))
		assert.NilError(t, err)
		assert.Equal(t, "parent/1/child/2", result.Value().(string))
	})
}
