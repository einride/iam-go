package iamexample

import (
	"context"
	"testing"

	"cloud.google.com/go/longrunning/autogen/longrunningpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gotest.tools/v3/assert"
)

func (ts *serverTestSuite) testLongRunningOperations(t *testing.T) {
	t.Parallel()
	ctx := withTestDeadline(context.Background(), t)

	t.Run("ListOperations", func(t *testing.T) {
		t.Run("authorized", func(t *testing.T) {
			const member = "user:test@example.com"
			fx := ts.newTestFixture(t)
			fx.iam.AddPolicyBinding(t, "shippers/1234", "roles/freight.admin", member)
			response, err := fx.longRunningClient.ListOperations(
				WithOutgoingMembers(ctx, member),
				&longrunningpb.ListOperationsRequest{Name: "shippers/1234"},
			)
			assert.NilError(t, err)
			assert.Assert(t, len(response.GetOperations()) == 0)
			assert.Equal(t, "", response.GetNextPageToken())
		})

		t.Run("unauthorized", func(t *testing.T) {
			const member = "user:test@example.com"
			fx := ts.newTestFixture(t)
			response, err := fx.longRunningClient.ListOperations(
				WithOutgoingMembers(ctx, member),
				&longrunningpb.ListOperationsRequest{Name: "shippers/1234"},
			)
			assert.Equal(t, codes.PermissionDenied, status.Code(err), status.Convert(err).Message())
			assert.Assert(t, response == nil)
		})
	})

	t.Run("GetOperation", func(t *testing.T) {
		t.Run("authorized", func(t *testing.T) {
			const member = "user:test@example.com"
			fx := ts.newTestFixture(t)
			fx.iam.AddPolicyBinding(t, "/", "roles/freight.admin", member)
			response, err := fx.longRunningClient.GetOperation(
				WithOutgoingMembers(ctx, member),
				&longrunningpb.GetOperationRequest{Name: "shippers/1234/operations/4567"},
			)
			assert.Equal(t, codes.NotFound, status.Code(err), status.Convert(err).Message())
			assert.Assert(t, response == nil)
		})

		t.Run("unauthorized", func(t *testing.T) {
			const member = "user:test@example.com"
			fx := ts.newTestFixture(t)
			response, err := fx.longRunningClient.GetOperation(
				WithOutgoingMembers(ctx, member),
				&longrunningpb.GetOperationRequest{Name: "shippers/1234/operations/4567"},
			)
			assert.Equal(t, codes.PermissionDenied, status.Code(err))
			assert.Assert(t, response == nil)
		})
	})

	t.Run("CancelOperation", func(t *testing.T) {
		t.Run("authorized", func(t *testing.T) {
			const member = "user:test@example.com"
			fx := ts.newTestFixture(t)
			fx.iam.AddPolicyBinding(t, "/", "roles/freight.admin", member)
			response, err := fx.longRunningClient.CancelOperation(
				WithOutgoingMembers(ctx, member),
				&longrunningpb.CancelOperationRequest{Name: "shippers/1234/operations/4567"},
			)
			assert.Equal(t, codes.NotFound, status.Code(err), status.Convert(err).Message())
			assert.Assert(t, response == nil)
		})

		t.Run("unauthorized", func(t *testing.T) {
			const member = "user:test@example.com"
			fx := ts.newTestFixture(t)
			response, err := fx.longRunningClient.CancelOperation(
				WithOutgoingMembers(ctx, member),
				&longrunningpb.CancelOperationRequest{Name: "shippers/1234/operations/4567"},
			)
			assert.Equal(t, codes.PermissionDenied, status.Code(err))
			assert.Assert(t, response == nil)
		})
	})

	t.Run("DeleteOperation", func(t *testing.T) {
		t.Run("authorized", func(t *testing.T) {
			const member = "user:test@example.com"
			fx := ts.newTestFixture(t)
			fx.iam.AddPolicyBinding(t, "/", "roles/freight.admin", member)
			response, err := fx.longRunningClient.DeleteOperation(
				WithOutgoingMembers(ctx, member),
				&longrunningpb.DeleteOperationRequest{Name: "shippers/1234/operations/4567"},
			)
			assert.Equal(t, codes.NotFound, status.Code(err), status.Convert(err).Message())
			assert.Assert(t, response == nil)
		})

		t.Run("unauthorized", func(t *testing.T) {
			const member = "user:test@example.com"
			fx := ts.newTestFixture(t)
			response, err := fx.longRunningClient.DeleteOperation(
				WithOutgoingMembers(ctx, member),
				&longrunningpb.DeleteOperationRequest{Name: "shippers/1234/operations/4567"},
			)
			assert.Equal(t, codes.PermissionDenied, status.Code(err))
			assert.Assert(t, response == nil)
		})
	})

	t.Run("WaitOperation", func(t *testing.T) {
		t.Run("authorized", func(t *testing.T) {
			const member = "user:test@example.com"
			fx := ts.newTestFixture(t)
			fx.iam.AddPolicyBinding(t, "/", "roles/freight.admin", member)
			response, err := fx.longRunningClient.WaitOperation(
				WithOutgoingMembers(ctx, member),
				&longrunningpb.WaitOperationRequest{Name: "shippers/1234/operations/4567"},
			)
			assert.Equal(t, codes.NotFound, status.Code(err), status.Convert(err).Message())
			assert.Assert(t, response == nil)
		})

		t.Run("unauthorized", func(t *testing.T) {
			const member = "user:test@example.com"
			fx := ts.newTestFixture(t)
			response, err := fx.longRunningClient.WaitOperation(
				WithOutgoingMembers(ctx, member),
				&longrunningpb.WaitOperationRequest{Name: "shippers/1234/operations/4567"},
			)
			assert.Equal(t, codes.PermissionDenied, status.Code(err))
			assert.Assert(t, response == nil)
		})
	})
}
