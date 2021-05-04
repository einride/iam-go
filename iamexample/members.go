package iamexample

import (
	"context"

	"go.einride.tech/iam/iamspanner"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// MemberHeader is the gRPC header used by the example server to determine IAM members of the caller.
const MemberHeader = "x-iam-example-members"

// NewMemberHeaderResolver returns an iamspanner.MemberResolver that resolves members from MemberHeader.
func NewMemberHeaderResolver() iamspanner.MemberResolver {
	return &memberHeaderResolver{}
}

var _ iamspanner.MemberResolver = &memberHeaderResolver{}

type memberHeaderResolver struct{}

// ResolveMember implements iamspanner.MemberResolver.
func (m *memberHeaderResolver) ResolveMember(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "missing members header: %s", MemberHeader)
	}
	values := md.Get(MemberHeader)
	if len(values) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "missing members header: %s", MemberHeader)
	}
	return values[0], nil
}

// WithOutgoingMembers appends the provided members to the outgoing gRPC context.
func WithOutgoingMembers(ctx context.Context, members ...string) context.Context {
	pairs := make([]string, 0, len(members)*2)
	for _, member := range members {
		pairs = append(pairs, MemberHeader, member)
	}
	return metadata.AppendToOutgoingContext(ctx, pairs...)
}
