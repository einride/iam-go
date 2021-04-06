package iampolicy

import (
	"fmt"

	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

// Validate an IAM policy.
func Validate(policy *iam.Policy) *errdetails.BadRequest {
	var result errdetails.BadRequest
	if len(policy.GetBindings()) == 0 {
		result.FieldViolations = append(result.FieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       "bindings",
			Description: "missing required field",
		})
	}
	for i, binding := range policy.GetBindings() {
		if len(binding.Role) == 0 {
			result.FieldViolations = append(result.FieldViolations, &errdetails.BadRequest_FieldViolation{
				Field:       fmt.Sprintf("bindings[%d].role", i),
				Description: "missing required field",
			})
		}
		if len(binding.Members) == 0 {
			result.FieldViolations = append(result.FieldViolations, &errdetails.BadRequest_FieldViolation{
				Field:       fmt.Sprintf("bindings[%d].members", i),
				Description: "missing required field",
			})
		}
		for j, member := range binding.Members {
			if len(member) == 0 {
				result.FieldViolations = append(result.FieldViolations, &errdetails.BadRequest_FieldViolation{
					Field:       fmt.Sprintf("bindings[%d].members[%d]", i, j),
					Description: "missing required field",
				})
			}
		}
	}
	if len(result.FieldViolations) > 0 {
		return &result
	}
	return nil
}
