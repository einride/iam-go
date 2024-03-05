package iamspanner

import (
	"context"

	"cloud.google.com/go/iam/apiv1/iampb"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamresource"
)

// GetIamPolicy implements iampb.IAMPolicyServer.
func (s *IAMServer) GetIamPolicy(
	ctx context.Context,
	request *iampb.GetIamPolicyRequest,
) (*iampb.Policy, error) {
	if err := validateGetIamPolicyRequest(request); err != nil {
		return nil, err
	}
	tx := s.client.Single()
	defer tx.Close()
	return s.ReadPolicyInTransaction(ctx, tx, request.GetResource())
}

func validateGetIamPolicyRequest(request *iampb.GetIamPolicyRequest) error {
	var result validation.MessageValidator
	switch request.GetResource() {
	case iamresource.Root: // OK
	case "":
		result.AddFieldViolation("resource", "missing required field")
	default:
		if err := resourcename.Validate(request.GetResource()); err != nil {
			result.AddFieldError("resource", err)
		}
		if resourcename.ContainsWildcard(request.GetResource()) {
			result.AddFieldViolation("resource", "must not contain wildcard")
		}
	}
	return result.Err()
}
