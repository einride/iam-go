package iammember

import (
	"context"
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func TestChainResolvers(t *testing.T) {
	t.Run("no resolvers", func(t *testing.T) {
		ctx, members, err := ChainResolvers().ResolveIAMMembers(context.Background())
		assert.Equal(t, ctx, context.Background())
		assert.Assert(t, members == nil)
		assert.NilError(t, err)
	})

	t.Run("single", func(t *testing.T) {
		expected := []string{"foo", "bar"}
		ctx, actual, err := ChainResolvers(constantResolver{expected}).ResolveIAMMembers(context.Background())
		assert.Equal(t, ctx, context.Background())
		assert.DeepEqual(t, expected, actual)
		assert.NilError(t, err)
	})

	t.Run("multi", func(t *testing.T) {
		expected := []string{"foo", "bar", "baz"}
		ctx, actual, err := ChainResolvers(
			constantResolver{members: []string{"foo", "bar"}},
			constantResolver{members: []string{"baz"}},
		).ResolveIAMMembers(context.Background())
		assert.Equal(t, ctx, context.Background())
		assert.DeepEqual(t, expected, actual)
		assert.NilError(t, err)
	})

	t.Run("multi duplicates", func(t *testing.T) {
		expected := []string{"foo", "bar", "baz"}
		ctx, actual, err := ChainResolvers(
			constantResolver{members: []string{"foo", "bar"}},
			constantResolver{members: []string{"bar", "baz"}},
		).ResolveIAMMembers(context.Background())
		assert.Equal(t, ctx, context.Background())
		assert.DeepEqual(t, expected, actual)
		assert.NilError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		ctx, actual, err := ChainResolvers(
			constantResolver{members: []string{"foo", "bar"}},
			errorResolver{err: errors.New("boom")},
		).ResolveIAMMembers(context.Background())
		assert.Assert(t, ctx == nil)
		assert.Assert(t, actual == nil)
		assert.Error(t, err, "boom")
	})
}

type constantResolver struct {
	members []string
}

func (c constantResolver) ResolveIAMMembers(ctx context.Context) (context.Context, []string, error) {
	return ctx, c.members, nil
}

type errorResolver struct {
	err error
}

func (e errorResolver) ResolveIAMMembers(ctx context.Context) (context.Context, []string, error) {
	return nil, nil, e.err
}
