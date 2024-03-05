package iamannotations

import (
	"fmt"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// ResolveLongRunningOperationsAuthorizationOptions resolves long-running operation authorization options for a service.
// The provided files registry is used for resolving resource patterns.
func ResolveLongRunningOperationsAuthorizationOptions(
	options *iamv1.LongRunningOperationsAuthorizationOptions,
	files *protoregistry.Files,
	startPackage protoreflect.FullName,
) (*iamv1.LongRunningOperationsAuthorizationOptions, error) {
	result := proto.Clone(options).(*iamv1.LongRunningOperationsAuthorizationOptions)
	for _, operationPermissions := range result.GetOperationPermissions() {
		operation, ok := resolveResource(files, startPackage, operationPermissions.GetOperation().GetType())
		if !ok {
			return nil, fmt.Errorf(
				"resolve long-running operations authorization options in %s: unknown resource %s",
				startPackage,
				operationPermissions.GetOperation().GetType(),
			)
		}
		operationPermissions.Operation.Pattern = append(operationPermissions.Operation.Pattern, operation.GetPattern()...)
	}
	return result, nil
}
