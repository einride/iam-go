package iamgooglemember

import (
	"context"
	"fmt"
	"strings"

	"go.einride.tech/iam/iammember"
	"google.golang.org/grpc/metadata"
)

// ResolveAuthorizationHeader returns an iammember.Resolver that uses the provided UserInfoResolver
// to resolve IAM members from the standard authorization metadata key.
func ResolveAuthorizationHeader(userInfoResolver UserInfoResolver) iammember.Resolver {
	return authorizationHeaderResolver{userInfoResolver: userInfoResolver}
}

type authorizationHeaderResolver struct {
	userInfoResolver UserInfoResolver
}

// ResolveIAMMembers implements iammember.Resolver.
func (a authorizationHeaderResolver) ResolveIAMMembers(ctx context.Context) ([]string, iammember.Metadata, error) {
	const key = "authorization"
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, iammember.Metadata{key: nil}, nil
	}
	values := md.Get(key)
	if len(values) == 0 {
		return nil, iammember.Metadata{key: nil}, nil
	}
	authorization := values[0]
	if !strings.HasPrefix(authorization, "bearer ") {
		return nil, iammember.Metadata{key: nil}, nil
	}
	var userInfo UserInfo
	if err := userInfo.UnmarshalAuthorization(authorization); err != nil {
		return nil, nil, fmt.Errorf("resolve IAM members from '%s' key: %w", key, err)
	}
	members, err := a.userInfoResolver.ResolveIAMMembersFromGoogleUserInfo(ctx, &userInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve IAM members from '%s' key: %w", key, err)
	}
	return members, iammember.Metadata{key: members}, nil
}
