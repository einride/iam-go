package iamexample

import (
	"context"
	"fmt"
	"testing"

	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"gotest.tools/v3/assert"
)

func (ts *serverTestSuite) testBatchGetSites(t *testing.T) {
	t.Parallel()
	ctx := withTestDeadline(context.Background(), t)

	t.Run("authorized", func(t *testing.T) {
		t.Run("permission on parent with requested parent", func(t *testing.T) {
			const (
				member = "user:test@example.com"
				parent = "shippers/1234"
				count  = 20
			)
			fx := ts.newTestFixture(t)
			fx.createShipper(t, parent)
			fx.iam.AddPolicyBinding(t, "*", "roles/freight.admin", member)
			expected := make([]*iamexamplev1.Site, 0, count)
			names := make([]string, 0, count)
			for i := 0; i < count; i++ {
				created, err := fx.client.CreateSite(
					WithOutgoingMembers(ctx, member),
					&iamexamplev1.CreateSiteRequest{
						Parent: parent,
						Site: &iamexamplev1.Site{
							DisplayName: fmt.Sprintf("Test Site %d", i),
						},
						SiteId: fmt.Sprintf("%04d", i),
					},
				)
				assert.NilError(t, err)
				expected = append(expected, created)
				names = append(names, created.Name)
			}
			response, err := fx.client.BatchGetSites(
				WithOutgoingMembers(ctx, member),
				&iamexamplev1.BatchGetSitesRequest{
					Parent: parent,
					Names:  names,
				},
			)
			assert.NilError(t, err)
			assert.DeepEqual(t, expected, response.Sites, protocmp.Transform())
		})

		t.Run("permission on parent without requested parent", func(t *testing.T) {
			const (
				member = "user:test@example.com"
				parent = "shippers/1234"
				count  = 20
			)
			fx := ts.newTestFixture(t)
			fx.createShipper(t, parent)
			fx.iam.AddPolicyBinding(t, "*", "roles/freight.admin", member)
			expected := make([]*iamexamplev1.Site, 0, count)
			names := make([]string, 0, count)
			for i := 0; i < count; i++ {
				created, err := fx.client.CreateSite(
					WithOutgoingMembers(ctx, member),
					&iamexamplev1.CreateSiteRequest{
						Parent: parent,
						Site: &iamexamplev1.Site{
							DisplayName: fmt.Sprintf("Test Site %d", i),
						},
						SiteId: fmt.Sprintf("%04d", i),
					},
				)
				assert.NilError(t, err)
				expected = append(expected, created)
				names = append(names, created.Name)
			}
			response, err := fx.client.BatchGetSites(
				WithOutgoingMembers(ctx, member),
				&iamexamplev1.BatchGetSitesRequest{
					Names: names,
				},
			)
			assert.NilError(t, err)
			assert.DeepEqual(t, expected, response.Sites, protocmp.Transform())
		})
	})

	t.Run("unauthorized", func(t *testing.T) {
		const member = "user:test@example.com"
		fx := ts.newTestFixture(t)
		response, err := fx.client.BatchGetSites(
			WithOutgoingMembers(ctx, member),
			&iamexamplev1.BatchGetSitesRequest{
				Names: []string{"shippers/1234/sites/5678"},
			},
		)
		assert.Equal(t, codes.PermissionDenied, status.Code(err), "unexpected status: %v", err)
		assert.Assert(t, response == nil)
	})
}
