package iamspanner

import (
	"context"
	"fmt"

	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamresource"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InsideSetIamPolicyTransaction describes a function that is called within the spanner.ReadWriteTransaction in
// IAMServer.SetIamPolicyWithFunctionsInTransaction. The policy provided is the request policy that is applied
// afterwards. If the function returns a non-nil error, the transaction will not be committed.
type InsideSetIamPolicyTransaction func(context.Context, *spanner.ReadWriteTransaction, *iampb.Policy) error

// SetIamPolicy implements iampb.IAMPolicyServer.
func (s *IAMServer) SetIamPolicy(
	ctx context.Context,
	request *iampb.SetIamPolicyRequest,
) (*iampb.Policy, error) {
	return s.SetIamPolicyWithFunctionsInTransaction(ctx, request)
}

// SetIamPolicyWithFunctionsInTransaction handles a SetIamPolicy request but allows for functions to be called
// within the spanner.ReadWriteTransaction.
func (s *IAMServer) SetIamPolicyWithFunctionsInTransaction(
	ctx context.Context,
	request *iampb.SetIamPolicyRequest,
	fns ...InsideSetIamPolicyTransaction,
) (*iampb.Policy, error) {
	if err := s.validateSetIamPolicyRequest(request); err != nil {
		return nil, err
	}
	var unfresh bool
	if _, err := s.client.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		if ok, err := s.ValidatePolicyFreshnessInTransaction(
			ctx, tx, request.GetResource(), request.GetPolicy().GetEtag(),
		); err != nil {
			return err
		} else if !ok {
			unfresh = true
			return nil
		}

		for _, fn := range fns {
			if err := fn(ctx, tx, request.GetPolicy()); err != nil {
				return err
			}
		}

		mutations := []*spanner.Mutation{deleteIAMPolicyMutation(request.GetResource())}
		mutations = append(mutations, insertIAMPolicyMutations(request.GetResource(), request.GetPolicy())...)
		return tx.BufferWrite(mutations)
	}); err != nil {
		return nil, s.handleStorageError(ctx, err)
	}
	if unfresh {
		return nil, status.Error(codes.Aborted, "resource freshness validation failed")
	}

	returned := request.GetPolicy()
	if returned == nil {
		returned = &iampb.Policy{}
	}

	returned.Etag = nil
	etag, err := computeETag(returned)
	if err != nil {
		return nil, err
	}
	returned.Etag = etag
	return returned, nil
}

func (s *IAMServer) validateSetIamPolicyRequest(request *iampb.SetIamPolicyRequest) error {
	var result validation.MessageValidator
	switch request.GetResource() {
	case iamresource.Root: // OK
	case "":
		result.AddFieldViolation("resource", "missing required field")
	default:
		if err := resourcename.Validate(request.GetResource()); err != nil {
			result.AddFieldError("resource", err)
		} else if resourcename.ContainsWildcard(request.GetResource()) {
			result.AddFieldViolation("resource", "must not contain wildcard")
		}
	}

	roleSet := map[string]bool{}
	for i, binding := range request.GetPolicy().GetBindings() {
		if binding.GetRole() == "" {
			result.AddFieldViolation(fmt.Sprintf("policy.bindings[%d].role", i), "missing required field")
		}
		if _, ok := s.roles.FindRoleByName(binding.GetRole()); !ok {
			result.AddFieldViolation(
				fmt.Sprintf("policy.bindings[%d].role", i),
				"unknown role: '%s'",
				binding.GetRole(),
			)
		}
		_, ok := roleSet[binding.GetRole()]
		if ok {
			result.AddFieldViolation(
				fmt.Sprintf("policy.bindings[%d].role", i),
				"duplicate role: '%s'",
				binding.GetRole(),
			)
		}
		roleSet[binding.GetRole()] = true

		if len(binding.GetMembers()) == 0 {
			result.AddFieldViolation(fmt.Sprintf("policy.bindings[%d].members", i), "missing required field")
		}
		memberSet := map[string]bool{}
		for j, member := range binding.GetMembers() {
			if err := s.validateMember(member); err != nil {
				result.AddFieldError(fmt.Sprintf("policy.bindings[%d].members[%d]", i, j), err)
			}
			_, ok := memberSet[member]
			if ok {
				// duplicate member
				result.AddFieldViolation(fmt.Sprintf("policy.bindings[%d].members[%d]", i, j), "duplicate member")
			}
			memberSet[member] = true
		}
	}
	return result.Err()
}
