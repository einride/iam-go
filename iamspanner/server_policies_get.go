package iamspanner

import (
	"context"

	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamresource"
	"google.golang.org/genproto/googleapis/iam/v1"
)

// GetIamPolicy implements iam.IAMPolicyServer.
func (s *IAMServer) GetIamPolicy(
	ctx context.Context,
	request *iam.GetIamPolicyRequest,
) (*iam.Policy, error) {
	if err := validateGetIamPolicyRequest(request); err != nil {
		return nil, err
	}
	tx := s.client.Single()
	defer tx.Close()
	return s.ReadPolicyInTransaction(ctx, tx, request.Resource)
}

func validateGetIamPolicyRequest(request *iam.GetIamPolicyRequest) error {
	var result validation.MessageValidator
	switch request.Resource {
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
