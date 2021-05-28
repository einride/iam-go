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

func (c chainResolver) ResolveIAMMembers(ctx context.Context) ([]string, Metadata, error) {
	resultMembers := make([]string, 0, len(c.resolvers))
	resultMetadata := make(map[string][]string, len(c.resolvers))
	for _, resolver := range c.resolvers {
		members, metadata, err := resolver.ResolveIAMMembers(ctx)
		if err != nil {
			return nil, nil, err
		}
		for _, member := range members {
			var hasMember bool
			for _, resultMember := range resultMembers {
				if member == resultMember {
					hasMember = true
					break
				}
			}
			if !hasMember {
				resultMembers = append(resultMembers, member)
			}
		}
		for key, keyMembers := range metadata {
			for _, keyMember := range keyMembers {
				var hasKeyMember bool
				for _, resultMember := range resultMetadata[key] {
					if keyMember == resultMember {
						hasKeyMember = true
						break
					}
				}
				if !hasKeyMember {
					resultMetadata[key] = append(resultMetadata[key], keyMember)
				}
			}
		}
	}
	return resultMembers, resultMetadata, nil
}
