package iamtest

import (
	"context"

	"go.einride.tech/iam/iamcaller"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// WithMembers returns a new context with resolved IAM members for test purposes.
func WithMembers(ctx context.Context, members ...string) context.Context {
	return iamcaller.WithResolvedContext(ctx, NewCaller(members...))
}

// NewCaller creates a new caller for test purposes.
func NewCaller(members ...string) *iamv1.Caller {
	return &iamv1.Caller{
		Members: members,
		Metadata: map[string]*iamv1.Caller_Metadata{
			"iamtest": {Members: members},
		},
	}
}
