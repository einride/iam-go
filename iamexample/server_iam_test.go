package iamexample

import (
	"context"
	"testing"

	"cloud.google.com/go/iam/apiv1/iampb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gotest.tools/v3/assert"
)

func (ts *serverTestSuite) testIAM(t *testing.T) {
	t.Parallel()
	ctx := withTestDeadline(context.Background(), t)

	t.Run("authorized", func(t *testing.T) {
		t.Run("unconfigured resource", func(t *testing.T) {
			const (
				member  = "user:test@example.com"
				shipper = "shippers/1234"
			)
			fx := ts.newTestFixture(t)
			fx.iam.AddPolicyBinding(t, "/", "roles/freight.admin", member)
			got, err := fx.client.GetIamPolicy(
				WithOutgoingMembers(ctx, member),
				&iampb.GetIamPolicyRequest{
					Resource: "resources/foo",
				},
			)
			assert.Equal(t, codes.PermissionDenied, status.Code(err), "unexpected status: %v", err)
			assert.Assert(t, got == nil)
		})
	})
}
