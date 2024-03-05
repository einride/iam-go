package iampolicy

import (
	"cloud.google.com/go/iam/apiv1/iampb"
)

// AddBinding adds the provided role and member binding to the policy.
// If the role and member already exists, no updates are made.
// No validation on the role or member is performed.
func AddBinding(policy *iampb.Policy, role, member string) {
	// Add binding to policy.
	var added bool
	for _, binding := range policy.GetBindings() {
		if binding.GetRole() == role {
			for _, bindingMember := range binding.GetMembers() {
				if bindingMember == member {
					return // already have this policy binding
				}
			}
			binding.Members = append(binding.Members, member)
			added = true
		}
	}
	if !added {
		policy.Bindings = append(policy.Bindings, &iampb.Binding{
			Role:    role,
			Members: []string{member},
		})
	}
}
