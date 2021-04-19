package spaniam

import (
	"context"
	"testing"

	"cloud.google.com/go/spanner"
	"go.einride.tech/spanner-aip/spantest"
	"google.golang.org/genproto/googleapis/iam/admin/v1"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestServer(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	var cancel func()
	if deadline, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(ctx, deadline)
		t.Cleanup(cancel)
	}
	fx := spantest.NewEmulatorFixture(t)
	newDatabase := func() *spanner.Client {
		return fx.NewDatabaseFromDDLFiles(t, "./schema.sql")
	}

	const (
		user1 = "email:user1@example.com"
		user2 = "email:user2@example.com"
		user3 = "email:user3@example.com"
	)
	roles := []*admin.Role{
		{
			Name:                "roles/admin",
			Title:               "Admin",
			Description:         "Test admin",
			IncludedPermissions: []string{"test.create", "test.get", "test.update", "test.delete"},
		},
		{
			Name:                "roles/user",
			Title:               "User",
			Description:         "Test user",
			IncludedPermissions: []string{"test.create", "test.get", "test.update"},
		},
		{
			Name:                "roles/viewer",
			Title:               "User",
			Description:         "Test user",
			IncludedPermissions: []string{"test.get"},
		},
	}

	t.Run("get non-existent returns empty policy", func(t *testing.T) {
		t.Parallel()
		server, err := NewServer(newDatabase(), ServerConfig{
			BuiltInRoles: roles,
			MemberFn: func(ctx context.Context) (string, error) {
				return user1, nil
			},
			ErrorHook: func(ctx context.Context, err error) {
				t.Log(err)
			},
		})
		assert.NilError(t, err)
		expected := &iam.Policy{
			Etag: []byte("W/0-00000000"),
		}
		actual, err := server.GetIamPolicy(ctx, &iam.GetIamPolicyRequest{Resource: "resources/1"})
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual, protocmp.Transform())
	})

	t.Run("set", func(t *testing.T) {
		t.Parallel()
		server, err := NewServer(newDatabase(), ServerConfig{
			BuiltInRoles: roles,
			MemberFn: func(ctx context.Context) (string, error) {
				return user1, nil
			},
			ErrorHook: func(ctx context.Context, err error) {
				t.Log(err)
			},
		})
		assert.NilError(t, err)
		policy := &iam.Policy{
			Bindings: []*iam.Binding{
				{Role: "roles/admin", Members: []string{user1, user2}},
				{Role: "roles/user", Members: []string{user3}},
			},
			Etag: []byte("W/0-00000000"),
		}
		expected := &iam.Policy{
			Bindings: []*iam.Binding{
				{Role: "roles/admin", Members: []string{user1, user2}},
				{Role: "roles/user", Members: []string{user3}},
			},
			Etag: []byte("W/104-14AAE092"),
		}
		actual, err := server.SetIamPolicy(ctx, &iam.SetIamPolicyRequest{
			Resource: "resources/1",
			Policy:   policy,
		})
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual, protocmp.Transform())
	})

	t.Run("set stale", func(t *testing.T) {
		t.Parallel()
		server, err := NewServer(newDatabase(), ServerConfig{
			BuiltInRoles: roles,
			MemberFn: func(ctx context.Context) (string, error) {
				return user1, nil
			},
			ErrorHook: func(ctx context.Context, err error) {
				t.Log(err)
			},
		})
		assert.NilError(t, err)
		policy := &iam.Policy{
			Bindings: []*iam.Binding{
				{Role: "roles/admin", Members: []string{user1, user2}},
				{Role: "roles/user", Members: []string{user3}},
			},
			// stale etag
			Etag: []byte("W/1234"),
		}
		actual, err := server.SetIamPolicy(ctx, &iam.SetIamPolicyRequest{
			Resource: "resources/1",
			Policy:   policy,
		})
		assert.Equal(t, codes.Aborted, status.Code(err))
		assert.Assert(t, actual == nil)
	})

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		server, err := NewServer(newDatabase(), ServerConfig{
			BuiltInRoles: roles,
			MemberFn: func(ctx context.Context) (string, error) {
				return user1, nil
			},
			ErrorHook: func(ctx context.Context, err error) {
				t.Log(err)
			},
		})
		assert.NilError(t, err)
		policy := &iam.Policy{
			Bindings: []*iam.Binding{
				{Role: "roles/admin", Members: []string{user1, user2}},
				{Role: "roles/user", Members: []string{user3}},
			},
			Etag: []byte("W/0-00000000"),
		}
		expected := &iam.Policy{
			Bindings: []*iam.Binding{
				{Role: "roles/admin", Members: []string{user1, user2}},
				{Role: "roles/user", Members: []string{user3}},
			},
			Etag: []byte("W/104-14AAE092"),
		}
		actual, err := server.SetIamPolicy(ctx, &iam.SetIamPolicyRequest{
			Resource: "resources/1",
			Policy:   policy,
		})
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual, protocmp.Transform())
		got, err := server.GetIamPolicy(ctx, &iam.GetIamPolicyRequest{Resource: "resources/1"})
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, got, protocmp.Transform())
	})

	t.Run("set and get other resource", func(t *testing.T) {
		t.Parallel()
		server, err := NewServer(newDatabase(), ServerConfig{
			BuiltInRoles: roles,
			MemberFn: func(ctx context.Context) (string, error) {
				return user1, nil
			},
			ErrorHook: func(ctx context.Context, err error) {
				t.Log(err)
			},
		})
		assert.NilError(t, err)
		policy := &iam.Policy{
			Bindings: []*iam.Binding{
				{Role: "roles/admin", Members: []string{user1, user2}},
				{Role: "roles/user", Members: []string{user3}},
			},
			Etag: []byte("W/0-00000000"),
		}
		expected := &iam.Policy{
			Bindings: []*iam.Binding{
				{Role: "roles/admin", Members: []string{user1, user2}},
				{Role: "roles/user", Members: []string{user3}},
			},
			Etag: []byte("W/104-14AAE092"),
		}
		actual, err := server.SetIamPolicy(ctx, &iam.SetIamPolicyRequest{
			Resource: "resources/1",
			Policy:   policy,
		})
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual, protocmp.Transform())
		emptyPolicy := &iam.Policy{
			Etag: []byte("W/0-00000000"),
		}
		got, err := server.GetIamPolicy(ctx, &iam.GetIamPolicyRequest{Resource: "resources/2"})
		assert.NilError(t, err)
		assert.DeepEqual(t, emptyPolicy, got, protocmp.Transform())
	})

	t.Run("test no permissions", func(t *testing.T) {
		t.Parallel()
		server, err := NewServer(newDatabase(), ServerConfig{
			BuiltInRoles: roles,
			MemberFn: func(ctx context.Context) (string, error) {
				return user1, nil
			},
			ErrorHook: func(ctx context.Context, err error) {
				t.Log(err)
			},
		})
		assert.NilError(t, err)
		response, err := server.TestIamPermissions(ctx, &iam.TestIamPermissionsRequest{
			Resource:    "resources/1",
			Permissions: []string{"test.create", "test.get", "test.update", "test.delete"},
		})
		assert.NilError(t, err)
		assert.Assert(t, cmp.Len(response.Permissions, 0))
	})

	t.Run("test all permissions", func(t *testing.T) {
		t.Parallel()
		server, err := NewServer(newDatabase(), ServerConfig{
			BuiltInRoles: roles,
			MemberFn: func(ctx context.Context) (string, error) {
				return user1, nil
			},
			ErrorHook: func(ctx context.Context, err error) {
				t.Log(err)
			},
		})
		assert.NilError(t, err)
		policy := &iam.Policy{
			Bindings: []*iam.Binding{
				{Role: "roles/admin", Members: []string{user1}},
			},
			Etag: []byte("W/0-00000000"),
		}
		_, err = server.SetIamPolicy(ctx, &iam.SetIamPolicyRequest{
			Resource: "resources/1",
			Policy:   policy,
		})
		assert.NilError(t, err)
		permissions := []string{"test.create", "test.get", "test.update", "test.delete"}
		response, err := server.TestIamPermissions(ctx, &iam.TestIamPermissionsRequest{
			Resource:    "resources/1",
			Permissions: permissions,
		})
		assert.NilError(t, err)
		assert.DeepEqual(t, permissions, response.Permissions)
	})

	t.Run("test some permissions", func(t *testing.T) {
		t.Parallel()
		server, err := NewServer(newDatabase(), ServerConfig{
			BuiltInRoles: roles,
			MemberFn: func(ctx context.Context) (string, error) {
				return user1, nil
			},
			ErrorHook: func(ctx context.Context, err error) {
				t.Log(err)
			},
		})
		assert.NilError(t, err)
		policy := &iam.Policy{
			Bindings: []*iam.Binding{
				{Role: "roles/viewer", Members: []string{user1}},
			},
			Etag: []byte("W/0-00000000"),
		}
		_, err = server.SetIamPolicy(ctx, &iam.SetIamPolicyRequest{
			Resource: "resources/1",
			Policy:   policy,
		})
		assert.NilError(t, err)
		permissions := []string{"test.create", "test.get", "test.update", "test.delete"}
		expected := []string{"test.get"}
		response, err := server.TestIamPermissions(ctx, &iam.TestIamPermissionsRequest{
			Resource:    "resources/1",
			Permissions: permissions,
		})
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, response.Permissions)
	})

	t.Run("test permissions wrong user", func(t *testing.T) {
		t.Parallel()
		server, err := NewServer(newDatabase(), ServerConfig{
			BuiltInRoles: roles,
			MemberFn: func(ctx context.Context) (string, error) {
				return user2, nil
			},
			ErrorHook: func(ctx context.Context, err error) {
				t.Log(err)
			},
		})
		assert.NilError(t, err)
		policy := &iam.Policy{
			Bindings: []*iam.Binding{
				{Role: "roles/admin", Members: []string{user1}},
			},
			Etag: []byte("W/0-00000000"),
		}
		_, err = server.SetIamPolicy(ctx, &iam.SetIamPolicyRequest{
			Resource: "resources/1",
			Policy:   policy,
		})
		assert.NilError(t, err)
		permissions := []string{"test.create", "test.get", "test.update", "test.delete"}
		response, err := server.TestIamPermissions(ctx, &iam.TestIamPermissionsRequest{
			Resource:    "resources/1",
			Permissions: permissions,
		})
		assert.NilError(t, err)
		assert.Assert(t, cmp.Len(response.Permissions, 0))
	})
}
