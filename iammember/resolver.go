package iammember

import (
	"context"

	"go.einride.tech/iam/iamjwt"
)

// Resolver resolves the IAM member identifiers for a caller context.
type Resolver interface {
	ResolveIAMMembers(context.Context) (ResolveResult, error)
}

// ResolveResult is the output from a Resolver.
type ResolveResult struct {
	// Metadata are the resolved IAM members partitioned by which metadata key they were resolved from.
	Metadata Metadata
}

// Metadata is a map from metadata keys to IAM members resolved from the metadata values.
type Metadata map[string]MetadataValue

// Members returns the set of all unique members resolved from all metadata keys.
func (r *ResolveResult) Members() []string {
	var size int
	for _, value := range r.Metadata {
		size = +len(value.Members)
	}
	if size == 0 {
		return nil
	}
	result := make([]string, 0, size)
	for _, value := range r.Metadata {
	MemberLoop:
		for _, member := range value.Members {
			for _, existingMember := range result {
				if member == existingMember {
					continue MemberLoop
				}
			}
			result = append(result, member)
		}
	}
	return result
}

// Add a metadata key and resolved metadata value to the result.
func (r *ResolveResult) Add(key string, value MetadataValue) {
	if r.Metadata == nil {
		r.Metadata = make(map[string]MetadataValue)
	}
	r.Metadata[key] = value
}

// MetadataValue is the resolve result from a single metatadata key.
type MetadataValue struct {
	// JWT is the JWT token parsed from the metadata value, if any.
	JWT *iamjwt.Token
	// Members are the members resolved from the metadata value.
	Members []string
}
