package iamexample

import (
	"context"
	"testing"

	"go.einride.tech/iam/iamregistry"
	"go.einride.tech/iam/iamspanner"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"go.einride.tech/spanner-aip/spantest"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"gotest.tools/v3/assert"
)

func TestServer(t *testing.T) {
	fx := spantest.NewEmulatorFixture(t)
	ctx := context.Background()
	if deadline, ok := t.Deadline(); ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, deadline)
		t.Cleanup(cancel)
	}
	newServer := func(resolver iamspanner.MemberResolver) iamexamplev1.FreightServiceServer {
		spannerClient := fx.NewDatabaseFromDDLFiles(t, "schema.sql", "../iamspanner/schema.sql")
		iamServer, err := iamspanner.NewServer(
			spannerClient,
			lookupPredefinedRoles(t),
			resolver,
			iamspanner.ServerConfig{
				ErrorHook: func(ctx context.Context, err error) {
					t.Log(err)
				},
			},
		)
		assert.NilError(t, err)
		server := &Server{
			IAM:     iamServer,
			Spanner: spannerClient,
			Config: Config{
				ErrorHook: func(ctx context.Context, err error) {
					t.Log(err)
				},
			},
		}
		return &Authorization{
			Next: server,
		}
	}
	testGet(ctx, t, newServer)
	testCreate(ctx, t, newServer)
}

func addPolicyBinding(
	ctx context.Context,
	t *testing.T,
	server iam.IAMPolicyServer,
	resource string,
	role string,
	member string,
) {
	// Bypass authorization.
	authorization, ok := server.(*Authorization)
	assert.Assert(t, ok)
	serverImpl, ok := authorization.Next.(*Server)
	assert.Assert(t, ok)
	// Get current policy.
	policy, err := serverImpl.IAM.GetIamPolicy(ctx, &iam.GetIamPolicyRequest{
		Resource: "resource",
	})
	assert.NilError(t, err)
	// Add binding to policy.
	var added bool
	for _, binding := range policy.Bindings {
		if binding.Role == role {
			for _, bindingMember := range binding.Members {
				if bindingMember == member {
					return // already have this policy binding
				}
			}
			binding.Members = append(binding.Members, member)
			added = true
		}
	}
	if !added {
		policy.Bindings = append(policy.Bindings, &iam.Binding{
			Role:    role,
			Members: []string{member},
		})
	}
	// Set updated policy.
	_, err = serverImpl.IAM.SetIamPolicy(ctx, &iam.SetIamPolicyRequest{
		Resource: resource,
		Policy:   policy,
	})
	assert.NilError(t, err)
}

func lookupPredefinedRoles(t *testing.T) *iamregistry.Roles {
	desc, err := protoregistry.GlobalFiles.FindDescriptorByName(
		protoreflect.FullName(iamexamplev1.FreightService_ServiceDesc.ServiceName),
	)
	assert.NilError(t, err)
	serviceDesc, ok := desc.(protoreflect.ServiceDescriptor)
	assert.Assert(t, ok)
	predefinedRoles := proto.GetExtension(serviceDesc.Options(), iamv1.E_PredefinedRoles).(*iamv1.Roles)
	assert.Assert(t, predefinedRoles != nil)
	roles, err := iamregistry.NewRoles(predefinedRoles)
	assert.NilError(t, err)
	return roles
}

type memberResolverFn func(context.Context) (string, error)

func (f memberResolverFn) ResolveMember(ctx context.Context) (string, error) {
	return f(ctx)
}

func constantMember(member string) iamspanner.MemberResolver {
	return memberResolverFn(func(ctx context.Context) (string, error) {
		return member, nil
	})
}
