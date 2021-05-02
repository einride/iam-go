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

func testGetSite(ctx context.Context, t *testing.T, newServer func(iamspanner.MemberResolver) iamexamplev1.FreightServiceServer) {
	t.Run("GetSite", func(t *testing.T) {
		t.Run("authorized", func(t *testing.T) {
			t.Run("not found", func(t *testing.T) {
				const (
					member = "user:test@example.com"
					site   = "shippers/1234/sites/5678"
				)
				server := newServer(constantMember(member))
				addPolicyBinding(ctx, t, server, "*", "roles/freight.admin", member)
				got, err := server.GetSite(ctx, &iamexamplev1.GetSiteRequest{
					Name: site,
				})
				assert.Equal(t, codes.NotFound, status.Code(err), "unexpected status: %v", err)
				assert.Assert(t, got == nil)
			})
		})

		t.Run("unauthorized", func(t *testing.T) {
			const (
				member = "user:test@example.com"
				site   = "shippers/1234/sites/5678"
			)
			server := newServer(constantMember(member))
			got, err := server.GetSite(ctx, &iamexamplev1.GetSiteRequest{
				Name: site,
			})
			assert.Equal(t, codes.PermissionDenied, status.Code(err), "unexpected status: %v", err)
			assert.Assert(t, got == nil)
		})
	})
}
