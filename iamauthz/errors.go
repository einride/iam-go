package iamauthz

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// forwardErrorCodes is a workaround for CEL custom functions not being able to return Go errors with gRPC codes.
// Instead, the CEL functions are expected to return the status code as a prefix of the error.
func forwardErrorCodes(err error) error {
	errStr := err.Error()
	for _, code := range []codes.Code{
		codes.InvalidArgument,
		codes.PermissionDenied,
		codes.Unauthenticated,
		codes.DeadlineExceeded,
	} {
		codeStr := code.String()
		if strings.HasPrefix(errStr, codeStr) {
			return status.Error(code, strings.TrimPrefix(strings.TrimPrefix(errStr, codeStr), ": "))
		}
	}
	return err
}
