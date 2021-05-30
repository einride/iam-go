package iammember

import (
	"context"
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func TestChainResolvers(t *testing.T) {
	t.Run("no resolvers", func(t *testing.T) {
		result, err := ChainResolvers().ResolveIAMMembers(context.Background())
		assert.NilError(t, err)
		assert.DeepEqual(t, ResolveResult{}, result)
	})

	t.Run("single", func(t *testing.T) {
		expected := ResolveResult{
			Members: []string{"foo", "bar"},
			Metadata: Metadata{
				"key1": {"foo"},
				"key2": {"bar"},
			},
		}
		actual, err := ChainResolvers(constantResult(expected)).ResolveIAMMembers(context.Background())
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual)
	})

	t.Run("multi", func(t *testing.T) {
		expected := ResolveResult{
			Members: []string{"foo", "bar", "baz"},
			Metadata: Metadata{
				"key1": {"foo", "bar"},
				"key2": {"baz"},
			},
		}
		actual, err := ChainResolvers(
			constantResult{
				Members: []string{"foo", "bar"},
				Metadata: Metadata{
					"key1": {"foo", "bar"},
				},
			},
			constantResult{
				Members: []string{"baz"},
				Metadata: Metadata{
					"key2": {"baz"},
				},
			},
		).ResolveIAMMembers(context.Background())
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual)
	})

	t.Run("multi duplicates", func(t *testing.T) {
		expected := ResolveResult{
			Members: []string{"foo", "bar", "baz"},
			Metadata: Metadata{
				"key1": {"foo", "bar", "baz"},
				"key2": {"bar", "baz"},
			},
		}
		actual, err := ChainResolvers(
			constantResult{
				Members: []string{"foo", "bar"},
				Metadata: Metadata{
					"key1": {"foo", "bar"},
				},
			},
			constantResult{
				Members: []string{"bar", "baz"},
				Metadata: Metadata{
					"key1": {"baz"},
					"key2": {"bar", "baz"},
				},
			},
		).ResolveIAMMembers(context.Background())
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual)
	})

	t.Run("error", func(t *testing.T) {
		actual, err := ChainResolvers(
			constantResult{Members: []string{"foo", "bar"}},
			errorResolver{err: errors.New("boom")},
		).ResolveIAMMembers(context.Background())
		assert.Error(t, err, "boom")
		assert.DeepEqual(t, ResolveResult{}, actual)
	})
}

type constantResult ResolveResult

func (c constantResult) ResolveIAMMembers(context.Context) (ResolveResult, error) {
	return ResolveResult(c), nil
}

type errorResolver struct {
	err error
}

func (e errorResolver) ResolveIAMMembers(context.Context) (ResolveResult, error) {
	return ResolveResult{}, e.err
}
