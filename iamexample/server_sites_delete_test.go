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

func (ts *serverTestSuite) testDeleteSite(t *testing.T) {
	t.Parallel()
	ctx := withTestDeadline(context.Background(), t)

	t.Run("authorized", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			const (
				member = "user:test@example.com"
				parent = "shippers/aaaa"
				siteID = "bbbb"
			)
			fx := ts.newTestFixture(t)
			fx.createShipper(t, parent)
			fx.iam.AddPolicyBinding(t, "/", "roles/freight.admin", member)
			input := &iamexamplev1.Site{
				DisplayName: "Test Site",
			}
			created, err := fx.client.CreateSite(
				WithOutgoingMembers(ctx, member),
				&iamexamplev1.CreateSiteRequest{
					Parent: parent,
					Site:   input,
					SiteId: siteID,
				},
			)
			assert.NilError(t, err)
			assert.Equal(t, input.GetDisplayName(), created.GetDisplayName())
			deleted, err := fx.client.DeleteSite(
				WithOutgoingMembers(ctx, member),
				&iamexamplev1.DeleteSiteRequest{
					Name: created.GetName(),
				},
			)
			assert.NilError(t, err)
			assert.Equal(t, created.GetName(), deleted.GetName())
			assert.Assert(t, time.Since(deleted.GetDeleteTime().AsTime()) < time.Second)
		})
	})

	t.Run("unauthorized", func(t *testing.T) {
		const (
			member = "user:test@example.com"
			site   = "shippers/aaaa/sites/bbbb"
		)
		fx := ts.newTestFixture(t)
		deleted, err := fx.client.DeleteSite(
			WithOutgoingMembers(ctx, member),
			&iamexamplev1.DeleteSiteRequest{
				Name: site,
			},
		)
		assert.Equal(t, codes.PermissionDenied, status.Code(err), "unexpected status: %v", err)
		assert.Assert(t, deleted == nil)
	})
}
