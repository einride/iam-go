package iamexample

import (
	"context"
	"fmt"
	"testing"

	"go.einride.tech/iam/iamspanner"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"gotest.tools/v3/assert"
)

func testListShippers(
	ctx context.Context,
	t *testing.T,
	newServer func(iamspanner.MemberResolver,
	) iamexamplev1.FreightServiceServer) {
	t.Run("List", func(t *testing.T) {
		t.Run("authorized", func(t *testing.T) {
			t.Run("ok", func(t *testing.T) {
				const member = "user:test@example.com"
				server := newServer(constantMember(member))
				addPolicyBinding(ctx, t, server, "*", "roles/freight.admin", member)
				const count = 20
				expected := make([]*iamexamplev1.Shipper, 0, count)
				for i := 0; i < count; i++ {
					created, err := server.CreateShipper(ctx, &iamexamplev1.CreateShipperRequest{
						Shipper: &iamexamplev1.Shipper{
							DisplayName: fmt.Sprintf("Test Shipper %d", i),
						},
						ShipperId: fmt.Sprintf("%03d", i),
					})
					assert.NilError(t, err)
					expected = append(expected, created)
				}
				actual := make([]*iamexamplev1.Shipper, 0, count)
				var pageToken string
				for {
					response, err := server.ListShippers(ctx, &iamexamplev1.ListShippersRequest{
						PageSize:  count / 6,
						PageToken: pageToken,
					})
					assert.NilError(t, err)
					actual = append(actual, response.Shippers...)
					pageToken = response.NextPageToken
					if pageToken == "" {
						break
					}
				}
				assert.DeepEqual(t, expected, actual, protocmp.Transform())
			})
		})

		t.Run("unauthorized", func(t *testing.T) {
			const member = "user:test@example.com"
			server := newServer(constantMember(member))
			response, err := server.ListShippers(ctx, &iamexamplev1.ListShippersRequest{
				PageSize: 10,
			})
			assert.Equal(t, codes.PermissionDenied, status.Code(err), "unexpected status: %v", err)
			assert.Assert(t, response == nil)
		})
	})
}
