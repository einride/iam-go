package iamexample

import (
	"context"
	"testing"

	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gotest.tools/v3/assert"
)

func (ts *serverTestSuite) testCreateShipment(t *testing.T) {
	t.Parallel()
	ctx := withTestDeadline(context.Background(), t)

	t.Run("authorized", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			const (
				member          = "user:test@example.com"
				parent          = "shippers/aaaa"
				originSite      = "shippers/aaaa/sites/origin"
				destinationSite = "shippers/aaaa/sites/destination"
				shipmentID      = "bbbb"
			)
			fx := ts.newTestFixture(t)
			fx.createShipper(t, parent)
			fx.createSite(t, originSite)
			fx.createSite(t, destinationSite)
			fx.iam.AddPolicyBinding(t, "/", "roles/freight.admin", member)
			input := &iamexamplev1.Shipment{
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
			}
			got, err := fx.client.CreateShipment(
				WithOutgoingMembers(ctx, member),
				&iamexamplev1.CreateShipmentRequest{
					Parent:     parent,
					Shipment:   input,
					ShipmentId: shipmentID,
				},
			)
			assert.NilError(t, err)
			assert.DeepEqual(t, input.LineItems, got.LineItems, protocmp.Transform())
		})
	})

	t.Run("unauthorized", func(t *testing.T) {
		const (
			member          = "user:test@example.com"
			parent          = "shippers/aaaa"
			originSite      = "shippers/aaaa/sites/origin"
			destinationSite = "shippers/aaaa/sites/destination"
			shipmentID      = "5678"
		)
		fx := ts.newTestFixture(t)
		input := &iamexamplev1.Shipment{
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
		}
		got, err := fx.client.CreateShipment(
			WithOutgoingMembers(ctx, member),
			&iamexamplev1.CreateShipmentRequest{
				Parent:     parent,
				Shipment:   input,
				ShipmentId: shipmentID,
			},
		)
		assert.Equal(t, codes.PermissionDenied, status.Code(err), "unexpected status: %v", err)
		assert.Assert(t, got == nil)
	})
}
