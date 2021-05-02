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

func testUpdateSite(
	ctx context.Context,
	t *testing.T,
	newServer func(iamspanner.MemberResolver) iamexamplev1.FreightServiceServer,
) {
	t.Run("UpdateSite", func(t *testing.T) {
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
				got, err := server.CreateSite(ctx, &iamexamplev1.CreateSiteRequest{
					Parent: parent,
					Site:   input,
					SiteId: siteID,
				})
				assert.NilError(t, err)
				assert.Equal(t, input.DisplayName, got.DisplayName)
				update := &iamexamplev1.Site{
					Name:        got.Name,
					DisplayName: "Updated Test Site",
				}
				updated, err := server.UpdateSite(ctx, &iamexamplev1.UpdateSiteRequest{
					Site: update,
				})
				assert.NilError(t, err)
				assert.Equal(t, update.DisplayName, updated.DisplayName)
			})
		})

		t.Run("unauthorized", func(t *testing.T) {
			const (
				member = "user:test@example.com"
				parent = "shippers/1234"
				siteID = "5678"
			)
			server := newServer(constantMember(member))
			got, err := server.CreateSite(ctx, &iamexamplev1.CreateSiteRequest{
				Parent: parent,
				Site: &iamexamplev1.Site{
					DisplayName: "Test Site",
				},
				SiteId: siteID,
			})
			assert.Equal(t, codes.PermissionDenied, status.Code(err), "unexpected status: %v", err)
			assert.Assert(t, got == nil)
		})
	})
}
