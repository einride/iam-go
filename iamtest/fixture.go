package iamtest

import (
	"context"
	"testing"

	"go.einride.tech/iam/iamspanner"
	"google.golang.org/genproto/googleapis/iam/v1"
	"gotest.tools/v3/assert"
)

// Fixture is a test fixture with helper methods for IAM testing.
type Fixture struct {
	server *iamspanner.Server
}

// NewFixture creates a new Fixture for the provided iamspanner.Server.
func NewFixture(server *iamspanner.Server) *Fixture {
	return &Fixture{server: server}
}

// AddPolicyBinding adds the provided policy binding.
func (fx *Fixture) AddPolicyBinding(t *testing.T, resource, role, member string) {
	// Get current policy.
	policy, err := fx.server.GetIamPolicy(context.TODO(), &iam.GetIamPolicyRequest{
		Resource: resource,
	})
	assert.NilError(t, err)
	// Add binding to policy.
	var added bool
	for _, binding := range policy.Bindings {
		if binding.Role == role {
			for _, bindingMember := range binding.Members {
				if bindingMember == member {
					return // already have this policy binding
				}
			}
			binding.Members = append(binding.Members, member)
			added = true
		}
	}
	if !added {
		policy.Bindings = append(policy.Bindings, &iam.Binding{
			Role:    role,
			Members: []string{member},
		})
	}
	// Set updated policy.
	_, err = fx.server.SetIamPolicy(context.TODO(), &iam.SetIamPolicyRequest{
		Resource: resource,
		Policy:   policy,
	})
	assert.NilError(t, err)
}
