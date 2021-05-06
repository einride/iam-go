package iamexample

import (
	"context"
	"testing"

	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gotest.tools/v3/assert"
)

func (ts *serverTestSuite) testGetSite(t *testing.T) {
	t.Parallel()
	ctx := withTestDeadline(context.Background(), t)

	t.Run("authorized", func(t *testing.T) {
		t.Run("not found", func(t *testing.T) {
			const (
				member = "user:test@example.com"
				site   = "shippers/1234/sites/5678"
			)
			fx := ts.newTestFixture(t)
			fx.iam.AddPolicyBinding(t, "/", "roles/freight.admin", member)
			got, err := fx.client.GetSite(
				WithOutgoingMembers(ctx, member),
				&iamexamplev1.GetSiteRequest{
					Name: site,
				},
			)
			assert.Equal(t, codes.NotFound, status.Code(err), "unexpected status: %v", err)
			assert.Assert(t, got == nil)
		})
	})

	t.Run("unauthorized", func(t *testing.T) {
		const (
			member = "user:test@example.com"
			site   = "shippers/1234/sites/5678"
		)
		fx := ts.newTestFixture(t)
		got, err := fx.client.GetSite(
			WithOutgoingMembers(ctx, member),
			&iamexamplev1.GetSiteRequest{
				Name: site,
			},
		)
		assert.Equal(t, codes.PermissionDenied, status.Code(err), "unexpected status: %v", err)
		assert.Assert(t, got == nil)
	})
}
