package iampolicy

import "google.golang.org/genproto/googleapis/iam/v1"

// AddBinding adds the provided role and member binding to the policy.
// If the role and member already exists, no updates are made.
// No validation on the role or member is performed.
func AddBinding(policy *iam.Policy, role, member string) {
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
}
