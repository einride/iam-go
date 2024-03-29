package iamtest

import (
	"context"
	"testing"

	"cloud.google.com/go/iam/apiv1/iampb"
	"go.einride.tech/iam/iampolicy"
	"go.einride.tech/iam/iamspanner"
	"gotest.tools/v3/assert"
)

// Fixture is a test fixture with helper methods for IAM testing.
type Fixture struct {
	server *iamspanner.IAMServer
}

// NewFixture creates a new Fixture for the provided iamspanner.IAMServer.
func NewFixture(server *iamspanner.IAMServer) *Fixture {
	return &Fixture{server: server}
}

// AddPolicyBinding adds the provided policy binding.
func (fx *Fixture) AddPolicyBinding(t *testing.T, resource, role, member string) {
	ctx := withTestDeadline(context.Background(), t)
	// Get current policy.
	policy, err := fx.server.GetIamPolicy(ctx, &iampb.GetIamPolicyRequest{
		Resource: resource,
	})
	assert.NilError(t, err)
	iampolicy.AddBinding(policy, role, member)
	// Set updated policy.
	_, err = fx.server.SetIamPolicy(ctx, &iampb.SetIamPolicyRequest{
		Resource: resource,
		Policy:   policy,
	})
	assert.NilError(t, err)
}

func withTestDeadline(ctx context.Context, t *testing.T) context.Context {
	deadline, ok := t.Deadline()
	if !ok {
		return ctx
	}
	ctx, cancel := context.WithDeadline(ctx, deadline)
	t.Cleanup(cancel)
	return ctx
}
