package iamregistry

import (
	"fmt"

	"go.einride.tech/iam/iamrole"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

// ValidateRoles validates a set of predefined roles.
func ValidateRoles(roles *iamv1.Roles) *errdetails.BadRequest {
	var result errdetails.BadRequest
	roleNames := make(map[string]struct{}, len(roles.Role))
	for i, role := range roles.Role {
		if _, ok := roleNames[role.Name]; ok {
			result.FieldViolations = append(result.FieldViolations, &errdetails.BadRequest_FieldViolation{
				Field:       fmt.Sprintf("role[%d].name", i),
				Description: "name must be unique among all predefined roles",
			})
		} else {
			roleNames[role.Name] = struct{}{}
		}
		if err := iamrole.Validate(role); err != nil {
			for _, fieldViolation := range err.FieldViolations {
				fieldViolation.Field = fmt.Sprintf("role[%d].%s", i, fieldViolation.Field)
				result.FieldViolations = append(result.FieldViolations, fieldViolation)
			}
		}
	}
	if len(result.FieldViolations) > 0 {
		return &result
	}
	return nil
}
