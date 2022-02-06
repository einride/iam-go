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
	if err := ValidateName(role.GetName()); err != nil {
		result.AddFieldError("name", err)
	}
	addTitleViolations(&result, role.GetTitle())
	addDescriptionViolations(&result, role.GetTitle())
	addIncludedPermissionsViolations(&result, role.GetIncludedPermissions())
	return result.Err()
}

// ValidateName checks that an IAM role name is valid.
func ValidateName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("role name must be non-empty")
	}
	if !strings.HasPrefix(name, "roles/") {
		return fmt.Errorf("role name '%s' is not on the format `roles/{service}.{role}`", name)
	}
	roleID := strings.TrimPrefix(name, "roles/")
	if len(roleID) > 64 {
		return fmt.Errorf("role name '%s' has a too long ID, it can be max 64 characters long", name)
	}
	indexOfPeriod := strings.IndexByte(roleID, '.')
	if indexOfPeriod == -1 {
		return fmt.Errorf("role name '%s' is not on the format `roles/{service}.{role}`", name)
	}
	service, role := roleID[:indexOfPeriod], roleID[indexOfPeriod+1:]
	if !isLowerCamelCase(service) || !isLowerCamelCase(role) {
		return fmt.Errorf("each part of role name '%s' must be valid lowerCamelCase", name)
	}
	return nil
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
