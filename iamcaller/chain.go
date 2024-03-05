package iamcaller

import (
	"context"
	"fmt"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// ChainResolvers creates a single resolver out of a chain of many resolvers.
//
// The resulting resolved caller data will be the union of the caller data resolved by each resolver.
//
// If any resolver returns an error, that error is immediately returned and no further resolvers are called.
//
// If multiple resolvers resolve the same metadata key, the only last encountered metadata value will be kept.
func ChainResolvers(resolvers ...Resolver) Resolver {
	return chainResolver{resolvers: resolvers}
}

type chainResolver struct {
	resolvers []Resolver
}

// ResolveCaller implements Resolver.
func (c chainResolver) ResolveCaller(ctx context.Context) (*iamv1.Caller, error) {
	var result iamv1.Caller
	for i, resolver := range c.resolvers {
		nextCaller, err := resolver.ResolveCaller(ctx)
		if err != nil {
			return nil, err
		}
		if err := Validate(nextCaller); err != nil {
			return nil, fmt.Errorf("chain callers: resolver %d: %w", i, err)
		}
		for key, value := range nextCaller.GetMetadata() {
			Add(&result, key, value)
		}
		// TODO: Remove this when CEL-Go supports async functions with context arguments.
		if result.GetContext() == nil && nextCaller.GetContext() != nil {
			result.Context = nextCaller.GetContext()
		}
	}
	return &result, nil
}
