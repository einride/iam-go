package iammember

import "context"

// ChainResolvers creates a single resolver out of a chain of many resolvers.
//
// The resulting resolved members will be the union of the members resolved by each resolver.
//
// If any resolver returns an error, that error is immediately returned and no further resolvers are called.
func ChainResolvers(resolvers ...Resolver) Resolver {
	return chainResolver{resolvers: resolvers}
}

type chainResolver struct {
	resolvers []Resolver
}

// ResolveIAMMembers implements Resolver.
func (c chainResolver) ResolveIAMMembers(ctx context.Context) (ResolveResult, error) {
	var result ResolveResult
	for _, resolver := range c.resolvers {
		nextResult, err := resolver.ResolveIAMMembers(ctx)
		if err != nil {
			return ResolveResult{}, err
		}
		result.AddAll(nextResult)
	}
	return result, nil
}
