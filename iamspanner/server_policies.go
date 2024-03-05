package iamspanner

import (
	"bytes"
	"context"
	"fmt"

	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/spanner"
	"go.einride.tech/iam/iamspanner/iamspannerdb"
)

// ReadPolicyInTransaction reads the IAM policy for a resource within the provided transaction.
func (s *IAMServer) ReadPolicyInTransaction(
	ctx context.Context,
	tx ReadTransaction,
	resource string,
) (*iampb.Policy, error) {
	var policy iampb.Policy
	var binding *iampb.Binding
	iamPolicyBindings := iamspannerdb.Descriptor().IamPolicyBindings()
	if err := tx.Read(
		ctx,
		iamPolicyBindings.TableName(),
		spanner.Key{resource}.AsPrefix(),
		[]string{
			iamPolicyBindings.BindingIndex().ColumnName(),
			iamPolicyBindings.Role().ColumnName(),
			iamPolicyBindings.Member().ColumnName(),
		},
	).Do(func(row *spanner.Row) error {
		var bindingIndex int64
		if err := row.Column(0, &bindingIndex); err != nil {
			return err
		}
		var role string
		if err := row.Column(1, &role); err != nil {
			return err
		}
		var member string
		if err := row.Column(2, &member); err != nil {
			return err
		}
		if binding == nil || int(bindingIndex) >= len(policy.GetBindings()) {
			binding = &iampb.Binding{Role: role}
			policy.Bindings = append(policy.Bindings, binding)
		}
		binding.Members = append(binding.Members, member)
		return nil
	}); err != nil {
		return nil, s.handleStorageError(ctx, err)
	}
	etag, err := computeETag(&policy)
	if err != nil {
		return nil, err
	}
	policy.Etag = etag
	return &policy, nil
}

// ReadWritePolicy enables the caller to modify a policy in a read-write transaction.
func (s *IAMServer) ReadWritePolicy(
	ctx context.Context,
	resource string,
	fn func(*iampb.Policy) (*iampb.Policy, error),
) (*iampb.Policy, error) {
	var result *iampb.Policy
	if _, err := s.client.ReadWriteTransaction(
		ctx,
		func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
			policy, err := s.ReadPolicyInTransaction(ctx, tx, resource)
			if err != nil {
				return err
			}
			policy, err = fn(policy)
			if err != nil {
				return err
			}
			result = policy
			mutations := []*spanner.Mutation{deleteIAMPolicyMutation(resource)}
			mutations = append(mutations, insertIAMPolicyMutations(resource, policy)...)
			return tx.BufferWrite(mutations)
		},
	); err != nil {
		return nil, s.handleStorageError(ctx, err)
	}
	result.Etag = nil
	etag, err := computeETag(result)
	if err != nil {
		return nil, err
	}
	result.Etag = etag
	return result, nil
}

// ValidatePolicyFreshnessInTransaction validates the freshness of an IAM policy for a resource
// within the provided transaction.
func (s *IAMServer) ValidatePolicyFreshnessInTransaction(
	ctx context.Context,
	tx ReadTransaction,
	resource string,
	etag []byte,
) (bool, error) {
	if len(etag) == 0 {
		return true, nil
	}
	existingPolicy, err := s.ReadPolicyInTransaction(ctx, tx, resource)
	if err != nil {
		return false, fmt.Errorf("validate freshness: %w", err)
	}
	return bytes.Equal(existingPolicy.GetEtag(), etag), nil
}
