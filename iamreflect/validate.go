package iamreflect

import (
	"fmt"

	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamcel"
	"go.einride.tech/iam/iampermission"
	"go.einride.tech/iam/iamrole"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
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

func ValidateMethodAuthorizationOptions(
	methodAuthorization *iamv1.MethodAuthorizationOptions,
	method protoreflect.MethodDescriptor,
	files *protoregistry.Files,
) error {
	var result validation.MessageValidator
	switch permissions := methodAuthorization.Permissions.(type) {
	case *iamv1.MethodAuthorizationOptions_Permission:
		if err := iampermission.Validate(permissions.Permission); err != nil {
			result.AddFieldError("permission", err)
		}
	case *iamv1.MethodAuthorizationOptions_ResourcePermissions:
		if err := validateResourcePermissions(permissions.ResourcePermissions, files); err != nil {
			result.AddFieldError("resource_permissions", err)
		}
	default:
		if !methodAuthorization.GetCustom() && !methodAuthorization.GetNone() {
			result.AddFieldViolation("permissions", "one of (permission|resource_permissions) must be specified")
		}
	}
	switch strategy := methodAuthorization.Strategy.(type) {
	case *iamv1.MethodAuthorizationOptions_Before:
		if err := validateBeforeStrategy(strategy, method); err != nil {
			result.AddFieldError("before", err)
		}
	case *iamv1.MethodAuthorizationOptions_After:
		if err := validateAfterStrategy(strategy, method); err != nil {
			result.AddFieldError("after", err)
		}
	case *iamv1.MethodAuthorizationOptions_Custom:
		if !strategy.Custom {
			result.AddFieldViolation("custom", "must be true")
		}
	case *iamv1.MethodAuthorizationOptions_None:
		if !strategy.None {
			result.AddFieldViolation("none", "must be true")
		}
	default:
		result.AddFieldViolation("strategy", "one of (before|after|custom|none) must be specified")
	}
	return result.Err()
}

func validateBeforeStrategy(
	before *iamv1.MethodAuthorizationOptions_Before,
	method protoreflect.MethodDescriptor,
) error {
	env, err := iamcel.NewBeforeEnv(method)
	if err != nil {
		return fmt.Errorf("internal error: %w", err)
	}
	ast, issues := env.Parse(before.Before.GetExpression())
	if issues.Err() != nil {
		return fmt.Errorf("parse error: %w", issues.Err())
	}
	if _, issues := env.Check(ast); issues.Err() != nil {
		return fmt.Errorf("type error: %w", issues.Err())
	}
	return nil
}

func validateAfterStrategy(
	after *iamv1.MethodAuthorizationOptions_After,
	method protoreflect.MethodDescriptor,
) error {
	env, err := iamcel.NewAfterEnv(method)
	if err != nil {
		return fmt.Errorf("internal error: %w", err)
	}
	ast, issues := env.Parse(after.After.GetExpression())
	if issues.Err() != nil {
		return fmt.Errorf("parse error: %w", err)
	}
	if _, issues := env.Check(ast); issues.Err() != nil {
		return fmt.Errorf("type error: %w", issues.Err())
	}
	return nil
}

func validateResourcePermissions(
	resourcePermissions *iamv1.ResourcePermissions,
	files *protoregistry.Files,
) error {
	// TODO: Implement me, and rewrite how resources are resolved during initialization of IAM descriptor.
	//       (To fix current cross-package resolution issue).
	return nil
}

// ValidateLongRunningOperationsAuthorization checks that a long-running operations authorization annotation is valid.
func ValidateLongRunningOperationsAuthorization(authorization *iamv1.LongRunningOperationsAuthorization) error {
	var result validation.MessageValidator
	switch strategy := authorization.Strategy.(type) {
	case *iamv1.LongRunningOperationsAuthorization_Before:
		if !strategy.Before {
			result.AddFieldViolation("before", "must be true")
		}
	case *iamv1.LongRunningOperationsAuthorization_Custom:
		if !strategy.Custom {
			result.AddFieldViolation("custom", "must be true")
		}
	case *iamv1.LongRunningOperationsAuthorization_None:
		if !strategy.None {
			result.AddFieldViolation("none", "must be true")
		}
	default:
		result.AddFieldViolation("strategy", "one of (before|custom|none) must be specified")
	}
	if len(authorization.OperationPermissions) == 0 {
		result.AddFieldViolation("operation_permissions", "required field")
	}
	for i, operationPermissions := range authorization.OperationPermissions {
		if err := validateOperationPermissions(operationPermissions); err != nil {
			result.AddFieldError(fmt.Sprintf("operation_permissions[%d]", i), err)
		}
	}
	return result.Err()
}

func validateOperationPermissions(operationPermissions *iamv1.LongRunningOperationPermissions) error {
	var result validation.MessageValidator
	if operationPermissions.Operation == nil {
		result.AddFieldViolation("operation", "required field")
	} else {
		if operationPermissions.Operation.GetType() == "" {
			result.AddFieldViolation("operation.type", "required field")
		}
		for i, pattern := range operationPermissions.Operation.GetPattern() {
			if err := resourcename.ValidatePattern(pattern); err != nil {
				result.AddFieldError(fmt.Sprintf("operation.type.pattern[%d]", i), err)
			}
		}
	}
	if operationPermissions.List != "" {
		if err := iampermission.Validate(operationPermissions.List); err != nil {
			result.AddFieldError("list", err)
		}
	}
	if operationPermissions.Get != "" {
		if err := iampermission.Validate(operationPermissions.Get); err != nil {
			result.AddFieldError("get", err)
		}
	}
	if operationPermissions.Cancel != "" {
		if err := iampermission.Validate(operationPermissions.Cancel); err != nil {
			result.AddFieldError("cancel", err)
		}
	}
	if operationPermissions.Delete != "" {
		if err := iampermission.Validate(operationPermissions.Delete); err != nil {
			result.AddFieldError("delete", err)
		}
	}
	if operationPermissions.Wait != "" {
		if err := iampermission.Validate(operationPermissions.Wait); err != nil {
			result.AddFieldError("wait", err)
		}
	}
	return result.Err()
}
