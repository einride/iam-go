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

func (ts *serverTestSuite) testListSites(t *testing.T) {
	t.Parallel()
	ctx := withTestDeadline(context.Background(), t)

	t.Run("authorized", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			const (
				member = "user:test@example.com"
				parent = "shippers/aaaa"
				count  = 20
			)
			fx := ts.newTestFixture(t)
			fx.createShipper(t, parent)
			fx.iam.AddPolicyBinding(t, "/", "roles/freight.admin", member)
			expected := make([]*iamexamplev1.Site, 0, count)
			for i := 0; i < count; i++ {
				created, err := fx.client.CreateSite(
					WithOutgoingMembers(ctx, member),
					&iamexamplev1.CreateSiteRequest{
						Parent: parent,
						Site: &iamexamplev1.Site{
							DisplayName: fmt.Sprintf("Test Site %d", i),
						},
						SiteId: fmt.Sprintf("site%04d", i),
					},
				)
				assert.NilError(t, err)
				expected = append(expected, created)
			}
			actual := make([]*iamexamplev1.Site, 0, count)
			var pageToken string
			for {
				response, err := fx.client.ListSites(
					WithOutgoingMembers(ctx, member),
					&iamexamplev1.ListSitesRequest{
						Parent:    parent,
						PageSize:  count / 6,
						PageToken: pageToken,
					},
				)
				assert.NilError(t, err)
				actual = append(actual, response.Sites...)
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
		fx := ts.newTestFixture(t)
		response, err := fx.client.ListSites(
			WithOutgoingMembers(ctx, member),
			&iamexamplev1.ListSitesRequest{
				PageSize: 10,
			},
		)
		assert.Equal(t, codes.PermissionDenied, status.Code(err), "unexpected status: %v", err)
		assert.Assert(t, response == nil)
	})
}
