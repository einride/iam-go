package iampolicy

import (
	"fmt"

	"cloud.google.com/go/iam/apiv1/iampb"
	"go.einride.tech/aip/validation"
)

// Validate an IAM policy.
func Validate(policy *iampb.Policy) error {
	var result validation.MessageValidator
	for i, binding := range policy.GetBindings() {
		if len(binding.Role) == 0 {
			result.AddFieldViolation(fmt.Sprintf("bindings[%d].role", i), "missing required field")
		}
		if len(binding.Members) == 0 {
			result.AddFieldViolation(fmt.Sprintf("bindings[%d].members", i), "missing required field")
		}
		for j, member := range binding.Members {
			if len(member) == 0 {
				result.AddFieldViolation(fmt.Sprintf("bindings[%d].members[%d]", i, j), "missing required field")
			}
		}
	}
	return result.Err()
}
