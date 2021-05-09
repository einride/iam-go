package iamgooglemember

import (
	"context"

	"go.einride.tech/iam/iammember"
	"google.golang.org/grpc/metadata"
)

// Known UserInfo headers.
const (
	GoogleCloudEndpointUserInfoHeader   = "x-endpoint-api-userinfo"
	GoogleCloudAPIGatewayUserInfoHeader = "x-apigateway-api-userinfo"
)

// ResolveUserInfoHeader returns an iammember.Resolver that uses the provided UserInfoResolver
// to resolve IAM members from a UserInfo header.
func ResolveUserInfoHeader(header string, userInfoResolver UserInfoResolver) iammember.Resolver {
	return userInfoHeaderResolver{
		header:           header,
		userInfoResolver: userInfoResolver,
	}
}

type userInfoHeaderResolver struct {
	header           string
	userInfoResolver UserInfoResolver
}

// ResolveIAMMembers implements iammember.Resolver.
func (u userInfoHeaderResolver) ResolveIAMMembers(ctx context.Context) (context.Context, []string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, nil, nil
	}
	if !ok {
		return ctx, nil, nil
	}
	values := md.Get(u.header)
	if len(values) == 0 {
		return ctx, nil, nil
	}
	var userInfo UserInfo
	if err := userInfo.UnmarshalBase64(values[0]); err != nil {
		return nil, nil, err
	}
	return u.userInfoResolver.ResolveIAMMembersFromGoogleUserInfo(ctx, &userInfo)
}
