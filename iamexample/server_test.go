package iamexample

import (
	"context"
	"fmt"
	"testing"

	"go.einride.tech/aip/resourcename"
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
			IAM:  iamServer,
			Next: server,
		}
	}
	// shippers
	testGetShipper(ctx, t, newServer)
	testCreateShipper(ctx, t, newServer)
	testUpdateShipper(ctx, t, newServer)
	testListShippers(ctx, t, newServer)
	testDeleteShipper(ctx, t, newServer)
	// sites
	testCreateSite(ctx, t, newServer)
	testGetSite(ctx, t, newServer)
	testDeleteSite(ctx, t, newServer)
	testListSites(ctx, t, newServer)
	testUpdateSite(ctx, t, newServer)
	testBatchGetSites(ctx, t, newServer)
	// shipments
	testCreateShipment(ctx, t, newServer)
	testBatchGetShipments(ctx, t, newServer)
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

func createShipper(
	ctx context.Context,
	t *testing.T,
	server iamexamplev1.FreightServiceServer,
	name string,
) {
	t.Helper()
	var id string
	assert.NilError(t, resourcename.Sscan(name, "shippers/{shipper}", &id))
	input := &iamexamplev1.Shipper{
		DisplayName: fmt.Sprintf("shipper %s", id),
	}
	got, err := server.CreateShipper(ctx, &iamexamplev1.CreateShipperRequest{
		Shipper:   input,
		ShipperId: id,
	})
	assert.NilError(t, err)
	assert.Equal(t, input.DisplayName, got.DisplayName)
}

func createSite(
	ctx context.Context,
	t *testing.T,
	server iamexamplev1.FreightServiceServer,
	name string,
) {
	t.Helper()
	var shipperID, siteID string
	assert.NilError(t, resourcename.Sscan(name, "shippers/{shipper}/sites/{site}", &shipperID, &siteID))
	input := &iamexamplev1.Site{
		DisplayName: fmt.Sprintf("site %s", siteID),
	}
	got, err := server.CreateSite(ctx, &iamexamplev1.CreateSiteRequest{
		Parent: resourcename.Sprint("shippers/{shipper}", shipperID),
		Site:   input,
		SiteId: siteID,
	})
	assert.NilError(t, err)
	assert.Equal(t, input.DisplayName, got.DisplayName)
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

func ptrMember(member *string) iamspanner.MemberResolver {
	return memberResolverFn(func(ctx context.Context) (string, error) {
		return *member, nil
	})
}
