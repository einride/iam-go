package iamannotations

import (
	"fmt"

	"go.einride.tech/iam/iamresource"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// ResolveMethodAuthorizationOptions loads method authorization options for a service.
// The provided files registry is used for resolving resource patterns.
func ResolveMethodAuthorizationOptions(
	options *iamv1.MethodAuthorizationOptions,
	files *protoregistry.Files,
	startPackage protoreflect.FullName,
) (*iamv1.MethodAuthorizationOptions, error) {
	result := proto.Clone(options).(*iamv1.MethodAuthorizationOptions)
	if permissions, ok := result.GetPermissions().(*iamv1.MethodAuthorizationOptions_ResourcePermissions); ok {
		for _, resourcePermission := range permissions.ResourcePermissions.GetResourcePermission() {
			switch {
			case resourcePermission.GetResource().GetType() == iamresource.Root:
				// Root resource requires no pattern resolution.
			case len(resourcePermission.GetResource().GetPattern()) > 0:
				// Resource is annotated with patterns manually. No need to resolve.
			default:
				resource, ok := resolveResource(files, startPackage, resourcePermission.GetResource().GetType())
				if !ok {
					return nil, fmt.Errorf(
						"resolve method authorization options in '%s': unable to resolve resource '%s'",
						startPackage,
						resourcePermission.GetResource().GetType(),
					)
				}
				if len(resource.GetPattern()) == 0 {
					return nil, fmt.Errorf(
						"resolve method authorization options in '%s': resource '%s' has no patterns",
						resourcePermission.GetResource().GetType(),
						startPackage,
					)
				}
				resourcePermission.Resource.Pattern = append(resourcePermission.Resource.Pattern, resource.GetPattern()...)
			}
		}
	}
	return result, nil
}
