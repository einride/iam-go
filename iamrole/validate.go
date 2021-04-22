package iamrole

import (
	"fmt"
	"strings"
	"unicode"

	"go.einride.tech/iam/iampermission"
	"google.golang.org/genproto/googleapis/iam/admin/v1"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

// Validate checks that an IAM role is valid.
func Validate(role *admin.Role) *errdetails.BadRequest {
	var result errdetails.BadRequest
	result.FieldViolations = appendNameViolations(result.FieldViolations, role.GetName())
	result.FieldViolations = appendTitleViolations(result.FieldViolations, role.GetTitle())
	result.FieldViolations = appendDescriptionViolations(result.FieldViolations, role.GetTitle())
	result.FieldViolations = appendIncludedPermissionsViolations(result.FieldViolations, role.GetIncludedPermissions())
	if len(result.FieldViolations) > 0 {
		return &result
	}
	return nil
}

func appendNameViolations(
	fieldViolations []*errdetails.BadRequest_FieldViolation,
	name string,
) []*errdetails.BadRequest_FieldViolation {
	if len(name) == 0 {
		fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       "name",
			Description: "must be non-empty",
		})
	}
	if !strings.HasPrefix(name, "roles/") {
		fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       "name",
			Description: "must have format `roles/{role}`",
		})
	}
	roleID := name[len("roles/"):]
	if len(roleID) > 64 {
		fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       "name",
			Description: "role ID can be max 64 characters long",
		})
	}
	for _, r := range roleID {
		if r > unicode.MaxASCII {
			fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
				Field:       "name",
				Description: "role ID must only contain ASCII characters",
			})
			break
		}
		if !(r == '.' || unicode.In(r, unicode.Letter, unicode.Digit)) {
			fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
				Field:       "name",
				Description: "role ID must only contain alphanumeric characters and periods",
			})
			break
		}
	}
	return fieldViolations
}

func appendTitleViolations(
	fieldViolations []*errdetails.BadRequest_FieldViolation,
	title string,
) []*errdetails.BadRequest_FieldViolation {
	if len(title) == 0 || len(title) > 100 {
		fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       "title",
			Description: "must be non-empty and <= 100 characters",
		})
	}
	return fieldViolations
}

func appendDescriptionViolations(
	fieldViolations []*errdetails.BadRequest_FieldViolation,
	description string,
) []*errdetails.BadRequest_FieldViolation {
	if len(description) == 0 || len(description) > 256 {
		fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       "description",
			Description: "must be non-empty and <= 256 characters",
		})
	}
	return fieldViolations
}

func appendIncludedPermissionsViolations(
	fieldViolations []*errdetails.BadRequest_FieldViolation,
	includedPermissions []string,
) []*errdetails.BadRequest_FieldViolation {
	if len(includedPermissions) == 0 {
		fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       "included_permissions",
			Description: "missing required field",
		})
	}
	for i, includedPermission := range includedPermissions {
		if err := iampermission.Validate(includedPermission); err != nil {
			fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
				Field:       fmt.Sprintf("included_permissions[%d]", i),
				Description: err.Error(),
			})
		}
	}
	return fieldViolations
}
