package iammember

import "context"

// Resolver resolves the IAM member identifiers for a caller context.
type Resolver interface {
	ResolveIAMMembers(context.Context) ([]string, Metadata, error)
}

// Metadata contains IAM members partitioned by which gRPC metadata key they were resolved from.
type Metadata map[string][]string
