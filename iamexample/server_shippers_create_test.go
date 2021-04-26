package iamexample

import (
	"context"
	"testing"

	"go.einride.tech/iam/iamspanner"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gotest.tools/v3/assert"
)

func testCreate(ctx context.Context, t *testing.T, newServer func(iamspanner.MemberResolver) iamexamplev1.FreightServiceServer) {
	t.Run("Create", func(t *testing.T) {
		t.Run("authorized", func(t *testing.T) {
			t.Run("ok", func(t *testing.T) {
				const (
					member    = "user:test@example.com"
					shipperID = "123"
					shipper   = "shippers/" + shipperID
				)
				server := newServer(constantMember(member))
				addPolicyBinding(ctx, t, server, "*", "roles/freight.admin", member)
				input := &iamexamplev1.Shipper{
					DisplayName: "Test Shipper",
				}
				got, err := server.CreateShipper(ctx, &iamexamplev1.CreateShipperRequest{
					Shipper:   input,
					ShipperId: shipperID,
				})
				assert.NilError(t, err)
				assert.Equal(t, input.DisplayName, got.DisplayName)
			})
		})

		t.Run("unauthorized", func(t *testing.T) {
			const (
				member    = "user:test@example.com"
				shipperID = "123"
			)
			server := newServer(constantMember(member))
			got, err := server.CreateShipper(ctx, &iamexamplev1.CreateShipperRequest{
				Shipper: &iamexamplev1.Shipper{
					DisplayName: "Test Shipper",
				},
				ShipperId: shipperID,
			})
			assert.Equal(t, codes.PermissionDenied, status.Code(err), "unexpected status: %v", err)
			assert.Assert(t, got == nil)
		})
	})
}
