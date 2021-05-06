package iamexample

import (
	"context"
	"testing"
	"time"

	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gotest.tools/v3/assert"
)

func (ts *serverTestSuite) testDeleteShipper(t *testing.T) {
	t.Parallel()
	ctx := withTestDeadline(context.Background(), t)

	t.Run("authorized", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			const (
				member    = "user:test@example.com"
				shipperID = "1234"
			)
			fx := ts.newTestFixture(t)
			fx.iam.AddPolicyBinding(t, "/", "roles/freight.admin", member)
			input := &iamexamplev1.Shipper{
				DisplayName: "Test Shipper",
			}
			created, err := fx.client.CreateShipper(
				WithOutgoingMembers(ctx, member),
				&iamexamplev1.CreateShipperRequest{
					Shipper:   input,
					ShipperId: shipperID,
				},
			)
			assert.NilError(t, err)
			assert.Equal(t, input.DisplayName, created.DisplayName)
			deleted, err := fx.client.DeleteShipper(
				WithOutgoingMembers(ctx, member),
				&iamexamplev1.DeleteShipperRequest{
					Name: created.Name,
				},
			)
			assert.NilError(t, err)
			assert.Equal(t, created.Name, deleted.Name)
			assert.Assert(t, time.Since(deleted.DeleteTime.AsTime()) < time.Second)
		})
	})

	t.Run("unauthorized", func(t *testing.T) {
		const (
			member    = "user:test@example.com"
			shipperID = "1234"
			shipper   = "shippers/" + shipperID
		)
		fx := ts.newTestFixture(t)
		deleted, err := fx.client.DeleteShipper(
			WithOutgoingMembers(ctx, member),
			&iamexamplev1.DeleteShipperRequest{
				Name: shipper,
			},
		)
		assert.Equal(t, codes.PermissionDenied, status.Code(err), "unexpected status: %v", err)
		assert.Assert(t, deleted == nil)
	})
}
