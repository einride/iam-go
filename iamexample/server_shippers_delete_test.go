package iamexample

import (
	"context"
	"testing"
	"time"

	"go.einride.tech/iam/iamspanner"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gotest.tools/v3/assert"
)

func testDeleteShippers(ctx context.Context, t *testing.T, newServer func(iamspanner.MemberResolver) iamexamplev1.FreightServiceServer) {
	t.Run("Delete", func(t *testing.T) {
		t.Run("authorized", func(t *testing.T) {
			t.Run("ok", func(t *testing.T) {
				const (
					member    = "user:test@example.com"
					shipperID = "123"
				)
				server := newServer(constantMember(member))
				addPolicyBinding(ctx, t, server, "*", "roles/freight.admin", member)
				input := &iamexamplev1.Shipper{
					DisplayName: "Test Shipper",
				}
				created, err := server.CreateShipper(ctx, &iamexamplev1.CreateShipperRequest{
					Shipper:   input,
					ShipperId: shipperID,
				})
				assert.NilError(t, err)
				assert.Equal(t, input.DisplayName, created.DisplayName)
				deleted, err := server.DeleteShipper(ctx, &iamexamplev1.DeleteShipperRequest{
					Name: created.Name,
				})
				assert.NilError(t, err)
				assert.Equal(t, created.Name, deleted.Name)
				assert.Assert(t, time.Since(deleted.DeleteTime.AsTime()) < time.Second)
			})
		})

		t.Run("unauthorized", func(t *testing.T) {
			const (
				member    = "user:test@example.com"
				shipperID = "123"
				shipper   = "shippers/" + shipperID
			)
			server := newServer(constantMember(member))
			deleted, err := server.DeleteShipper(ctx, &iamexamplev1.DeleteShipperRequest{
				Name: shipper,
			})
			assert.Equal(t, codes.PermissionDenied, status.Code(err), "unexpected status: %v", err)
			assert.Assert(t, deleted == nil)
		})
	})
}
