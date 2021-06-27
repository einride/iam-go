package iamexample

import (
	"context"

	"go.einride.tech/iam/iamcaller"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MemberHeader is the gRPC header used by the example server to determine IAM members of the caller.
const MemberHeader = "x-iam-example-members"

// NewMemberHeaderResolver returns an iammember.Resolver that resolves members from MemberHeader.
func NewMemberHeaderResolver() iamcaller.Resolver {
	return &memberHeaderResolver{}
}

var _ iamcaller.Resolver = &memberHeaderResolver{}

type memberHeaderResolver struct{}

// ResolveCaller implements iamcaller.Resolver.
func (m *memberHeaderResolver) ResolveCaller(ctx context.Context) (*iamv1.Caller, error) {
	var result iamv1.Caller
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		iamcaller.Add(&result, MemberHeader, &iamv1.Caller_Metadata{
			Members: md.Get(MemberHeader),
		})
	}
	if deadline, ok := ctx.Deadline(); ok {
		result.Context = &iamv1.Caller_Context{
			Deadline: timestamppb.New(deadline),
		}
	}
	return &result, nil
}

// WithOutgoingMembers appends the provided members to the outgoing gRPC context.
func WithOutgoingMembers(ctx context.Context, members ...string) context.Context {
	pairs := make([]string, 0, len(members)*2)
	for _, member := range members {
		pairs = append(pairs, MemberHeader, member)
	}
	return metadata.AppendToOutgoingContext(ctx, pairs...)
}
