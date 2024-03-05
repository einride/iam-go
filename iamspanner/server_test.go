package iamspanner

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"cloud.google.com/go/iam/admin/apiv1/adminpb"
	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/spanner"
	"go.einride.tech/iam/iamcaller"
	"go.einride.tech/iam/iampolicy"
	"go.einride.tech/iam/iamregistry"
	"go.einride.tech/iam/iamresource"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"go.einride.tech/spanner-aip/spantest"
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
	newDatabase := func(t *testing.T) *spanner.Client {
		return fx.NewDatabaseFromDDLFiles(t, "./schema.sql")
	}
	const (
		user1 = "email:user1@example.com"
		user2 = "email:user2@example.com"
		user3 = "email:user3@example.com"
	)
	roles := []*adminpb.Role{
		{
			Name:        "roles/test.admin",
			Title:       "Admin",
			Description: "Test admin",
			IncludedPermissions: []string{
				"test.resources.create",
				"test.resources.get",
				"test.resources.update",
				"test.resources.delete",
			},
		},
		{
			Name:        "roles/test.user",
			Title:       "User",
			Description: "Test user",
			IncludedPermissions: []string{
				"test.resources.create",
				"test.resources.get",
				"test.resources.update",
			},
		},
		{
			Name:        "roles/test.viewer",
			Title:       "User",
			Description: "Test user",
			IncludedPermissions: []string{
				"test.resources.get",
			},
		},
	}
	rolesRegistry, err := iamregistry.NewRoles(roles...)
	assert.NilError(t, err)

	t.Run("get non-existent returns empty policy", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		expected := &iampb.Policy{
			Etag: []byte("W/0-00000000"),
		}
		actual, err := server.GetIamPolicy(
			withMembers(ctx, user1),
			&iampb.GetIamPolicyRequest{Resource: "resources/1"},
		)
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual, protocmp.Transform())
	})

	t.Run("get invalid resource", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		actual, err := server.GetIamPolicy(
			withMembers(ctx, user1),
			&iampb.GetIamPolicyRequest{Resource: "ice cream is best"},
		)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Assert(t, actual == nil)
	})

	t.Run("get wildcard resource", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		actual, err := server.GetIamPolicy(
			withMembers(ctx, user1),
			&iampb.GetIamPolicyRequest{Resource: "resources/-"},
		)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Assert(t, actual == nil)
	})

	t.Run("set", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		policy := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.admin", Members: []string{user1, user2}},
				{Role: "roles/test.user", Members: []string{user3}},
			},
		}
		expected := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.admin", Members: []string{user1, user2}},
				{Role: "roles/test.user", Members: []string{user3}},
			},
			Etag: []byte("W/114-946EB3AA"),
		}
		actual, err := server.SetIamPolicy(
			withMembers(ctx, user1),
			&iampb.SetIamPolicyRequest{
				Resource: "resources/1",
				Policy:   policy,
			},
		)
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual, protocmp.Transform())
	})

	t.Run("set stale", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		policy := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.admin", Members: []string{user1, user2}},
				{Role: "roles/test.user", Members: []string{user3}},
			},
			// stale etag
			Etag: []byte("W/1234"),
		}
		actual, err := server.SetIamPolicy(
			withMembers(ctx, user1),
			&iampb.SetIamPolicyRequest{
				Resource: "resources/1",
				Policy:   policy,
			},
		)
		assert.Equal(t, codes.Aborted, status.Code(err))
		assert.Assert(t, actual == nil)
	})

	t.Run("set and get", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		policy := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.admin", Members: []string{user1, user2}},
				{Role: "roles/test.user", Members: []string{user3}},
			},
			Etag: []byte("W/0-00000000"),
		}
		expected := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.admin", Members: []string{user1, user2}},
				{Role: "roles/test.user", Members: []string{user3}},
			},
			Etag: []byte("W/114-946EB3AA"),
		}
		actual, err := server.SetIamPolicy(
			withMembers(ctx, user1),
			&iampb.SetIamPolicyRequest{
				Resource: "resources/1",
				Policy:   policy,
			},
		)
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual, protocmp.Transform())
		got, err := server.GetIamPolicy(
			withMembers(ctx, user1),
			&iampb.GetIamPolicyRequest{Resource: "resources/1"},
		)
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, got, protocmp.Transform())
	})

	t.Run("set and get other resource", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		policy := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.admin", Members: []string{user1, user2}},
				{Role: "roles/test.user", Members: []string{user3}},
			},
			Etag: []byte("W/0-00000000"),
		}
		expected := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.admin", Members: []string{user1, user2}},
				{Role: "roles/test.user", Members: []string{user3}},
			},
			Etag: []byte("W/114-946EB3AA"),
		}
		actual, err := server.SetIamPolicy(
			withMembers(ctx, user1),
			&iampb.SetIamPolicyRequest{
				Resource: "resources/1",
				Policy:   policy,
			},
		)
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual, protocmp.Transform())
		emptyPolicy := &iampb.Policy{
			Etag: []byte("W/0-00000000"),
		}
		got, err := server.GetIamPolicy(
			withMembers(ctx, user1),
			&iampb.GetIamPolicyRequest{Resource: "resources/2"},
		)
		assert.NilError(t, err)
		assert.DeepEqual(t, emptyPolicy, got, protocmp.Transform())
	})

	t.Run("set invalid member", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
				ValidateMember: func(s string) error {
					if s == "invalid:member" {
						return fmt.Errorf("invalid member")
					}
					return nil
				},
			},
		)
		assert.NilError(t, err)
		policy := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.admin", Members: []string{user1, "invalid:member"}},
			},
		}
		actual, err := server.SetIamPolicy(
			withMembers(ctx, user1),
			&iampb.SetIamPolicyRequest{
				Resource: "resources/1",
				Policy:   policy,
			},
		)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.ErrorContains(t, err, "invalid member")
		assert.Assert(t, actual == nil)
	})

	t.Run("set policy with duplicate members", func(t *testing.T) {
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		policy := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.user", Members: []string{user3, user3}},
			},
			Etag: []byte("W/0-00000000"),
		}

		_, err = server.SetIamPolicy(
			withMembers(ctx, user1),
			&iampb.SetIamPolicyRequest{
				Resource: "resources/1",
				Policy:   policy,
			},
		)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Error(t, err, "field violation on policy.bindings[0].members[1]: duplicate member")
	})

	t.Run("set policy with duplicate roles", func(t *testing.T) {
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		policy := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.user", Members: []string{user3}},
				{Role: "roles/test.user", Members: []string{user3}},
			},
			Etag: []byte("W/0-00000000"),
		}

		_, err = server.SetIamPolicy(
			withMembers(ctx, user1),
			&iampb.SetIamPolicyRequest{
				Resource: "resources/1",
				Policy:   policy,
			},
		)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Error(t, err, "field violation on policy.bindings[1].role: duplicate role: 'roles/test.user'")
	})

	t.Run("set invalid role", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		policy := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.fooBar", Members: []string{user1}},
			},
		}
		actual, err := server.SetIamPolicy(
			withMembers(ctx, user1),
			&iampb.SetIamPolicyRequest{
				Resource: "resources/1",
				Policy:   policy,
			},
		)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.ErrorContains(t, err, "unknown role")
		assert.Assert(t, actual == nil)
	})

	t.Run("set with transaction functions", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)

		policy := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.admin", Members: []string{user1, user2}},
				{Role: "roles/test.user", Members: []string{user3}},
			},
		}
		expected := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.admin", Members: []string{user1, user2}},
				{Role: "roles/test.user", Members: []string{user3}},
			},
			Etag: []byte("W/114-946EB3AA"),
		}

		var calledOne, calledTwo, calledThree bool
		actual, err := server.SetIamPolicyWithFunctionsInTransaction(
			withMembers(ctx, user1),
			&iampb.SetIamPolicyRequest{
				Resource: "resources/1",
				Policy:   policy,
			},
			func(context.Context, *spanner.ReadWriteTransaction, *iampb.Policy) error {
				calledOne = true
				return nil
			},
			func(context.Context, *spanner.ReadWriteTransaction, *iampb.Policy) error {
				calledTwo = true
				return nil
			},
			func(context.Context, *spanner.ReadWriteTransaction, *iampb.Policy) error {
				calledThree = true
				return nil
			},
		)
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual, protocmp.Transform())
		assert.Equal(t, calledOne, true)
		assert.Equal(t, calledTwo, true)
		assert.Equal(t, calledThree, true)
	})

	t.Run("set with transaction function failing", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		policy := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.admin", Members: []string{user1, user2}},
				{Role: "roles/test.user", Members: []string{user3}},
			},
		}

		actual, err := server.SetIamPolicyWithFunctionsInTransaction(
			withMembers(ctx, user1),
			&iampb.SetIamPolicyRequest{
				Resource: "resources/1",
				Policy:   policy,
			},
			func(context.Context, *spanner.ReadWriteTransaction, *iampb.Policy) error {
				return nil
			},
			func(context.Context, *spanner.ReadWriteTransaction, *iampb.Policy) error {
				return fmt.Errorf("transaction function error")
			},
			func(context.Context, *spanner.ReadWriteTransaction, *iampb.Policy) error {
				return nil
			},
		)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.ErrorContains(t, err, "storage error")
		assert.Assert(t, actual == nil)
	})

	t.Run("set empty", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		expected := &iampb.Policy{
			Bindings: []*iampb.Binding{},
			Etag:     []byte("W/0-00000000"),
		}
		actual, err := server.SetIamPolicy(
			withMembers(ctx, user1),
			&iampb.SetIamPolicyRequest{
				Resource: "resources/1",
				Policy:   nil,
			},
		)
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual, protocmp.Transform())
	})

	t.Run("test no permissions", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		response, err := server.TestIamPermissions(
			withMembers(ctx, user1),
			&iampb.TestIamPermissionsRequest{
				Resource: "resources/1",
				Permissions: []string{
					"test.resources.create",
					"test.resources.get",
					"test.resources.update",
					"test.resources.delete",
				},
			},
		)
		assert.NilError(t, err)
		assert.Assert(t, cmp.Len(response.GetPermissions(), 0))
	})

	t.Run("test all permissions", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		policy := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.admin", Members: []string{user1}},
			},
			Etag: []byte("W/0-00000000"),
		}
		_, err = server.SetIamPolicy(
			withMembers(ctx, user1),
			&iampb.SetIamPolicyRequest{
				Resource: "resources/1",
				Policy:   policy,
			},
		)
		assert.NilError(t, err)
		permissions := []string{
			"test.resources.create",
			"test.resources.get",
			"test.resources.update",
			"test.resources.delete",
		}
		response, err := server.TestIamPermissions(
			withMembers(ctx, user1),
			&iampb.TestIamPermissionsRequest{
				Resource:    "resources/1",
				Permissions: permissions,
			},
		)
		assert.NilError(t, err)
		assert.DeepEqual(t, permissions, response.GetPermissions())
	})

	t.Run("test some permissions", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		policy := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.viewer", Members: []string{user1}},
			},
			Etag: []byte("W/0-00000000"),
		}
		_, err = server.SetIamPolicy(
			withMembers(ctx, user1),
			&iampb.SetIamPolicyRequest{
				Resource: "resources/1",
				Policy:   policy,
			},
		)
		assert.NilError(t, err)
		permissions := []string{
			"test.resources.create",
			"test.resources.get",
			"test.resources.update",
			"test.resources.delete",
		}
		expected := []string{"test.resources.get"}
		response, err := server.TestIamPermissions(
			withMembers(ctx, user1),
			&iampb.TestIamPermissionsRequest{
				Resource:    "resources/1",
				Permissions: permissions,
			},
		)
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, response.GetPermissions())
	})

	t.Run("test permissions different user", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		policy := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.admin", Members: []string{user1}},
			},
			Etag: []byte("W/0-00000000"),
		}
		_, err = server.SetIamPolicy(
			withMembers(ctx, user1),
			&iampb.SetIamPolicyRequest{
				Resource: "resources/1",
				Policy:   policy,
			},
		)
		assert.NilError(t, err)
		permissions := []string{
			"test.resources.create",
			"test.resources.get",
			"test.resources.update",
			"test.resources.delete",
		}
		response, err := server.TestIamPermissions(
			withMembers(ctx, user2),
			&iampb.TestIamPermissionsRequest{
				Resource:    "resources/1",
				Permissions: permissions,
			},
		)
		assert.NilError(t, err)
		assert.Assert(t, cmp.Len(response.GetPermissions(), 0))
	})

	t.Run("test permissions on parent", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		policy := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.viewer", Members: []string{user1}},
			},
			Etag: []byte("W/0-00000000"),
		}
		_, err = server.SetIamPolicy(
			withMembers(ctx, user1),
			&iampb.SetIamPolicyRequest{
				Resource: "parents/1",
				Policy:   policy,
			},
		)
		assert.NilError(t, err)
		permissions := []string{
			"test.resources.create",
			"test.resources.get",
			"test.resources.update",
			"test.resources.delete",
		}
		expected := []string{"test.resources.get"}
		response, err := server.TestIamPermissions(
			withMembers(ctx, user1),
			&iampb.TestIamPermissionsRequest{
				Resource:    "parents/1/resources/1",
				Permissions: permissions,
			},
		)
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, response.GetPermissions())
	})

	t.Run("test permissions on root", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		policy := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{Role: "roles/test.viewer", Members: []string{user1}},
			},
			Etag: []byte("W/0-00000000"),
		}
		_, err = server.SetIamPolicy(
			withMembers(ctx, user1),
			&iampb.SetIamPolicyRequest{
				Resource: iamresource.Root,
				Policy:   policy,
			},
		)
		assert.NilError(t, err)
		permissions := []string{
			"test.resources.create",
			"test.resources.get",
			"test.resources.update",
			"test.resources.delete",
		}
		expected := []string{"test.resources.get"}
		response, err := server.TestIamPermissions(
			withMembers(ctx, user1),
			&iampb.TestIamPermissionsRequest{
				Resource:    "parents/1/resources/1",
				Permissions: permissions,
			},
		)
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, response.GetPermissions())
	})

	t.Run("get role", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		expected, ok := rolesRegistry.FindRoleByName("roles/test.admin")
		assert.Assert(t, ok)
		actual, err := server.GetRole(
			withMembers(ctx, user1),
			&adminpb.GetRoleRequest{
				Name: "roles/test.admin",
			},
		)
		assert.NilError(t, err)
		assert.DeepEqual(t, expected, actual, protocmp.Transform())
	})

	t.Run("list roles", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		expected := make([]*adminpb.Role, 0, rolesRegistry.Count())
		rolesRegistry.RangeRoles(func(role *adminpb.Role) bool {
			expected = append(expected, role)
			return true
		})
		sort.Slice(expected, func(i, j int) bool {
			return expected[i].GetName() < expected[j].GetName()
		})
		actual := make([]*adminpb.Role, 0, rolesRegistry.Count())
		var nextPageToken string
		for {
			response, err := server.ListRoles(
				withMembers(ctx, user1),
				&adminpb.ListRolesRequest{
					PageSize:  1,
					PageToken: nextPageToken,
					View:      adminpb.RoleView_FULL,
				},
			)
			assert.NilError(t, err)
			actual = append(actual, response.GetRoles()...)
			nextPageToken = response.GetNextPageToken()
			if nextPageToken == "" {
				break
			}
		}
		assert.DeepEqual(t, expected, actual, protocmp.Transform())
	})

	t.Run("read and write", func(t *testing.T) {
		t.Parallel()
		server, err := NewIAMServer(
			newDatabase(t),
			roles,
			iamcaller.FromContextResolver(),
			ServerConfig{
				ErrorHook: func(_ context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		expected := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{
					Role:    "roles/test.admin",
					Members: []string{"user:user1"},
				},
			},
		}
		actual, err := server.ReadWritePolicy(ctx, "resources/test1", func(policy *iampb.Policy) (*iampb.Policy, error) {
			iampolicy.AddBinding(policy, "roles/test.admin", "user:user1")
			return policy, nil
		})
		assert.NilError(t, err)
		assert.DeepEqual(t, expected.GetBindings(), actual.GetBindings(), protocmp.Transform())
		expected2 := &iampb.Policy{
			Bindings: []*iampb.Binding{
				{
					Role:    "roles/test.admin",
					Members: []string{"user:user1"},
				},
				{
					Role:    "roles/test.user",
					Members: []string{"user:user2"},
				},
			},
		}
		actual2, err := server.ReadWritePolicy(ctx, "resources/test1", func(policy *iampb.Policy) (*iampb.Policy, error) {
			assert.DeepEqual(t, actual, policy, protocmp.Transform())
			iampolicy.AddBinding(policy, "roles/test.user", "user:user2")
			return policy, nil
		})
		assert.NilError(t, err)
		assert.DeepEqual(t, expected2.GetBindings(), actual2.GetBindings(), protocmp.Transform())
	})
}

func withMembers(ctx context.Context, members ...string) context.Context {
	return iamcaller.WithResolvedContext(ctx, &iamv1.Caller{
		Members: members,
		Metadata: map[string]*iamv1.Caller_Metadata{
			"test": {Members: members},
		},
	})
}
