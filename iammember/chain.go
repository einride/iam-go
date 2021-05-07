package iammember

import "context"

// ChainResolvers creates a single resolver out of a chain of many resolvers.
//
// The resulting resolved members will be the union of the members resolved by each resolver.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainResolvers(one, two, three) will execute one before two before three, and three
// will see context changes of one and two.
//
// If any resolver returns an error, that error is immediately returned and no further resolvers are called.
func ChainResolvers(resolvers ...Resolver) Resolver {
	return chainResolver{resolvers: resolvers}
}

type chainResolver struct {
	resolvers []Resolver
}

func (c chainResolver) ResolveIAMMembers(ctx context.Context) (context.Context, []string, error) {
	var result, members []string
	var err error
	for _, resolver := range c.resolvers {
		ctx, members, err = resolver.ResolveIAMMembers(ctx)
		if err != nil {
			return nil, nil, err
		}
		for _, member := range members {
			var hasMember bool
			for _, resultMember := range result {
				if member == resultMember {
					hasMember = true
					break
				}
			}
			if !hasMember {
				result = append(result, member)
			}
		}
	}
	return ctx, result, nil
}
