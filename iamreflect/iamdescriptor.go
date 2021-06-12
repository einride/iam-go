package iamreflect

import (
	"fmt"

	"go.einride.tech/iam/iampermission"
	"go.einride.tech/iam/iamresource"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// IAMDescriptor represents an RPC service's IAM specification.
type IAMDescriptor struct {
	// PredefinedRoles are the service's predefined IAM roles.
	PredefinedRoles *iamv1.PredefinedRoles
	// LongRunningOperationsAuthorization is the service's configuration for authorization of long-running operations.
	LongRunningOperationsAuthorization *iamv1.LongRunningOperationsAuthorization
	// MethodAuthorizationOptions is a mapping from full method name to the method's authorization options.
	MethodAuthorizationOptions map[protoreflect.FullName]*iamv1.MethodAuthorizationOptions
	// RequestAuthorizationOptions is a mapping from full request name to the method's authorization options.
	RequestAuthorizationOptions map[protoreflect.FullName]*iamv1.MethodAuthorizationOptions
}

// NewIAMDescriptor creates a new IAMDescriptor from the provided service descriptor.
// Uses the provided files to resolve resource name patterns.
func NewIAMDescriptor(service protoreflect.ServiceDescriptor, files *protoregistry.Files) (*IAMDescriptor, error) {
	result := IAMDescriptor{
		MethodAuthorizationOptions: make(
			map[protoreflect.FullName]*iamv1.MethodAuthorizationOptions, service.Methods().Len(),
		),
		RequestAuthorizationOptions: make(
			map[protoreflect.FullName]*iamv1.MethodAuthorizationOptions, service.Methods().Len(),
		),
	}
	if predefinedRoles := proto.GetExtension(
		service.Options(), iamv1.E_PredefinedRoles,
	).(*iamv1.PredefinedRoles); predefinedRoles != nil {
		result.PredefinedRoles = proto.Clone(predefinedRoles).(*iamv1.PredefinedRoles)
	}
	if longRunningOperationsAuthorization := proto.GetExtension(
		service.Options(), iamv1.E_LongRunningOperationsAuthorization,
	).(*iamv1.LongRunningOperationsAuthorization); longRunningOperationsAuthorization != nil {
		result.LongRunningOperationsAuthorization = proto.Clone(
			longRunningOperationsAuthorization,
		).(*iamv1.LongRunningOperationsAuthorization)
	}
	for i := 0; i < service.Methods().Len(); i++ {
		method := service.Methods().Get(i)
		if methodAuthorizationOptions := proto.GetExtension(
			method.Options(), iamv1.E_MethodAuthorization,
		).(*iamv1.MethodAuthorizationOptions); methodAuthorizationOptions != nil {
			methodAuthorizationOptions = proto.Clone(methodAuthorizationOptions).(*iamv1.MethodAuthorizationOptions)
			result.MethodAuthorizationOptions[method.FullName()] = methodAuthorizationOptions
			if _, ok := result.RequestAuthorizationOptions[method.Input().FullName()]; ok {
				return nil, fmt.Errorf(
					"new %s IAM descriptor: service uses the same request %s for multiple methods",
					service.Name(),
					method.Input().FullName(),
				)
			}
			result.RequestAuthorizationOptions[method.Input().FullName()] = methodAuthorizationOptions
		}
	}
	// Resolve method authorization resources.
	for method, methodAuthorizationOptions := range result.MethodAuthorizationOptions {
		resourcePermissions := methodAuthorizationOptions.GetResourcePermissions()
		if resourcePermissions == nil {
			continue
		}
		for _, resourcePermission := range resourcePermissions.ResourcePermission {
			switch {
			case resourcePermission.Resource.GetType() == iamresource.Root:
				// Root resource requires no pattern resolution.
				continue
			case len(resourcePermission.Resource.GetPattern()) > 0:
				// Resource is annotated with patterns manually. No need to resolve.
				continue
			}
			resource, ok := resolveResource(resourcePermission.Resource.GetType(), service, files)
			if !ok {
				return nil, fmt.Errorf(
					"new %s IAM descriptor: unable to resolve resource '%s' patterns for method '%s'",
					service.Name(),
					resourcePermission.Resource.GetType(),
					method,
				)
			}
			if len(resource.Pattern) == 0 {
				return nil, fmt.Errorf(
					"new %s IAM descriptor: resource '%s' has no patterns for method '%s'",
					service.Name(),
					resource.GetType(),
					method,
				)
			}
			resourcePermission.Resource = proto.Clone(resource).(*annotations.ResourceDescriptor)
		}
	}
	// Resolve long-running operation authorization operations.
	for _, operationPermissions := range result.LongRunningOperationsAuthorization.GetOperationPermissions() {
		if len(operationPermissions.Operation.GetPattern()) > 0 {
			// Operation is annotated with patterns manually. No need to resolve.
			continue
		}
		operation, ok := resolveResource(operationPermissions.Operation.GetType(), service, files)
		if !ok {
			return nil, fmt.Errorf(
				"new %s IAM descriptor: unable to resolve operation '%s' patterns",
				service.Name(),
				operation.GetType(),
			)
		}
		if len(operation.Pattern) == 0 {
			return nil, fmt.Errorf(
				"new %s IAM descriptor: operation '%s' has no patterns",
				service.Name(),
				operation.GetType(),
			)
		}
		operationPermissions.Operation = proto.Clone(operation).(*annotations.ResourceDescriptor)
	}
	return &result, nil
}

func (d *IAMDescriptor) ResolvePermissionByRequestAndResource(
	request proto.Message,
	resource string,
) (string, bool) {
	methodAuthorizationOptions, ok := d.FindMethodAuthorizationOptionsByRequest(request)
	if !ok {
		return "", false
	}
	return iampermission.ResolveMethodPermission(methodAuthorizationOptions, resource)
}

func (d *IAMDescriptor) FindMethodAuthorizationOptionsByRequest(
	request proto.Message,
) (*iamv1.MethodAuthorizationOptions, bool) {
	result, ok := d.RequestAuthorizationOptions[request.ProtoReflect().Descriptor().FullName()]
	return result, ok
}

func resolveResource(
	resourceType string,
	service protoreflect.ServiceDescriptor,
	files *protoregistry.Files,
) (*annotations.ResourceDescriptor, bool) {
	var result *annotations.ResourceDescriptor
	var searchMessagesFn func(protoreflect.MessageDescriptors) bool
	searchMessagesFn = func(messages protoreflect.MessageDescriptors) bool {
		for i := 0; i < messages.Len(); i++ {
			message := messages.Get(i)
			if resource := proto.GetExtension(
				message.Options(), annotations.E_Resource,
			).(*annotations.ResourceDescriptor); resource != nil {
				if resource.Type == resourceType {
					result = resource
					return false
				}
			}
			if !searchMessagesFn(message.Messages()) {
				return false
			}
		}
		return true
	}
	searchFileFn := func(file protoreflect.FileDescriptor) bool {
		// Search file annotations.
		for _, resource := range proto.GetExtension(
			file.Options(), annotations.E_ResourceDefinition,
		).([]*annotations.ResourceDescriptor) {
			if resource.Type == resourceType {
				result = resource
				return false
			}
		}
		return searchMessagesFn(file.Messages())
	}
	// Start with a narrow search in the same package.
	files.RangeFilesByPackage(service.ParentFile().Package(), searchFileFn)
	if result != nil {
		return result, true
	}
	// Fall back to a broad search of all files.
	files.RangeFiles(searchFileFn)
	return result, result != nil
}
