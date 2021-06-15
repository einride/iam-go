package iamspanner

import (
	"context"
	"fmt"
	"hash/crc32"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iammember"
	"go.einride.tech/iam/iamregistry"
	"go.einride.tech/iam/iamresource"
	"go.einride.tech/iam/iamrole"
	"go.einride.tech/iam/iamspanner/iamspannerdb"
	"google.golang.org/genproto/googleapis/iam/admin/v1"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// IAMServer is a Spanner implementation of the iam.IAMPolicyServer interface.
type IAMServer struct {
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

// NewIAMServer creates a new Spanner IAM policy server.
func NewIAMServer(
	client *spanner.Client,
	roles []*admin.Role,
	memberResolver iammember.Resolver,
	config ServerConfig,
) (*IAMServer, error) {
	rolesRegistry, err := iamregistry.NewRoles(roles...)
	if err != nil {
		return nil, fmt.Errorf("new IAM server: %w", err)
	}
	s := &IAMServer{
		client:         client,
		config:         config,
		roles:          rolesRegistry,
		memberResolver: memberResolver,
	}
	return s, nil
}

// TestPermissionOnResource tests if the caller has the specified permission on the specified resource.
func (s *IAMServer) TestPermissionOnResource(
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

func (s *IAMServer) TestResourcePermission(
	ctx context.Context,
	members []string,
	resource string,
	permission string,
) (bool, error) {
	result, err := s.TestResourcePermissions(ctx, members, map[string]string{resource: permission})
	if err != nil {
		return false, err
	}
	return result[resource], nil
}

func (s *IAMServer) TestResourcePermissions(
	ctx context.Context,
	members []string,
	resourcePermissions map[string]string,
) (map[string]bool, error) {
	result := make(map[string]bool, len(resourcePermissions))
	tx := s.client.Single()
	defer tx.Close()
	resources := make([]string, 0, len(resourcePermissions))
	for resource := range resourcePermissions {
		resources = append(resources, resource)
	}
	if err := s.ReadBindingsByResourcesAndMembersInTransaction(
		ctx,
		tx,
		resources,
		members,
		func(ctx context.Context, boundResource string, role *admin.Role, _ string) error {
			for resource, permission := range resourcePermissions {
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
	for resource := range resourcePermissions {
		if _, ok := result[resource]; !ok {
			result[resource] = false
		}
	}
	return result, nil
}

// TestPermissionOnResources tests if the caller has the specified permission on the specified resources.
func (s *IAMServer) TestPermissionOnResources(
	ctx context.Context,
	permission string,
	resources []string,
) (map[string]bool, error) {
	memberResolveResult, err := s.memberResolver.ResolveIAMMembers(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[string]bool, len(resources))
	tx := s.client.Single()
	defer tx.Close()
	if err := s.ReadBindingsByResourcesAndMembersInTransaction(
		ctx,
		tx,
		resources,
		memberResolveResult.Members(),
		func(ctx context.Context, boundResource string, role *admin.Role, _ string) error {
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

func (s *IAMServer) logError(ctx context.Context, err error) {
	if s.config.ErrorHook != nil {
		s.config.ErrorHook(ctx, err)
	}
}

func (s *IAMServer) handleStorageError(ctx context.Context, err error) error {
	s.logError(ctx, err)
	switch code := status.Code(err); code {
	case codes.Aborted, codes.Canceled, codes.DeadlineExceeded, codes.Unavailable:
		return status.Error(code, "transient storage error")
	default:
		return status.Error(codes.Internal, "storage error")
	}
}

func computeETag(policy *iam.Policy) ([]byte, error) {
	data, err := proto.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("compute etag: %w", err)
	}
	return []byte(fmt.Sprintf("W/%d-%08X", len(data), crc32.ChecksumIEEE(data))), nil
}
