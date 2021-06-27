package iamtest

import (
	"context"
	"testing"

	"gotest.tools/v3/assert"
)

func TestPermissionTester(t *testing.T) {
	t.Run("allow none", func(t *testing.T) {
		const (
			member     = "email:foo@example.com"
			resource   = "resources/1234"
			permission = "test.resources.foo"
		)
		mock := NewPermissionTester()
		result, err := mock.TestPermissions(context.Background(), NewCaller(member), map[string]string{
			resource: permission,
		})
		assert.NilError(t, err)
		assert.DeepEqual(t, map[string]bool{}, result)
	})

	t.Run("allow all", func(t *testing.T) {
		const (
			member     = "email:foo@example.com"
			resource   = "resources/1234"
			permission = "test.resources.foo"
		)
		mock := NewPermissionTester()
		mock.AllowAll()
		result, err := mock.TestPermissions(context.Background(), NewCaller(member), map[string]string{
			resource: permission,
		})
		assert.NilError(t, err)
		assert.DeepEqual(t, map[string]bool{
			resource: true,
		}, result)
	})

	t.Run("allow some", func(t *testing.T) {
		const (
			member     = "email:foo@example.com"
			resource1  = "resources/1234"
			resource2  = "resources/5678"
			permission = "test.resources.foo"
		)
		mock := NewPermissionTester()
		mock.Allow(member, permission, resource2)
		result, err := mock.TestPermissions(context.Background(), NewCaller(member), map[string]string{
			resource1: permission,
			resource2: permission,
		})
		assert.NilError(t, err)
		assert.DeepEqual(t, map[string]bool{
			resource2: true,
		}, result)
	})

	t.Run("reset", func(t *testing.T) {
		const (
			member     = "email:foo@example.com"
			resource   = "resources/1234"
			permission = "test.resources.foo"
		)
		mock := NewPermissionTester()
		mock.AllowAll()
		mock.Reset()
		result, err := mock.TestPermissions(context.Background(), NewCaller(member), map[string]string{
			resource: permission,
		})
		assert.NilError(t, err)
		assert.DeepEqual(t, map[string]bool{}, result)
	})
}
