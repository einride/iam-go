package iammember

import "context"

// Resolver resolves the IAM member identifiers for a caller context.
type Resolver interface {
	ResolveIAMMembers(context.Context) (ResolveResult, error)
}

// Metadata contains IAM members partitioned by which gRPC metadata key they were resolved from.
type Metadata map[string][]string

// ResolveResult is the output from a Resolver.
type ResolveResult struct {
	// Members are the resolved IAM members.
	Members []string
	// Metadata are the resolved IAM members partitioned by which metadata key they were resolved from.
	Metadata Metadata
}

// Add a member resolved from the provided metadata key.
func (r *ResolveResult) Add(key string, member string) {
	var hasMember bool
	for _, existingMember := range r.Members {
		if member == existingMember {
			hasMember = true
			break
		}
	}
	if !hasMember {
		r.Members = append(r.Members, member)
	}
	var hasMetadataMember bool
	for _, existingMetadataMember := range r.Metadata[key] {
		if member == existingMetadataMember {
			hasMetadataMember = true
			break
		}
	}
	if !hasMetadataMember {
		if r.Metadata == nil {
			r.Metadata = make(Metadata)
		}
		r.Metadata[key] = append(r.Metadata[key], member)
	}
}

// AddAll adds all the resolved members from another ResolveResult.
func (r *ResolveResult) AddAll(other ResolveResult) {
	// Add ordered members first to maintain order.
	for _, member := range other.Members {
		var hasMember bool
		for _, existingMember := range r.Members {
			if member == existingMember {
				hasMember = true
				break
			}
		}
		if !hasMember {
			r.Members = append(r.Members, member)
		}
	}
	// Add metadata.
	for key, members := range other.Metadata {
		for _, member := range members {
			r.Add(key, member)
		}
	}
}
