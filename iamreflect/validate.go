package iamreflect

import (
	"fmt"

	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamrole"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// ValidatePredefinedRoles validates a set of predefined roles.
func ValidatePredefinedRoles(roles *iamv1.PredefinedRoles) error {
	var result validation.MessageValidator
	roleNames := make(map[string]struct{}, len(roles.Role))
	for i, role := range roles.Role {
		if _, ok := roleNames[role.Name]; ok {
			result.AddFieldViolation(fmt.Sprintf("role[%d].name", i), "name must be unique among all predefined roles")
		} else {
			roleNames[role.Name] = struct{}{}
		}
		if err := iamrole.Validate(role); err != nil {
			result.AddFieldError(fmt.Sprintf("role[%d]", i), err)
		}
	}
	return result.Err()
}
