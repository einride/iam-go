package iamspanner

import (
	"context"
	"fmt"
	"hash/crc32"

	"cloud.google.com/go/spanner"
	"go.einride.tech/iam/iamcaller"
	"go.einride.tech/iam/iammember"
	"go.einride.tech/iam/iamregistry"
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
	callerResolver iamcaller.Resolver
	config         ServerConfig
}

// ServerConfig configures a Spanner IAM policy server.
type ServerConfig struct {
	// ErrorHook is called when errors occur in the IAMServer.
	ErrorHook func(context.Context, error)
	// ValidateMember is a custom IAM member validator.
	// When not provided, iammember.Validate will be used.
	ValidateMember func(string) error
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
	callerResolver iamcaller.Resolver,
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
		callerResolver: callerResolver,
	}
	return s, nil
}

func (s *IAMServer) validateMember(member string) error {
	if s.config.ValidateMember != nil {
		return s.config.ValidateMember(member)
	}
	return iammember.Validate(member)
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
