package iamtest

import (
	"context"
	"fmt"
	"sync"

	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iammember"
	"go.einride.tech/iam/iampermission"
	"go.einride.tech/iam/iamresource"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// PermissionTester is a mock permission tester.
type PermissionTester struct {
	mu                 sync.Mutex
	allowAll           bool
	allowedPermissions []allowedPermission
}

// NewPermissionTester creates a new mock permission tester.
func NewPermissionTester() *PermissionTester {
	return &PermissionTester{}
}

type allowedPermission struct {
	member, resource, permission string
}

// AllowAll allows all testable permissions for all members and all resources.
func (m *PermissionTester) AllowAll() *PermissionTester {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.allowAll = true
	return m
}

// Allow a permission for a member and resource.
func (m *PermissionTester) Allow(member, permission string, resources ...string) *PermissionTester {
	if err := iammember.Validate(member); err != nil {
		panic(fmt.Errorf("permission tester: allow on invalid member %s: %w", member, err))
	}
	if err := iampermission.Validate(permission); err != nil {
		panic(fmt.Errorf("permission tester: allow on invalid permission %s: %w", permission, err))
	}
	for _, resource := range resources {
		if err := resourcename.Validate(resource); err != nil {
			panic(fmt.Errorf("permission tester: allow on invalid resource %s: %w", resource, err))
		}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.allowAll {
		return m
	}
	for _, resource := range resources {
		m.allowedPermissions = append(m.allowedPermissions, allowedPermission{
			member: member, resource: resource, permission: permission,
		})
	}
	return m
}

// Reset all allowed permissions.
func (m *PermissionTester) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.allowAll = false
	m.allowedPermissions = m.allowedPermissions[:0]
}

// TestPermissions implements iamcel.PermissionTester.
func (m *PermissionTester) TestPermissions(
	_ context.Context,
	caller *iamv1.Caller,
	resourcePermissions map[string]string,
) (map[string]bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make(map[string]bool, len(resourcePermissions))
ResourcePermissionLoop:
	for resource, permission := range resourcePermissions {
		if m.allowAll {
			result[resource] = true
			continue ResourcePermissionLoop
		}
		for _, member := range caller.GetMembers() {
			for _, p := range m.allowedPermissions {
				isResource :=
					p.resource == iamresource.Root ||
						p.resource == resource ||
						resourcename.HasParent(resource, p.resource)
				if isResource && p.member == member && p.permission == permission {
					result[resource] = true
					continue ResourcePermissionLoop
				}
			}
		}
	}
	return result, nil
}
