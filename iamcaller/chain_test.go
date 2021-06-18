package iamcaller

import (
	"context"
	"errors"
	"testing"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/protobuf/testing/protocmp"
	"gotest.tools/v3/assert"
)

func TestChainResolvers(t *testing.T) {
	t.Run("no resolvers", func(t *testing.T) {
		result, err := ChainResolvers().ResolveCaller(context.Background())
		assert.NilError(t, err)
		assert.DeepEqual(t, &iamv1.Caller{}, result, protocmp.Transform())
	})

	t.Run("single", func(t *testing.T) {
		expected := &iamv1.Caller{
			Members: []string{"test:bar", "test:foo"},
			Metadata: map[string]*iamv1.Caller_Metadata{
				"key1": {Members: []string{"test:foo"}},
				"key2": {Members: []string{"test:bar"}},
			},
		}
		actual, err := ChainResolvers(constant(expected)).ResolveCaller(context.Background())
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual, protocmp.Transform())
	})

	t.Run("multi", func(t *testing.T) {
		expected := &iamv1.Caller{
			Members: []string{"test:bar", "test:baz", "test:foo"},
			Metadata: map[string]*iamv1.Caller_Metadata{
				"key1": {Members: []string{"test:foo", "test:bar"}},
				"key2": {Members: []string{"test:baz"}},
			},
		}
		actual, err := ChainResolvers(
			constant(&iamv1.Caller{
				Members: []string{"test:foo", "test:bar"},
				Metadata: map[string]*iamv1.Caller_Metadata{
					"key1": {Members: []string{"test:foo", "test:bar"}},
				},
			}),
			constant(&iamv1.Caller{
				Members: []string{"test:baz"},
				Metadata: map[string]*iamv1.Caller_Metadata{
					"key2": {Members: []string{"test:baz"}},
				},
			}),
		).ResolveCaller(context.Background())
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual, protocmp.Transform())
	})

	t.Run("multi duplicates", func(t *testing.T) {
		expected := &iamv1.Caller{
			Members: []string{"test:bar", "test:baz", "test:foo"},
			Metadata: map[string]*iamv1.Caller_Metadata{
				"key1": {Members: []string{"test:bar"}},
				"key2": {Members: []string{"test:bar", "test:baz"}},
			},
		}
		actual, err := ChainResolvers(
			constant(&iamv1.Caller{
				Members: []string{"test:foo", "test:bar"},
				Metadata: map[string]*iamv1.Caller_Metadata{
					"key1": {Members: []string{"test:foo", "test:bar"}},
				},
			}),
			constant(&iamv1.Caller{
				Members: []string{"test:bar", "test:baz"},
				Metadata: map[string]*iamv1.Caller_Metadata{
					"key1": {Members: []string{"test:bar"}},
					"key2": {Members: []string{"test:bar", "test:baz"}},
				},
			}),
		).ResolveCaller(context.Background())
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual, protocmp.Transform())
	})

	t.Run("error", func(t *testing.T) {
		actual, err := ChainResolvers(
			constant(&iamv1.Caller{
				Members: []string{"test:baz", "test:bar"},
				Metadata: map[string]*iamv1.Caller_Metadata{
					"key1": {Members: []string{"test:baz"}},
					"key2": {Members: []string{"test:bar", "test:baz"}},
				},
			}),
			errorResolver{err: errors.New("boom")},
		).ResolveCaller(context.Background())
		assert.Error(t, err, "boom")
		assert.Assert(t, actual == nil)
	})
}

type constantResolver struct {
	caller *iamv1.Caller
}

func constant(caller *iamv1.Caller) Resolver {
	return constantResolver{caller: caller}
}

func (c constantResolver) ResolveCaller(context.Context) (*iamv1.Caller, error) {
	return c.caller, nil
}

type errorResolver struct {
	err error
}

func (e errorResolver) ResolveCaller(context.Context) (*iamv1.Caller, error) {
	return nil, e.err
}
