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
			Metadata: Metadata{
				"key1": {Members: []string{"foo"}},
				"key2": {Members: []string{"bar"}},
			},
		}
		actual, err := ChainResolvers(constantResult(expected)).ResolveIAMMembers(context.Background())
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual)
	})

	t.Run("multi", func(t *testing.T) {
		expected := ResolveResult{
			Metadata: Metadata{
				"key1": {Members: []string{"foo", "bar"}},
				"key2": {Members: []string{"baz"}},
			},
		}
		actual, err := ChainResolvers(
			constantResult{
				Metadata: Metadata{
					"key1": {Members: []string{"foo", "bar"}},
				},
			},
			constantResult{
				Metadata: Metadata{
					"key2": {Members: []string{"baz"}},
				},
			},
		).ResolveIAMMembers(context.Background())
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual)
	})

	t.Run("multi duplicates", func(t *testing.T) {
		expected := ResolveResult{
			Metadata: Metadata{
				"key1": {Members: []string{"bar"}},
				"key2": {Members: []string{"bar", "baz"}},
			},
		}
		actual, err := ChainResolvers(
			constantResult{
				Metadata: Metadata{
					"key1": {Members: []string{"foo", "bar"}},
				},
			},
			constantResult{
				Metadata: Metadata{
					"key1": {Members: []string{"bar"}},
					"key2": {Members: []string{"bar", "baz"}},
				},
			},
		).ResolveIAMMembers(context.Background())
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual)
	})

	t.Run("error", func(t *testing.T) {
		actual, err := ChainResolvers(
			constantResult{
				Metadata: Metadata{
					"key1": {Members: []string{"baz"}},
					"key2": {Members: []string{"bar", "baz"}},
				},
			},
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
