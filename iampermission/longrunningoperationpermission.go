package iampermission

import (
	"strings"

	"cloud.google.com/go/longrunning/autogen/longrunningpb"
	"go.einride.tech/aip/resourcename"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// LongRunningOperationRequest is an interface for long-running operation requests.
type LongRunningOperationRequest interface {
	GetName() string
}

// ResolveLongRunningOperationPermission resolves a permission for a long-running operation.
func ResolveLongRunningOperationPermission(
	operationsPermissions []*iamv1.LongRunningOperationPermissions,
	operationRequest LongRunningOperationRequest,
) (string, bool) {
	_, isListRequest := operationRequest.(*longrunningpb.ListOperationsRequest)
	var match *iamv1.LongRunningOperationPermissions
OperationLoop:
	for _, operationPermissions := range operationsPermissions {
		for _, pattern := range operationPermissions.GetOperation().GetPattern() {
			if isListRequest {
				pattern = parentPattern(pattern)
			}
			if resourcename.Match(pattern, operationRequest.GetName()) {
				match = operationPermissions
				break OperationLoop
			}
		}
	}
	if match == nil {
		return "", false
	}
	switch operationRequest.(type) {
	case *longrunningpb.GetOperationRequest:
		return match.GetGet(), match.GetGet() != ""
	case *longrunningpb.ListOperationsRequest:
		return match.GetList(), match.GetList() != ""
	case *longrunningpb.CancelOperationRequest:
		return match.GetCancel(), match.GetCancel() != ""
	case *longrunningpb.DeleteOperationRequest:
		return match.GetDelete(), match.GetDelete() != ""
	case *longrunningpb.WaitOperationRequest:
		return match.GetWait(), match.GetWait() != ""
	default:
		return "", false
	}
}

func parentPattern(pattern string) string {
	return trimSegment(trimSegment(pattern))
}

func trimSegment(pattern string) string {
	lastIndexSlash := strings.LastIndexByte(pattern, '/')
	if lastIndexSlash == -1 {
		return ""
	}
	return pattern[:lastIndexSlash]
}
