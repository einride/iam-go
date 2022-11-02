package iamspanner

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamresource"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SetIamPolicy implements iam.IAMPolicyServer.
func (s *IAMServer) SetIamPolicy(
	ctx context.Context,
	request *iam.SetIamPolicyRequest,
) (*iam.Policy, error) {
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

		for _, fn := range getInsideSetTransactionFuncs(ctx) {
			if err := fn(ctx, tx, request.Policy); err != nil {
				return err
			}
		}

		mutations := []*spanner.Mutation{deleteIAMPolicyMutation(request.Resource)}
		mutations = append(mutations, insertIAMPolicyMutations(request.Resource, request.Policy)...)
		return tx.BufferWrite(mutations)
	}); err != nil {
		return nil, s.handleStorageError(ctx, err)
	}
	if unfresh {
		return nil, status.Error(codes.Aborted, "resource freshness validation failed")
	}
	request.Policy.Etag = nil
	etag, err := computeETag(request.Policy)
	if err != nil {
		return nil, err
	}
	request.Policy.Etag = etag
	return request.Policy, nil
}

// InsideSetIamPolicyTransaction describes a function that is called within the spanner.ReadWriteTransaction in
// IAMServer.SetIamPolicy. The policy provided is the request policy that is applied afterwards. If the function returns
// a non-nil error, the transaction will not be committed.
type InsideSetIamPolicyTransaction func(context.Context, *spanner.ReadWriteTransaction, *iam.Policy) error

type insideSetIamPolicyTransactionCtxKey struct{}

// WithFuncInsideSetIamPolicyTransaction creates a context with function pointers that are called within the
// spanned.ReadWriteTransaction in IAMServer.SetIamPolicy. This makes it possible to write data to other tables
// within the same transaction, e.g. to implement an outbox pattern.
func WithFuncInsideSetIamPolicyTransaction(ctx context.Context, fns ...InsideSetIamPolicyTransaction) context.Context {
	ctxFns := getInsideSetTransactionFuncs(ctx)
	ctxFns = append(ctxFns, fns...)
	return context.WithValue(ctx, insideSetIamPolicyTransactionCtxKey{}, ctxFns)
}

func getInsideSetTransactionFuncs(ctx context.Context) []InsideSetIamPolicyTransaction {
	val := ctx.Value(insideSetIamPolicyTransactionCtxKey{})
	if val == nil {
		return nil
	}
	fns, ok := val.([]InsideSetIamPolicyTransaction)
	if !ok {
		return nil
	}
	return fns
}

func (s *IAMServer) validateSetIamPolicyRequest(request *iam.SetIamPolicyRequest) error {
	var result validation.MessageValidator
	switch request.Resource {
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
		if len(binding.Members) == 0 {
			result.AddFieldViolation(fmt.Sprintf("policy.bindings[%d].members", i), "missing required field")
		}
		for j, member := range binding.Members {
			if err := s.validateMember(member); err != nil {
				result.AddFieldError(fmt.Sprintf("policy.bindings[%d].members[%d]", i, j), err)
			}
		}
	}
	return result.Err()
}
