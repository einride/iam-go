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

func testDeleteSite(
	ctx context.Context,
	t *testing.T,
	newServer func(iamspanner.MemberResolver) iamexamplev1.FreightServiceServer,
) {
	t.Run("DeleteSite", func(t *testing.T) {
		t.Run("authorized", func(t *testing.T) {
			t.Run("ok", func(t *testing.T) {
				const (
					member = "user:test@example.com"
					parent = "shippers/1234"
					siteID = "5678"
				)
				server := newServer(constantMember(member))
				addPolicyBinding(ctx, t, server, "*", "roles/freight.admin", member)
				createShipper(ctx, t, server, parent)
				input := &iamexamplev1.Site{
					DisplayName: "Test Site",
				}
				created, err := server.CreateSite(ctx, &iamexamplev1.CreateSiteRequest{
					Parent: parent,
					Site:   input,
					SiteId: siteID,
				})
				assert.NilError(t, err)
				assert.Equal(t, input.DisplayName, created.DisplayName)
				deleted, err := server.DeleteSite(ctx, &iamexamplev1.DeleteSiteRequest{
					Name: created.Name,
				})
				assert.NilError(t, err)
				assert.Equal(t, created.Name, deleted.Name)
				assert.Assert(t, time.Since(deleted.DeleteTime.AsTime()) < time.Second)
			})
		})

		t.Run("unauthorized", func(t *testing.T) {
			const (
				member = "user:test@example.com"
				site   = "shippers/1234/sites/5678"
			)
			server := newServer(constantMember(member))
			deleted, err := server.DeleteSite(ctx, &iamexamplev1.DeleteSiteRequest{
				Name: site,
			})
			assert.Equal(t, codes.PermissionDenied, status.Code(err), "unexpected status: %v", err)
			assert.Assert(t, deleted == nil)
		})
	})
}
