package iamauthz

import "go.einride.tech/iam/iamcel"

// PermissionTester is an interface for the permission testing functionality needed for authorization.
type PermissionTester interface {
	iamcel.ResourcePermissionTester
	iamcel.ResourcePermissionsTester
}
