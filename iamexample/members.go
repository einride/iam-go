package iamexample

import (
	"context"

	"go.einride.tech/iam/iammember"
	"google.golang.org/grpc/metadata"
)

// MemberHeader is the gRPC header used by the example server to determine IAM members of the caller.
const MemberHeader = "x-iam-example-members"

// NewIAMMemberHeaderResolver returns an iammember.Resolver that resolves members from MemberHeader.
func NewIAMMemberHeaderResolver() iammember.Resolver {
	return &iamMemberHeaderResolver{}
}

var _ iammember.Resolver = &iamMemberHeaderResolver{}

type iamMemberHeaderResolver struct{}

// ResolveIAMMembers implements iammember.Resolver.
func (m *iamMemberHeaderResolver) ResolveIAMMembers(ctx context.Context) (iammember.ResolveResult, error) {
	var result iammember.ResolveResult
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return result, nil
	}
	result.Add(MemberHeader, iammember.MetadataValue{Members: md.Get(MemberHeader)})
	return result, nil
}

// WithOutgoingMembers appends the provided members to the outgoing gRPC context.
func WithOutgoingMembers(ctx context.Context, members ...string) context.Context {
	pairs := make([]string, 0, len(members)*2)
	for _, member := range members {
		pairs = append(pairs, MemberHeader, member)
	}
	return metadata.AppendToOutgoingContext(ctx, pairs...)
}
