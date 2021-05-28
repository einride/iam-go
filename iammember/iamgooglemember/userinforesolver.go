package iamgooglemember

import "context"

// UserInfoResolver resolves IAM members from Google ID token UserInfo.
type UserInfoResolver interface {
	ResolveIAMMembersFromGoogleUserInfo(context.Context, *UserInfo) ([]string, error)
}
