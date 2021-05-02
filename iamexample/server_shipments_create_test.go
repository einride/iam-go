package iamexample

import (
	"context"
	"testing"

	"go.einride.tech/iam/iamspanner"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gotest.tools/v3/assert"
)

func testCreateShipment(
	ctx context.Context,
	t *testing.T,
	newServer func(iamspanner.MemberResolver) iamexamplev1.FreightServiceServer,
) {
	t.Run("CreateShipment", func(t *testing.T) {
		t.Run("authorized", func(t *testing.T) {
			t.Run("ok", func(t *testing.T) {
				const (
					member          = "user:test@example.com"
					parent          = "shippers/1234"
					originSite      = "shippers/1234/sites/origin"
					destinationSite = "shippers/1234/sites/destination"
					shipmentID      = "5678"
				)
				server := newServer(constantMember(member))
				addPolicyBinding(ctx, t, server, "*", "roles/freight.admin", member)
				createShipper(ctx, t, server, parent)
				createSite(ctx, t, server, originSite)
				createSite(ctx, t, server, destinationSite)
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
				got, err := server.CreateShipment(ctx, &iamexamplev1.CreateShipmentRequest{
					Parent:     parent,
					Shipment:   input,
					ShipmentId: shipmentID,
				})
				assert.NilError(t, err)
				assert.DeepEqual(t, input.LineItems, got.LineItems, protocmp.Transform())
			})
		})

		t.Run("unauthorized", func(t *testing.T) {
			const (
				member          = "user:test@example.com"
				parent          = "shippers/1234"
				originSite      = "shippers/1234/sites/origin"
				destinationSite = "shippers/1234/sites/destination"
				shipmentID      = "5678"
			)
			server := newServer(constantMember(member))
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
			got, err := server.CreateShipment(ctx, &iamexamplev1.CreateShipmentRequest{
				Parent:     parent,
				Shipment:   input,
				ShipmentId: shipmentID,
			})
			assert.Equal(t, codes.PermissionDenied, status.Code(err), "unexpected status: %v", err)
			assert.Assert(t, got == nil)
		})
	})
}
