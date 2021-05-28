package iammember

import (
	"context"
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func TestChainResolvers(t *testing.T) {
	t.Run("no resolvers", func(t *testing.T) {
		actualMembers, actualMetadata, err := ChainResolvers().ResolveIAMMembers(context.Background())
		assert.Assert(t, len(actualMembers) == 0)
		assert.Assert(t, len(actualMetadata) == 0)
		assert.NilError(t, err)
	})

	t.Run("single", func(t *testing.T) {
		expectedMembers := []string{"foo", "bar"}
		expectedMetadata := Metadata{
			"key1": {"foo"},
			"key2": {"bar"},
		}
		actualMembers, actualMetadata, err := ChainResolvers(constantResolver{
			members:  expectedMembers,
			metadata: expectedMetadata,
		}).ResolveIAMMembers(context.Background())
		assert.DeepEqual(t, expectedMembers, actualMembers)
		assert.DeepEqual(t, expectedMetadata, actualMetadata)
		assert.NilError(t, err)
	})

	t.Run("multi", func(t *testing.T) {
		expectedMembers := []string{"foo", "bar", "baz"}
		expectedMetadata := Metadata{
			"key1": {"foo", "bar"},
			"key2": {"baz"},
		}
		actualMembers, actualMetadata, err := ChainResolvers(
			constantResolver{
				members: []string{"foo", "bar"},
				metadata: Metadata{
					"key1": {"foo", "bar"},
				},
			},
			constantResolver{
				members: []string{"baz"},
				metadata: Metadata{
					"key2": {"baz"},
				},
			},
		).ResolveIAMMembers(context.Background())
		assert.DeepEqual(t, expectedMembers, actualMembers)
		assert.DeepEqual(t, expectedMetadata, actualMetadata)
		assert.NilError(t, err)
	})

	t.Run("multi duplicates", func(t *testing.T) {
		expectedMembers := []string{"foo", "bar", "baz"}
		expectedMetadata := Metadata{
			"key1": {"foo", "bar", "baz"},
			"key2": {"bar", "baz"},
		}
		actualMembers, actualMetadata, err := ChainResolvers(
			constantResolver{
				members: []string{"foo", "bar"},
				metadata: Metadata{
					"key1": {"foo", "bar"},
				},
			},
			constantResolver{
				members: []string{"bar", "baz"},
				metadata: Metadata{
					"key1": {"baz"},
					"key2": {"bar", "baz"},
				},
			},
		).ResolveIAMMembers(context.Background())
		assert.DeepEqual(t, expectedMembers, actualMembers)
		assert.DeepEqual(t, expectedMetadata, actualMetadata)
		assert.NilError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		actualMembers, actualMetadata, err := ChainResolvers(
			constantResolver{members: []string{"foo", "bar"}},
			errorResolver{err: errors.New("boom")},
		).ResolveIAMMembers(context.Background())
		assert.Assert(t, actualMembers == nil)
		assert.Assert(t, actualMetadata == nil)
		assert.Error(t, err, "boom")
	})
}

type constantResolver struct {
	members  []string
	metadata Metadata
}

func (c constantResolver) ResolveIAMMembers(ctx context.Context) ([]string, Metadata, error) {
	return c.members, c.metadata, nil
}

type errorResolver struct {
	err error
}

func (e errorResolver) ResolveIAMMembers(ctx context.Context) ([]string, Metadata, error) {
	return nil, nil, e.err
}
