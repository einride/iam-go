package iamgooglemember

import (
	"context"
	"encoding/base64"
	"fmt"

	"go.einride.tech/iam/iammember"
	"google.golang.org/grpc/metadata"
)

// Known UserInfo headers.
const (
	GoogleCloudEndpointUserInfoHeader   = "x-endpoint-api-userinfo"
	GoogleCloudAPIGatewayUserInfoHeader = "x-apigateway-api-userinfo"
)

// ResolveUserInfoHeader returns an iammember.Resolver that uses the provided UserInfoResolver
// to resolve IAM members from a UserInfo key.
func ResolveUserInfoHeader(key string, userInfoResolver UserInfoResolver) iammember.Resolver {
	return userInfoHeaderResolver{
		key:              key,
		userInfoResolver: userInfoResolver,
	}
}

type userInfoHeaderResolver struct {
	key              string
	userInfoResolver UserInfoResolver
}

// ResolveIAMMembers implements iammember.Resolver.
func (u userInfoHeaderResolver) ResolveIAMMembers(ctx context.Context) ([]string, iammember.Metadata, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, iammember.Metadata{u.key: nil}, nil
	}
	if !ok {
		return nil, iammember.Metadata{u.key: nil}, nil
	}
	values := md.Get(u.key)
	if len(values) == 0 {
		return nil, iammember.Metadata{u.key: nil}, nil
	}
	var userInfo UserInfo
	if err := userInfo.UnmarshalBase64(values[0], base64.URLEncoding); err != nil {
		return nil, nil, fmt.Errorf("resolve IAM members from '%s' key: %w", u.key, err)
	}
	members, err := u.userInfoResolver.ResolveIAMMembersFromGoogleUserInfo(ctx, &userInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve IAM members from '%s' key: %w", u.key, err)
	}
	return members, iammember.Metadata{u.key: members}, nil
}
