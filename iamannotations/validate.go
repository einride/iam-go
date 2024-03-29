package iamannotations

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamcel"
	"go.einride.tech/iam/iampermission"
	"go.einride.tech/iam/iamresource"
	"go.einride.tech/iam/iamrole"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// ValidatePredefinedRoles validates a set of predefined roles.
func ValidatePredefinedRoles(roles *iamv1.PredefinedRoles) error {
	var result validation.MessageValidator
	roleNames := make(map[string]struct{}, len(roles.GetRole()))
	for i, role := range roles.GetRole() {
		if _, ok := roleNames[role.GetName()]; ok {
			result.AddFieldViolation(fmt.Sprintf("role[%d].name", i), "name must be unique among all predefined roles")
		} else {
			roleNames[role.GetName()] = struct{}{}
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
	switch permissions := methodAuthorization.GetPermissions().(type) {
	case *iamv1.MethodAuthorizationOptions_Permission:
		if err := iampermission.Validate(permissions.Permission); err != nil {
			result.AddFieldError("permission", err)
		}
	case *iamv1.MethodAuthorizationOptions_ResourcePermissions:
		if err := validateResourcePermissions(
			permissions.ResourcePermissions, files, method.ParentFile().Package(),
		); err != nil {
			result.AddFieldError("resource_permissions", err)
		}
	default:
		if !methodAuthorization.GetCustom() && !methodAuthorization.GetNone() {
			result.AddFieldViolation("permissions", "one of (permission|resource_permissions) must be specified")
		}
	}
	switch strategy := methodAuthorization.GetStrategy().(type) {
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
	checkedAst, issues := env.Check(ast)
	if issues.Err() != nil {
		return fmt.Errorf("type error: %w", issues.Err())
	}
	if checkedAst.OutputType() != cel.BoolType {
		return fmt.Errorf("invalid output type: %v", ast.OutputType())
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
	checkedAst, issues := env.Check(ast)
	if issues.Err() != nil {
		return fmt.Errorf("type error: %w", issues.Err())
	}
	if checkedAst.OutputType() != cel.BoolType {
		return fmt.Errorf("invalid output type: %v", ast.OutputType())
	}
	return nil
}

func validateResourcePermissions(
	resourcePermissions *iamv1.ResourcePermissions,
	files *protoregistry.Files,
	startPackage protoreflect.FullName,
) error {
	var result validation.MessageValidator
	if len(resourcePermissions.GetResourcePermission()) == 0 {
		result.AddFieldViolation("resource_permission", "at least one resource permission is required")
	}
	for i, resourcePermission := range resourcePermissions.GetResourcePermission() {
		switch {
		case resourcePermission.GetResource().GetType() == "":
			result.AddFieldViolation(fmt.Sprintf("resource_permission[%d].resource.type", i), "missing required field")
		case resourcePermission.GetResource().GetType() == iamresource.Root:
			if len(resourcePermission.GetResource().GetPattern()) > 0 {
				result.AddFieldViolation(
					fmt.Sprintf("resource_permission[%d]", i), "root resource must not have patterns",
				)
			}
		case len(resourcePermission.GetResource().GetPattern()) > 0:
			for j, pattern := range resourcePermission.GetResource().GetPattern() {
				if err := resourcename.ValidatePattern(pattern); err != nil {
					result.AddFieldError(fmt.Sprintf("resource_permission[%d].pattern[%d]", i, j), err)
				}
			}
		default:
			if resource, ok := resolveResource(files, startPackage, resourcePermission.GetResource().GetType()); ok {
				if len(resource.GetPattern()) == 0 {
					result.AddFieldViolation(
						fmt.Sprintf("resource_permission[%d].resource.type", i),
						"resolved resource '%s' has no patterns",
						resourcePermission.GetResource().GetType(),
					)
				}
			} else {
				result.AddFieldViolation(
					fmt.Sprintf("resource_permission[%d].resource.type", i),
					"unable to resolve resource '%s'",
					resourcePermission.GetResource().GetType(),
				)
			}
		}
	}
	return result.Err()
}

// ValidateLongRunningOperationsAuthorization checks that a long-running operations authorization annotation is valid.
func ValidateLongRunningOperationsAuthorization(
	options *iamv1.LongRunningOperationsAuthorizationOptions,
) error {
	var result validation.MessageValidator
	switch strategy := options.GetStrategy().(type) {
	case *iamv1.LongRunningOperationsAuthorizationOptions_Before:
		if !strategy.Before {
			result.AddFieldViolation("before", "must be true")
		}
	case *iamv1.LongRunningOperationsAuthorizationOptions_Custom:
		if !strategy.Custom {
			result.AddFieldViolation("custom", "must be true")
		}
	case *iamv1.LongRunningOperationsAuthorizationOptions_None:
		if !strategy.None {
			result.AddFieldViolation("none", "must be true")
		}
	default:
		result.AddFieldViolation("strategy", "one of (before|custom|none) must be specified")
	}
	if len(options.GetOperationPermissions()) == 0 {
		result.AddFieldViolation("operation_permissions", "required field")
	}
	for i, operationPermissions := range options.GetOperationPermissions() {
		if err := validateOperationPermissions(operationPermissions); err != nil {
			result.AddFieldError(fmt.Sprintf("operation_permissions[%d]", i), err)
		}
	}
	return result.Err()
}

func validateOperationPermissions(operationPermissions *iamv1.LongRunningOperationPermissions) error {
	var result validation.MessageValidator
	if operationPermissions.GetOperation() == nil {
		result.AddFieldViolation("operation", "required field")
	} else {
		if operationPermissions.GetOperation().GetType() == "" {
			result.AddFieldViolation("operation.type", "required field")
		}
		for i, pattern := range operationPermissions.GetOperation().GetPattern() {
			if err := resourcename.ValidatePattern(pattern); err != nil {
				result.AddFieldError(fmt.Sprintf("operation.type.pattern[%d]", i), err)
			}
		}
	}
	if operationPermissions.GetList() != "" {
		if err := iampermission.Validate(operationPermissions.GetList()); err != nil {
			result.AddFieldError("list", err)
		}
	}
	if operationPermissions.GetGet() != "" {
		if err := iampermission.Validate(operationPermissions.GetGet()); err != nil {
			result.AddFieldError("get", err)
		}
	}
	if operationPermissions.GetCancel() != "" {
		if err := iampermission.Validate(operationPermissions.GetCancel()); err != nil {
			result.AddFieldError("cancel", err)
		}
	}
	if operationPermissions.GetDelete() != "" {
		if err := iampermission.Validate(operationPermissions.GetDelete()); err != nil {
			result.AddFieldError("delete", err)
		}
	}
	if operationPermissions.GetWait() != "" {
		if err := iampermission.Validate(operationPermissions.GetWait()); err != nil {
			result.AddFieldError("wait", err)
		}
	}
	return result.Err()
}
