package iammember

import "context"

// Resolver resolves the IAM member identifiers for a caller context.
type Resolver interface {
	ResolveIAMMembers(context.Context) (context.Context, []string, error)
}
