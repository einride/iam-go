package iampolicy

import (
	"cloud.google.com/go/iam/apiv1/iampb"
)

// RemoveBinding removes the provided role and member binding from the policy.
// If a binding of the role and member don't exist, no updates are made.
// No validation on the role or member is performed.
func RemoveBinding(policy *iampb.Policy, role, member string) {
	for _, binding := range policy.Bindings {
		if binding.Role == role {
			binding.Members = removeMember(binding.Members, member)
			if len(binding.Members) == 0 {
				policy.Bindings = removeRole(policy.Bindings, role)
			}
			return
		}
	}
}

func removeMember(members []string, member string) []string {
	for i, candidate := range members {
		if candidate == member {
			return append(members[:i], members[i+1:]...)
		}
	}
	return members
}

func removeRole(bindings []*iampb.Binding, role string) []*iampb.Binding {
	for i, binding := range bindings {
		if binding.Role == role {
			return append(bindings[:i], bindings[i+1:]...)
		}
	}
	return bindings
}
