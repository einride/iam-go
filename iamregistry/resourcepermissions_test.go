package iamregistry

import (
	"testing"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/genproto/googleapis/api/annotations"
	"gotest.tools/v3/assert"
)

func TestResourcePermissions(t *testing.T) {
	t.Run("empty registry", func(t *testing.T) {
		resourcePermissions, err := NewResourcePermissions(
			&iamv1.ResourcePermissions{},
			nil,
		)
		assert.NilError(t, err)
		t.Run("root", func(t *testing.T) {
			permission, ok := resourcePermissions.FindPermissionByResourceName("/")
			assert.Equal(t, "", permission)
			assert.Assert(t, !ok)
		})
		t.Run("empty name", func(t *testing.T) {
			permission, ok := resourcePermissions.FindPermissionByResourceName("")
			assert.Equal(t, "", permission)
			assert.Assert(t, !ok)
		})
		t.Run("valid name", func(t *testing.T) {
			permission, ok := resourcePermissions.FindPermissionByResourceName("users/1234")
			assert.Equal(t, "", permission)
			assert.Assert(t, !ok)
		})
	})

	t.Run("non-empty registry", func(t *testing.T) {
		resourcePermissions, err := NewResourcePermissions(
			&iamv1.ResourcePermissions{
				ResourcePermission: []*iamv1.ResourcePermission{
					{Resource: &annotations.ResourceDescriptor{Type: "/"}, Permission: "test.root.examplePermission"},
					{Resource: &annotations.ResourceDescriptor{Type: "example.com/User"}, Permission: "test.users.examplePermission"},
					{Resource: &annotations.ResourceDescriptor{Type: "example.com/Book"}, Permission: "test.books.examplePermission"},
				},
			},
			[]*annotations.ResourceDescriptor{
				{
					Type:    "example.com/User",
					Pattern: []string{"users/{user}", "tenants/{tenant}/users/{user}"},
				},
				{
					Type:    "example.com/Book",
					Pattern: []string{"books/{book}"},
				},
			},
		)
		assert.NilError(t, err)
		t.Run("root", func(t *testing.T) {
			permission, ok := resourcePermissions.FindPermissionByResourceName("/")
			assert.Equal(t, "test.root.examplePermission", permission)
			assert.Assert(t, ok)
		})
		t.Run("empty name", func(t *testing.T) {
			permission, ok := resourcePermissions.FindPermissionByResourceName("")
			assert.Equal(t, "", permission)
			assert.Assert(t, !ok)
		})
		t.Run("valid name", func(t *testing.T) {
			permission, ok := resourcePermissions.FindPermissionByResourceName("users/1234")
			assert.Equal(t, "test.users.examplePermission", permission)
			assert.Assert(t, ok)
		})
		t.Run("multiple patterns", func(t *testing.T) {
			permission, ok := resourcePermissions.FindPermissionByResourceName("tenants/1234/users/1234")
			assert.Equal(t, "test.users.examplePermission", permission)
			assert.Assert(t, ok)
		})
	})

	t.Run("double registration", func(t *testing.T) {
		resourcePermissions, err := NewResourcePermissions(
			&iamv1.ResourcePermissions{
				ResourcePermission: []*iamv1.ResourcePermission{
					{Resource: &annotations.ResourceDescriptor{Type: "example.com/User"}, Permission: "test.users.examplePermission"},
					{Resource: &annotations.ResourceDescriptor{Type: "example.com/Book"}, Permission: "test.books.examplePermission"},
					{Resource: &annotations.ResourceDescriptor{Type: "example.com/User"}, Permission: "test.users.otherPermission"},
				},
			},
			[]*annotations.ResourceDescriptor{
				{
					Type:    "example.com/User",
					Pattern: []string{"users/{user}", "tenants/{tenant}/users/{user}"},
				},
				{
					Type:    "example.com/Book",
					Pattern: []string{"books/{book}"},
				},
			},
		)
		assert.ErrorContains(t, err, "pattern users/{user} already mapped to permission test.users.examplePermission")
		assert.Assert(t, resourcePermissions == nil)
	})

	t.Run("missing resource", func(t *testing.T) {
		resourcePermissions, err := NewResourcePermissions(
			&iamv1.ResourcePermissions{
				ResourcePermission: []*iamv1.ResourcePermission{
					{Resource: &annotations.ResourceDescriptor{Type: "/"}, Permission: "test.root.examplePermission"},
					{Resource: &annotations.ResourceDescriptor{Type: "example.com/User"}, Permission: "test.users.examplePermission"},
					{Resource: &annotations.ResourceDescriptor{Type: "example.com/Book"}, Permission: "test.books.examplePermission"},
				},
			},
			[]*annotations.ResourceDescriptor{
				{
					Type:    "example.com/User",
					Pattern: []string{"users/{user}", "tenants/{tenant}/users/{user}"},
				},
			},
		)
		assert.ErrorContains(t, err, "found no resource with type example.com/Book")
		assert.Assert(t, resourcePermissions == nil)
	})
}
