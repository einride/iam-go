package iamexample

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	"cloud.google.com/go/longrunning/autogen/longrunningpb"
	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iamauthz"
	"go.einride.tech/iam/iamcaller"
	"go.einride.tech/iam/iamspanner"
	"go.einride.tech/iam/iamtest"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"go.einride.tech/spanner-aip/spantest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gotest.tools/v3/assert"
)

func TestServer(t *testing.T) {
	t.Parallel()
	ts := newServerTestSuite(t)
	// shippers
	t.Run("CreateShipper", ts.testCreateShipper)
	t.Run("DeleteShipper", ts.testDeleteShipper)
	t.Run("GetShipper", ts.testGetShipper)
	t.Run("ListShippers", ts.testListShippers)
	t.Run("UpdateShipper", ts.testUpdateShipper)
	// sites
	t.Run("CreateSite", ts.testCreateSite)
	t.Run("DeleteSite", ts.testDeleteSite)
	t.Run("GetSite", ts.testGetSite)
	t.Run("ListSites", ts.testListSites)
	t.Run("UpdateSite", ts.testUpdateSite)
	t.Run("BatchGetSites", ts.testBatchGetSites)
	// shipments
	t.Run("CreateShipment", ts.testCreateShipment)
	t.Run("BatchGetShipments", ts.testBatchGetShipments)
	// long-running operations
	t.Run("LongRunningOperations", ts.testLongRunningOperations)
	// IAM
	t.Run("IAM", ts.testIAM)
}

type serverTestSuite struct {
	spanner spantest.Fixture
}

func newServerTestSuite(t *testing.T) *serverTestSuite {
	return &serverTestSuite{
		spanner: spantest.NewEmulatorFixture(t),
	}
}

func (ts *serverTestSuite) newTestFixture(t *testing.T) *serverTestFixture {
	iamDescriptor, err := iamexamplev1.NewFreightServiceIAMDescriptor()
	assert.NilError(t, err)
	spannerClient := ts.spanner.NewDatabaseFromDDLFiles(t, "schema.sql", "../iamspanner/schema.sql")
	iamServer, err := iamspanner.NewIAMServer(
		spannerClient,
		iamDescriptor.PredefinedRoles.GetRole(),
		iamcaller.FromContextResolver(),
		iamspanner.ServerConfig{
			ErrorHook: func(_ context.Context, err error) {
				t.Log(err)
			},
		},
	)
	assert.NilError(t, err)
	server := &Server{
		IAM:     iamServer,
		Spanner: spannerClient,
		Config: Config{
			ErrorHook: func(_ context.Context, err error) {
				t.Log(err)
			},
		},
	}
	authorization, err := iamexamplev1.NewFreightServiceAuthorization(
		server, iamServer, iamcaller.FromContextResolver(),
	)
	assert.NilError(t, err)
	serverWithAuthorization := &Authorization{
		Next:                        server,
		IAMServer:                   iamServer,
		IAMDescriptor:               iamDescriptor,
		FreightServiceAuthorization: authorization,
		CallerResolver:              iamcaller.FromContextResolver(),
	}
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			iamcaller.ResolveContextUnaryInterceptor(NewMemberHeaderResolver()),
			iamauthz.RequireAuthorizationUnaryInterceptor,
		),
	)
	iamexamplev1.RegisterFreightServiceServer(grpcServer, serverWithAuthorization)
	longrunningpb.RegisterOperationsServer(grpcServer, serverWithAuthorization)
	errChan := make(chan error)
	go func() {
		if err := grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			errChan <- err
			return
		}
		errChan <- nil
	}()
	t.Cleanup(func() {
		assert.NilError(t, <-errChan)
	})
	t.Cleanup(func() {
		grpcServer.GracefulStop()
	})
	ctx := withTestDeadline(context.Background(), t)
	conn, err := grpc.DialContext(
		ctx,
		lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	assert.NilError(t, err)
	serviceClient := iamexamplev1.NewFreightServiceClient(conn)
	longRunningClient := longrunningpb.NewOperationsClient(conn)
	return &serverTestFixture{
		iam:               iamtest.NewFixture(iamServer),
		spanner:           spannerClient,
		client:            serviceClient,
		longRunningClient: longRunningClient,
	}
}

type serverTestFixture struct {
	iam               *iamtest.Fixture
	spanner           *spanner.Client
	client            iamexamplev1.FreightServiceClient
	longRunningClient longrunningpb.OperationsClient
}

func (fx *serverTestFixture) createShipper(t *testing.T, name string) {
	t.Helper()
	const member = "user:fixture@example.com"
	fx.iam.AddPolicyBinding(t, "/", "roles/freight.admin", member)
	ctx := WithOutgoingMembers(withTestDeadline(context.Background(), t), member)
	var id string
	assert.NilError(t, resourcename.Sscan(name, "shippers/{shipper}", &id))
	input := &iamexamplev1.Shipper{
		DisplayName: fmt.Sprintf("shipper %s", id),
	}
	got, err := fx.client.CreateShipper(ctx, &iamexamplev1.CreateShipperRequest{
		Shipper:   input,
		ShipperId: id,
	})
	assert.NilError(t, err)
	assert.Equal(t, input.GetDisplayName(), got.GetDisplayName())
}

func (fx *serverTestFixture) createSite(t *testing.T, name string) {
	t.Helper()
	const member = "user:fixture@example.com"
	fx.iam.AddPolicyBinding(t, "/", "roles/freight.admin", member)
	ctx := WithOutgoingMembers(withTestDeadline(context.Background(), t), member)
	var shipperID, siteID string
	assert.NilError(t, resourcename.Sscan(name, "shippers/{shipper}/sites/{site}", &shipperID, &siteID))
	input := &iamexamplev1.Site{
		DisplayName: fmt.Sprintf("site %s", siteID),
	}
	got, err := fx.client.CreateSite(ctx, &iamexamplev1.CreateSiteRequest{
		Parent: resourcename.Sprint("shippers/{shipper}", shipperID),
		Site:   input,
		SiteId: siteID,
	})
	assert.NilError(t, err)
	assert.Equal(t, input.GetDisplayName(), got.GetDisplayName())
}

func withTestDeadline(ctx context.Context, t *testing.T) context.Context {
	deadline, ok := t.Deadline()
	if !ok {
		return ctx
	}
	ctx, cancel := context.WithDeadline(ctx, deadline)
	t.Cleanup(cancel)
	return ctx
}
