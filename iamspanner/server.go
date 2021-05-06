package iamspanner

import (
	"bytes"
	"context"
	"fmt"
	"hash/crc32"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iammember"
	"go.einride.tech/iam/iampolicy"
	"go.einride.tech/iam/iamregistry"
	"go.einride.tech/iam/iamresource"
	"go.einride.tech/iam/iamrole"
	"go.einride.tech/iam/iamspanner/iamspannerdb"
	"google.golang.org/genproto/googleapis/iam/admin/v1"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Server is a Spanner implementation of the iam.IAMPolicyServer interface.
type Server struct {
	iam.UnimplementedIAMPolicyServer
	admin.UnimplementedIAMServer
	client         *spanner.Client
	roles          *iamregistry.Roles
	memberResolver iammember.Resolver
	config         ServerConfig
}

// ServerConfig configures a Spanner IAM policy server.
type ServerConfig struct {
	ErrorHook func(context.Context, error)
}

// ReadTransaction is an interface for Spanner read transactions.
type ReadTransaction interface {
	Read(context.Context, string, spanner.KeySet, []string) *spanner.RowIterator
	ReadWithOptions(context.Context, string, spanner.KeySet, []string, *spanner.ReadOptions) *spanner.RowIterator
}

// NewServer creates a new Spanner IAM policy server.
func NewServer(
	client *spanner.Client,
	roles *iamregistry.Roles,
	memberResolver iammember.Resolver,
	config ServerConfig,
) (*Server, error) {
	s := &Server{
		client:         client,
		config:         config,
		roles:          roles,
		memberResolver: memberResolver,
	}
	return s, nil
}

// SetIamPolicy implements iam.IAMPolicyServer.
func (s *Server) SetIamPolicy(
	ctx context.Context,
	request *iam.SetIamPolicyRequest,
) (*iam.Policy, error) {
	if err := s.validateSetIamPolicyRequest(ctx, request); err != nil {
		return nil, err
	}
	var unfresh bool
	if _, err := s.client.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		if ok, err := s.ValidateIamPolicyFreshnessInTransaction(
			ctx, tx, request.GetResource(), request.GetPolicy().GetEtag(),
		); err != nil {
			return err
		} else if !ok {
			unfresh = true
			return nil
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

// GetIamPolicy implements iam.IAMPolicyServer.
func (s *Server) GetIamPolicy(
	ctx context.Context,
	request *iam.GetIamPolicyRequest,
) (*iam.Policy, error) {
	tx := s.client.Single()
	defer tx.Close()
	return s.QueryIamPolicyInTransaction(ctx, tx, request.Resource)
}

// TestIamPermissions implements iam.IAMPolicyServer.
func (s *Server) TestIamPermissions(
	ctx context.Context,
	request *iam.TestIamPermissionsRequest,
) (*iam.TestIamPermissionsResponse, error) {
	members, err := s.resolveMembers(ctx)
	if err != nil {
		return nil, err
	}
	permissions := make(map[string]struct{}, len(request.Permissions))
	tx := s.client.Single()
	defer tx.Close()
	if err := s.ReadRolesBoundToMembersAndResourcesInTransaction(
		ctx,
		tx,
		members,
		[]string{request.Resource},
		func(ctx context.Context, _, _ string, role *admin.Role) error {
			for _, permission := range request.Permissions {
				if s.roles.RoleHasPermission(role.Name, permission) {
					permissions[permission] = struct{}{}
				}
			}
			return nil
		},
	); err != nil {
		return nil, s.handleStorageError(ctx, err)
	}
	response := &iam.TestIamPermissionsResponse{
		Permissions: make([]string, 0, len(permissions)),
	}
	for _, permission := range request.Permissions {
		if _, ok := permissions[permission]; ok {
			response.Permissions = append(response.Permissions, permission)
		}
	}
	return response, nil
}

// TestPermissionOnResource tests if the caller has the specified permission on the specified resource.
func (s *Server) TestPermissionOnResource(
	ctx context.Context,
	permission string,
	resource string,
) (bool, error) {
	result, err := s.TestPermissionOnResources(ctx, permission, []string{resource})
	if err != nil {
		return false, err
	}
	return result[resource], nil
}

// TestPermissionOnResources tests if the caller has the specified permission on the specified resources.
func (s *Server) TestPermissionOnResources(
	ctx context.Context,
	permission string,
	resources []string,
) (map[string]bool, error) {
	members, err := s.resolveMembers(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[string]bool, len(resources))
	tx := s.client.Single()
	defer tx.Close()
	if err := s.ReadRolesBoundToMembersAndResourcesInTransaction(
		ctx,
		tx,
		members,
		resources,
		func(ctx context.Context, _ string, boundResource string, role *admin.Role) error {
			for _, resource := range resources {
				result[resource] = result[resource] ||
					(boundResource == iamresource.Root ||
						resource == boundResource ||
						resourcename.HasParent(resource, boundResource) &&
							iamrole.HasPermission(role, permission))
			}
			return nil
		},
	); err != nil {
		return nil, s.handleStorageError(ctx, err)
	}
	return result, nil
}

// ReadRolesBoundToMembersAndResources reads all roles bound to the provided members and resources.
func (s *Server) ReadRolesBoundToMembersAndResources(
	ctx context.Context,
	members []string,
	resources []string,
	fn func(ctx context.Context, member, resource string, role *admin.Role) error,
) error {
	tx := s.client.Single()
	defer tx.Close()
	return s.ReadRolesBoundToMembersAndResourcesInTransaction(ctx, tx, members, resources, fn)
}

// ReadRolesBoundToMembersAndResourcesInTransaction reads all roles bound to members and resources
// within the provided Spanner transaction.
// Also considers roles bound to parent resources.
func (s *Server) ReadRolesBoundToMembersAndResourcesInTransaction(
	ctx context.Context,
	tx ReadTransaction,
	members []string,
	resources []string,
	fn func(ctx context.Context, member, resource string, role *admin.Role) error,
) error {
	// Deduplicate resources and parents to read.
	resourcesAndParents := make(map[string]struct{}, len(resources))
	// Include root resource.
	resourcesAndParents[iamresource.Root] = struct{}{}
	for _, resource := range resources {
		if resource == iamresource.Root {
			continue
		}
		resourcesAndParents[resource] = struct{}{}
		resourcename.RangeParents(resource, func(parent string) bool {
			resourcesAndParents[parent] = struct{}{}
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
			iamPolicyBindingsByMemberAndResource.Member().ColumnName(),
			iamPolicyBindingsByMemberAndResource.Resource().ColumnName(),
			iamPolicyBindingsByMemberAndResource.Role().ColumnName(),
		},
		&spanner.ReadOptions{
			Index: iamPolicyBindingsByMemberAndResource.IndexName(),
		},
	).Do(func(r *spanner.Row) error {
		var member string
		if err := r.Column(0, &member); err != nil {
			return err
		}
		var resource string
		if err := r.Column(1, &resource); err != nil {
			return err
		}
		var roleName string
		if err := r.Column(2, &roleName); err != nil {
			return err
		}
		role, ok := s.roles.FindRoleByName(roleName)
		if !ok {
			return status.Errorf(codes.Internal, "missing built-in role: %s", roleName)
		}
		return fn(ctx, member, resource, role)
	})
}

// QueryResourcesBoundToMemberAndPermission reads all resources bound to the member and permission.
func (s *Server) QueryResourcesBoundToMemberAndPermission(
	ctx context.Context,
	member string,
	permission string,
) ([]string, error) {
	tx := s.client.Single()
	defer tx.Close()
	return s.QueryResourcesBoundToMemberAndPermissionInTransaction(ctx, tx, member, permission)
}

// QueryResourcesBoundToMemberAndPermissionInTransaction reads all resources bound to the member and permission,
// within the provided Spanner transaction.
func (s *Server) QueryResourcesBoundToMemberAndPermissionInTransaction(
	ctx context.Context,
	tx ReadTransaction,
	member string,
	permission string,
) ([]string, error) {
	memberRoleKeySets := make([]spanner.KeySet, 0, s.roles.Count())
	s.roles.RangeRolesByPermission(permission, func(role *admin.Role) bool {
		memberRoleKeySets = append(memberRoleKeySets, spanner.Key{member, role}.AsPrefix())
		return true
	})
	iamPolicyBindings := iamspannerdb.Descriptor().IamPolicyBindings()
	iamPolicyBindingsByMemberAndRole := iamspannerdb.Descriptor().IamPolicyBindingsByMemberAndRole()
	var resources []string
	if err := tx.ReadWithOptions(
		ctx,
		iamPolicyBindings.TableName(),
		spanner.KeySets(memberRoleKeySets...),
		[]string{
			iamPolicyBindingsByMemberAndRole.Resource().ColumnName(),
		},
		&spanner.ReadOptions{
			Index: iamPolicyBindingsByMemberAndRole.IndexName(),
		},
	).Do(func(r *spanner.Row) error {
		var resource string
		if err := r.Column(0, &resource); err != nil {
			return err
		}
		resources = append(resources, resource)
		return nil
	}); err != nil {
		return nil, s.handleStorageError(ctx, err)
	}
	return resources, nil
}

// QueryIamPolicyInTransaction queries the IAM policy for a resource within the provided transaction.
func (s *Server) QueryIamPolicyInTransaction(
	ctx context.Context,
	tx ReadTransaction,
	resource string,
) (*iam.Policy, error) {
	var policy iam.Policy
	var binding *iam.Binding
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
		if binding == nil || int(bindingIndex) >= len(policy.Bindings) {
			binding = &iam.Binding{Role: role}
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

// ValidateIamPolicyFreshnessInTransaction validates the freshness of an IAM policy for a resource
// within the provided transaction.
func (s *Server) ValidateIamPolicyFreshnessInTransaction(
	ctx context.Context,
	tx ReadTransaction,
	resource string,
	etag []byte,
) (bool, error) {
	if len(etag) == 0 {
		return true, nil
	}
	existingPolicy, err := s.QueryIamPolicyInTransaction(ctx, tx, resource)
	if err != nil {
		return false, fmt.Errorf("validate freshness: %w", err)
	}
	return bytes.Equal(existingPolicy.Etag, etag), nil
}

func deleteIAMPolicyMutation(resource string) *spanner.Mutation {
	return spanner.Delete(
		iamspannerdb.Descriptor().IamPolicyBindings().TableName(),
		spanner.Key{resource}.AsPrefix(),
	)
}

func insertIAMPolicyMutations(resource string, policy *iam.Policy) []*spanner.Mutation {
	var mutations []*spanner.Mutation
	for i, binding := range policy.GetBindings() {
		for j, member := range binding.GetMembers() {
			row := iamspannerdb.IamPolicyBindingsRow{
				Resource:     resource,
				BindingIndex: int64(i),
				Role:         binding.Role,
				MemberIndex:  int64(j),
				Member:       member,
			}
			mutations = append(mutations, spanner.Insert(row.Mutate()))
		}
	}
	return mutations
}

func (s *Server) validateSetIamPolicyRequest(ctx context.Context, request *iam.SetIamPolicyRequest) error {
	var fieldViolations []*errdetails.BadRequest_FieldViolation
	if len(request.Resource) == 0 {
		fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       "resource",
			Description: "missing required field",
		})
	}
	for _, fieldViolation := range iampolicy.Validate(request.GetPolicy()).GetFieldViolations() {
		fieldViolation.Field = fmt.Sprintf("policy.%s", fieldViolation.Field)
		fieldViolations = append(fieldViolations, fieldViolation)
	}
	if len(fieldViolations) > 0 {
		result, err := status.New(codes.InvalidArgument, "bad request").WithDetails(&errdetails.BadRequest{
			FieldViolations: fieldViolations,
		})
		if err != nil {
			s.logError(ctx, err)
			return status.Error(codes.Internal, "failed to attach error details")
		}
		return result.Err()
	}
	return nil
}

func (s *Server) logError(ctx context.Context, err error) {
	if s.config.ErrorHook != nil {
		s.config.ErrorHook(ctx, err)
	}
}

func (s *Server) handleStorageError(ctx context.Context, err error) error {
	s.logError(ctx, err)
	switch code := status.Code(err); code {
	case codes.Aborted, codes.Canceled, codes.DeadlineExceeded, codes.Unavailable:
		return status.Error(code, "transient storage error")
	default:
		return status.Error(codes.Internal, "storage error")
	}
}

func (s *Server) resolveMembers(ctx context.Context) ([]string, error) {
	members, err := s.memberResolver.ResolveIAMMembers(ctx)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func computeETag(policy *iam.Policy) ([]byte, error) {
	data, err := proto.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("compute etag: %w", err)
	}
	return []byte(fmt.Sprintf("W/%d-%08X", len(data), crc32.ChecksumIEEE(data))), nil
}
