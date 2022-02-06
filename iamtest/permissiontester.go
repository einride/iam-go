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
func (p *PermissionTester) AllowAll() *PermissionTester {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.allowAll = true
	return p
}

// Allow a permission for a member and resource.
func (p *PermissionTester) Allow(member, permission string, resources ...string) *PermissionTester {
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
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.allowAll {
		return p
	}
	for _, resource := range resources {
		p.allowedPermissions = append(p.allowedPermissions, allowedPermission{
			member: member, resource: resource, permission: permission,
		})
	}
	return p
}

// Reset all allowed permissions.
func (p *PermissionTester) Reset() *PermissionTester {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.allowAll = false
	p.allowedPermissions = p.allowedPermissions[:0]
	return p
}

// TestPermissions implements iamcel.PermissionTester.
func (p *PermissionTester) TestPermissions(
	_ context.Context,
	caller *iamv1.Caller,
	resourcePermissions map[string]string,
) (map[string]bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	result := make(map[string]bool, len(resourcePermissions))
ResourcePermissionLoop:
	for resource, permission := range resourcePermissions {
		if p.allowAll {
			result[resource] = true
			continue ResourcePermissionLoop
		}
		for _, member := range caller.GetMembers() {
			for _, p := range p.allowedPermissions {
				isResource := p.resource == iamresource.Root ||
					p.resource == resource ||
					resourcename.HasParent(resource, p.resource)
				if isResource && p.member == member && p.permission == permission {
					result[resource] = true
					continue ResourcePermissionLoop
				}
			}
		}
		result[resource] = false
	}
	return result, nil
}
