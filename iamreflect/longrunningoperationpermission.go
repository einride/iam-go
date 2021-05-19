package iamreflect

import (
	"strings"

	"go.einride.tech/aip/resourcename"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/genproto/googleapis/longrunning"
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
	_, isListRequest := operationRequest.(*longrunning.ListOperationsRequest)
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
	case *longrunning.GetOperationRequest:
		return match.Get, match.Get != ""
	case *longrunning.ListOperationsRequest:
		return match.List, match.List != ""
	case *longrunning.CancelOperationRequest:
		return match.Cancel, match.Cancel != ""
	case *longrunning.DeleteOperationRequest:
		return match.Delete, match.Delete != ""
	case *longrunning.WaitOperationRequest:
		return match.Wait, match.Wait != ""
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
