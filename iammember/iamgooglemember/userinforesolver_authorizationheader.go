package iamgooglemember

import (
	"context"
	"strings"

	"go.einride.tech/iam/iammember"
	"google.golang.org/grpc/metadata"
)

// ResolveAuthorizationHeader returns an iammember.Resolver that uses the provided UserInfoResolver
// to resolve IAM members from the standard authorization header.
func ResolveAuthorizationHeader(userInfoResolver UserInfoResolver) iammember.Resolver {
	return authorizationHeaderResolver{userInfoResolver: userInfoResolver}
}

type authorizationHeaderResolver struct {
	userInfoResolver UserInfoResolver
}

// ResolveIAMMembers implements iammember.Resolver.
func (a authorizationHeaderResolver) ResolveIAMMembers(ctx context.Context) (context.Context, []string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, nil, nil
	}
	values := md.Get("authorization")
	if len(values) == 0 {
		return ctx, nil, nil
	}
	authorization := values[0]
	indexOfSpace := strings.IndexByte(authorization, ' ')
	if indexOfSpace == -1 {
		return ctx, nil, nil
	}
	if !strings.EqualFold(authorization[:indexOfSpace], "bearer") {
		return ctx, nil, nil
	}
	var userInfo UserInfo
	if err := userInfo.UnmarshalJWT(authorization[indexOfSpace+1:]); err != nil {
		return nil, nil, err
	}
	return a.userInfoResolver.ResolveIAMMembersFromGoogleUserInfo(ctx, &userInfo)
}
