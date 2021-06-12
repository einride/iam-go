package iamrole

import (
	"fmt"
	"strings"
	"unicode"

	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iampermission"
	"google.golang.org/genproto/googleapis/iam/admin/v1"
)

// Validate checks that an IAM role is valid.
func Validate(role *admin.Role) error {
	var result validation.MessageValidator
	addNameViolations(&result, role.GetName())
	addTitleViolations(&result, role.GetTitle())
	addDescriptionViolations(&result, role.GetTitle())
	addIncludedPermissionsViolations(&result, role.GetIncludedPermissions())
	return result.Err()
}

func addNameViolations(result *validation.MessageValidator, name string) {
	if len(name) == 0 {
		result.AddFieldViolation("name", "must be non-empty")
		return
	}
	if !strings.HasPrefix(name, "roles/") {
		result.AddFieldViolation("name", "must have format `roles/{service}.{role}`")
		return
	}
	roleID := strings.TrimPrefix(name, "roles/")
	if len(roleID) > 64 {
		result.AddFieldViolation("name", "role ID can be max 64 characters long")
	}
	if indexOfPeriod := strings.IndexByte(roleID, '.'); indexOfPeriod == -1 {
		result.AddFieldViolation("name", "must be on the format `roles/{service}.{role}`")
	} else {
		service, role := roleID[:indexOfPeriod], roleID[indexOfPeriod+1:]
		if !isLowerCamelCase(service) || !isLowerCamelCase(role) {
			result.AddFieldViolation("name", "each part of `roles/{service}.{role}` must be valid lowerCamelCase")
		}
	}
}

func addTitleViolations(result *validation.MessageValidator, title string) {
	if len(title) == 0 || len(title) > 100 {
		result.AddFieldViolation("title", "must be non-empty and <= 100 characters")
	}
}

func addDescriptionViolations(result *validation.MessageValidator, description string) {
	if len(description) == 0 || len(description) > 256 {
		result.AddFieldViolation("description", "must be non-empty and <= 256 characters")
	}
}

func addIncludedPermissionsViolations(result *validation.MessageValidator, includedPermissions []string) {
	if len(includedPermissions) == 0 {
		result.AddFieldViolation("included_permissions", "missing required field")
	}
	for i, includedPermission := range includedPermissions {
		if err := iampermission.Validate(includedPermission); err != nil {
			result.AddFieldError(fmt.Sprintf("included_permissions[%d]", i), err)
		}
	}
}

func isLowerCamelCase(s string) bool {
	for i, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
		if i == 0 {
			if !unicode.IsLower(r) {
				return false
			}
		} else {
			if !unicode.In(r, unicode.Letter, unicode.Digit) {
				return false
			}
		}
	}
	return true
}
