package iamexample

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/protobuf/types/known/timestamppb"

	"go.einride.tech/iam/iamspanner"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"gotest.tools/v3/assert"
)

func testBatchGetShipments(
	ctx context.Context,
	t *testing.T,
	newServer func(iamspanner.MemberResolver) iamexamplev1.FreightServiceServer,
) {
	t.Run("BatchGetShipments", func(t *testing.T) {
		t.Run("authorized", func(t *testing.T) {
			t.Run("permission on parent with requested parent", func(t *testing.T) {
				const (
					member          = "user:test@example.com"
					parent          = "shippers/1234"
					originSite      = "shippers/1234/sites/origin"
					destinationSite = "shippers/1234/sites/destination"
					count           = 20
				)
				server := newServer(constantMember(member))
				addPolicyBinding(ctx, t, server, "*", "roles/freight.admin", member)
				createShipper(ctx, t, server, parent)
				createSite(ctx, t, server, originSite)
				createSite(ctx, t, server, destinationSite)
				expected := make([]*iamexamplev1.Shipment, 0, count)
				names := make([]string, 0, count)
				for i := 0; i < count; i++ {
					created, err := server.CreateShipment(ctx, &iamexamplev1.CreateShipmentRequest{
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
				response, err := server.BatchGetShipments(ctx, &iamexamplev1.BatchGetShipmentsRequest{
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
				server := newServer(constantMember(member))
				addPolicyBinding(ctx, t, server, "*", "roles/freight.admin", member)
				createShipper(ctx, t, server, parent)
				createSite(ctx, t, server, originSite)
				createSite(ctx, t, server, destinationSite)
				expected := make([]*iamexamplev1.Shipment, 0, count)
				names := make([]string, 0, count)
				for i := 0; i < count; i++ {
					created, err := server.CreateShipment(ctx, &iamexamplev1.CreateShipmentRequest{
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
				response, err := server.BatchGetShipments(ctx, &iamexamplev1.BatchGetShipmentsRequest{
					Names: names,
				})
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
			member := admin
			server := newServer(ptrMember(&member))
			addPolicyBinding(ctx, t, server, "*", "roles/freight.admin", admin)
			createShipper(ctx, t, server, parent)
			createSite(ctx, t, server, originSite)
			createSite(ctx, t, server, destinationSite)
			created, err := server.CreateShipment(ctx, &iamexamplev1.CreateShipmentRequest{
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
			})
			assert.NilError(t, err)
			member = user
			response, err := server.BatchGetShipments(ctx, &iamexamplev1.BatchGetShipmentsRequest{
				Names: []string{created.Name},
			})
			assert.Equal(t, codes.PermissionDenied, status.Code(err), "unexpected status: %v", err)
			assert.Assert(t, response == nil)
		})
	})
}
