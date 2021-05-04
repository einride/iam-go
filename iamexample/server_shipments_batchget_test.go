package iamexample

import (
	"context"
	"fmt"
	"testing"

	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gotest.tools/v3/assert"
)

func (ts *serverTestSuite) testBatchGetShipments(t *testing.T) {
	t.Parallel()
	ctx := withTestDeadline(context.Background(), t)

	t.Run("authorized", func(t *testing.T) {
		t.Run("permission on parent with requested parent", func(t *testing.T) {
			const (
				member          = "user:test@example.com"
				parent          = "shippers/1234"
				originSite      = "shippers/1234/sites/origin"
				destinationSite = "shippers/1234/sites/destination"
				count           = 20
			)
			fx := ts.newTestFixture(t)
			fx.createShipper(t, parent)
			fx.createSite(t, originSite)
			fx.createSite(t, destinationSite)
			fx.iam.AddPolicyBinding(t, "*", "roles/freight.admin", member)
			ctx := WithOutgoingMembers(ctx, member)
			expected := make([]*iamexamplev1.Shipment, 0, count)
			names := make([]string, 0, count)
			for i := 0; i < count; i++ {
				created, err := fx.client.CreateShipment(ctx, &iamexamplev1.CreateShipmentRequest{
					Parent: parent,
					Shipment: &iamexamplev1.Shipment{
						OriginSite:           originSite,
						DestinationSite:      destinationSite,
						PickupEarliestTime:   timestamppb.Now(),
						PickupLatestTime:     timestamppb.Now(),
						DeliveryEarliestTime: timestamppb.Now(),
						DeliveryLatestTime:   timestamppb.Now(),
						LineItems: []*iamexamplev1.LineItem{
							{Title: "test 1", Quantity: 1},
							{Title: "test 2", Quantity: 2},
						},
					},
					ShipmentId: fmt.Sprintf("%04d", i),
				})
				assert.NilError(t, err)
				expected = append(expected, created)
				names = append(names, created.Name)
			}
			response, err := fx.client.BatchGetShipments(ctx, &iamexamplev1.BatchGetShipmentsRequest{
				Parent: parent,
				Names:  names,
			})
			assert.NilError(t, err)
			assert.DeepEqual(t, expected, response.Shipments, protocmp.Transform())
		})

		t.Run("permission on parent without requested parent", func(t *testing.T) {
			const (
				member          = "user:test@example.com"
				parent          = "shippers/1234"
				originSite      = "shippers/1234/sites/origin"
				destinationSite = "shippers/1234/sites/destination"
				count           = 20
			)
			fx := ts.newTestFixture(t)
			fx.iam.AddPolicyBinding(t, "*", "roles/freight.admin", member)
			fx.createShipper(t, parent)
			fx.createSite(t, originSite)
			fx.createSite(t, destinationSite)
			expected := make([]*iamexamplev1.Shipment, 0, count)
			names := make([]string, 0, count)
			for i := 0; i < count; i++ {
				created, err := fx.client.CreateShipment(
					WithOutgoingMembers(ctx, member),
					&iamexamplev1.CreateShipmentRequest{
						Parent: parent,
						Shipment: &iamexamplev1.Shipment{
							OriginSite:           originSite,
							DestinationSite:      destinationSite,
							PickupEarliestTime:   timestamppb.Now(),
							PickupLatestTime:     timestamppb.Now(),
							DeliveryEarliestTime: timestamppb.Now(),
							DeliveryLatestTime:   timestamppb.Now(),
							LineItems: []*iamexamplev1.LineItem{
								{Title: "test 1", Quantity: 1},
								{Title: "test 2", Quantity: 2},
							},
						},
						ShipmentId: fmt.Sprintf("%04d", i),
					},
				)
				assert.NilError(t, err)
				expected = append(expected, created)
				names = append(names, created.Name)
			}
			response, err := fx.client.BatchGetShipments(
				WithOutgoingMembers(ctx, member),
				&iamexamplev1.BatchGetShipmentsRequest{
					Names: names,
				},
			)
			assert.NilError(t, err)
			assert.DeepEqual(t, expected, response.Shipments, protocmp.Transform())
		})
	})

	t.Run("unauthorized", func(t *testing.T) {
		const (
			admin           = "user:admin@example.com"
			user            = "user:user@example.com"
			parent          = "shippers/1234"
			originSite      = "shippers/1234/sites/origin"
			destinationSite = "shippers/1234/sites/destination"
			shipmentID      = "test"
		)
		fx := ts.newTestFixture(t)
		fx.createShipper(t, parent)
		fx.createSite(t, originSite)
		fx.createSite(t, destinationSite)
		fx.iam.AddPolicyBinding(t, "*", "roles/freight.admin", admin)
		created, err := fx.client.CreateShipment(
			WithOutgoingMembers(ctx, admin),
			&iamexamplev1.CreateShipmentRequest{
				Parent: parent,
				Shipment: &iamexamplev1.Shipment{
					OriginSite:           originSite,
					DestinationSite:      destinationSite,
					PickupEarliestTime:   timestamppb.Now(),
					PickupLatestTime:     timestamppb.Now(),
					DeliveryEarliestTime: timestamppb.Now(),
					DeliveryLatestTime:   timestamppb.Now(),
					LineItems: []*iamexamplev1.LineItem{
						{Title: "test 1", Quantity: 1},
						{Title: "test 2", Quantity: 2},
					},
				},
				ShipmentId: shipmentID,
			},
		)
		assert.NilError(t, err)
		response, err := fx.client.BatchGetShipments(
			WithOutgoingMembers(ctx, user),
			&iamexamplev1.BatchGetShipmentsRequest{
				Names: []string{created.Name},
			},
		)
		assert.Equal(t, codes.PermissionDenied, status.Code(err), "unexpected status: %v", err)
		assert.Assert(t, response == nil)
	})
}
