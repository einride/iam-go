package iamcaller

import (
	"context"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// Resolver resolves the IAM caller identity for the current context.
type Resolver interface {
	ResolveCaller(context.Context) (*iamv1.Caller, error)
}
