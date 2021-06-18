package iamcel

import (
	"context"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// PermissionTester is an interface for testing IAM permissions.
type PermissionTester interface {
	TestPermissions(
		ctx context.Context, caller *iamv1.Caller, resourcePermissions map[string]string,
	) (map[string]bool, error)
}
