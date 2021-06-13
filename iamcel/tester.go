package iamcel

import (
	"context"
)

type PermissionTester interface {
	ResourcePermissionTester
	ResourcePermissionsTester
}

// ResourcePermissionTester is an interface for testing the presence for a single resource permission binding.
type ResourcePermissionTester interface {
	TestResourcePermission(
		ctx context.Context, members []string, resource string, permission string,
	) (bool, error)
}

// ResourcePermissionsTester is an interface for testing the presence for multiple resource permission bindings.
type ResourcePermissionsTester interface {
	TestResourcePermissions(
		ctx context.Context, members []string, resourcePermissions map[string]string,
	) (map[string]bool, error)
}
