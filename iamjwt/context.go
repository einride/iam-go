package iamjwt

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

// FromIncomingContext looks up a JWT token in the incoming gRPC request metadata by key.
func FromIncomingContext(ctx context.Context, key string) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}
	values := md.Get(key)
	if len(values) == 0 {
		return "", false
	}
	value := values[0]
	const prefix = "bearer "
	isBearerToken := len(value) > len(prefix) && strings.EqualFold(value[:len(prefix)], prefix)
	if !isBearerToken {
		return "", false
	}
	return value[len(prefix):], true
}
