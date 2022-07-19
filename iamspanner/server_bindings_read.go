package iamspanner

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iamresource"
	"go.einride.tech/iam/iamspanner/iamspannerdb"
	"google.golang.org/genproto/googleapis/iam/admin/v1"
)

// ReadBindingsByResourcesAndMembers reads all roles bound to the provided members and resources.
func (s *IAMServer) ReadBindingsByResourcesAndMembers(
	ctx context.Context,
	resources []string,
	members []string,
	fn func(ctx context.Context, resource string, role *admin.Role, member string) error,
) error {
	tx := s.client.Single()
	defer tx.Close()
	return s.ReadBindingsByResourcesAndMembersInTransaction(ctx, tx, resources, members, fn)
}

// ReadBindingsByResourcesAndMembersInTransaction reads all roles bound to members and resources
// within the provided Spanner transaction.
// Also considers roles bound to parent resources.
func (s *IAMServer) ReadBindingsByResourcesAndMembersInTransaction(
	ctx context.Context,
	tx ReadTransaction,
	resources []string,
	members []string,
	fn func(ctx context.Context, resource string, role *admin.Role, member string) error,
) error {
	// Deduplicate resources and parents to read.
	resourcesAndParents := make(map[string]struct{}, len(resources))
	// Include root resource.
	resourcesAndParents[iamresource.Root] = struct{}{}
	for _, resource := range resources {
		if resource == iamresource.Root {
			continue
		}
		if !resourcename.ContainsWildcard(resource) {
			resourcesAndParents[resource] = struct{}{}
		}
		resourcename.RangeParents(resource, func(parent string) bool {
			if !resourcename.ContainsWildcard(parent) {
				resourcesAndParents[parent] = struct{}{}
			}
			return true
		})
	}
	// Build deduplicated key ranges to read.
	memberResourceKeySets := make([]spanner.KeySet, 0, len(resources))
	for resource := range resourcesAndParents {
		for _, member := range members {
			memberResourceKeySets = append(memberResourceKeySets, spanner.Key{member, resource}.AsPrefix())
		}
	}
	iamPolicyBindings := iamspannerdb.Descriptor().IamPolicyBindings()
	iamPolicyBindingsByMemberAndResource := iamspannerdb.Descriptor().IamPolicyBindingsByMemberAndResource()
	return tx.ReadWithOptions(
		ctx,
		iamPolicyBindings.TableName(),
		spanner.KeySets(memberResourceKeySets...),
		[]string{
			iamPolicyBindingsByMemberAndResource.Resource().ColumnName(),
			iamPolicyBindingsByMemberAndResource.Role().ColumnName(),
			iamPolicyBindingsByMemberAndResource.Member().ColumnName(),
		},
		&spanner.ReadOptions{
			Index: iamPolicyBindingsByMemberAndResource.IndexName(),
		},
	).Do(func(r *spanner.Row) error {
		var resource string
		if err := r.Column(0, &resource); err != nil {
			return err
		}
		var roleName string
		if err := r.Column(1, &roleName); err != nil {
			return err
		}
		role, ok := s.roles.FindRoleByName(roleName)
		if !ok {
			s.logError(ctx, fmt.Errorf("missing built-in role: %s", roleName))
			return nil
		}
		var member string
		if err := r.Column(2, &member); err != nil {
			return err
		}
		return fn(ctx, resource, role, member)
	})
}

// ReadBindingsByMembersAndPermissions reads all bindings for the provided members and permissions.
func (s *IAMServer) ReadBindingsByMembersAndPermissions(
	ctx context.Context,
	members []string,
	permissions []string,
	fn func(ctx context.Context, resource string, role *admin.Role, member string) error,
) error {
	tx := s.client.Single().WithTimestampBound(spanner.MaxStaleness(5 * time.Second))
	defer tx.Close()
	return s.ReadBindingsByMembersAndPermissionsInTransaction(ctx, tx, members, permissions, fn)
}

// ReadBindingsByMembersAndPermissionsInTransaction reads all bindings for the provided members and permissions,
// within the provided Spanner transaction.
func (s *IAMServer) ReadBindingsByMembersAndPermissionsInTransaction(
	ctx context.Context,
	tx ReadTransaction,
	members []string,
	permissions []string,
	fn func(ctx context.Context, resource string, role *admin.Role, member string) error,
) error {
	memberRoleKeySets := make([]spanner.KeySet, 0, len(members)*len(permissions))
	for _, member := range members {
		for _, permission := range permissions {
			s.roles.RangeRolesByPermission(permission, func(role *admin.Role) bool {
				for _, existingKeySet := range memberRoleKeySets {
					existingKey := existingKeySet.(spanner.KeyRange).Start
					if member == existingKey[0] && role.Name == existingKey[1].(string) {
						return true // already have this member and role
					}
				}
				memberRoleKeySets = append(memberRoleKeySets, spanner.Key{member, role.Name}.AsPrefix())
				return true
			})
		}
	}
	iamPolicyBindings := iamspannerdb.Descriptor().IamPolicyBindings()
	iamPolicyBindingsByMemberAndRole := iamspannerdb.Descriptor().IamPolicyBindingsByMemberAndRole()
	return tx.ReadWithOptions(
		ctx,
		iamPolicyBindings.TableName(),
		spanner.KeySets(memberRoleKeySets...),
		[]string{
			iamPolicyBindingsByMemberAndRole.Resource().ColumnName(),
			iamPolicyBindingsByMemberAndRole.Role().ColumnName(),
			iamPolicyBindingsByMemberAndRole.Member().ColumnName(),
		},
		&spanner.ReadOptions{
			Index: iamPolicyBindingsByMemberAndRole.IndexName(),
		},
	).Do(func(r *spanner.Row) error {
		var resource string
		if err := r.Column(0, &resource); err != nil {
			return err
		}
		var roleName string
		if err := r.Column(1, &roleName); err != nil {
			return err
		}
		role, ok := s.roles.FindRoleByName(roleName)
		if !ok {
			s.logError(ctx, fmt.Errorf("missing built-in role: %s", roleName))
			return nil
		}
		var member string
		if err := r.Column(2, &member); err != nil {
			return err
		}
		return fn(ctx, resource, role, member)
	})
}
