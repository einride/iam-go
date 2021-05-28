package iammember

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// FromResolvedContext returns the resolved IAM members and metadata from the provided context.
func FromResolvedContext(ctx context.Context) ([]string, Metadata, bool) {
	value, ok := ctx.Value(resolvedContextKey{}).(resolvedContextValue)
	if !ok {
		return nil, nil, false
	}
	return value.members, value.metadata, true
}

// WithResolvedContext returns a new context with cached resolved IAM members and metadata.
func WithResolvedContext(ctx context.Context, members []string, memberMetadata Metadata) context.Context {
	return context.WithValue(
		ctx, resolvedContextKey{}, resolvedContextValue{members: members, metadata: memberMetadata},
	)
}

type resolvedContextKey struct{}

type resolvedContextValue struct {
	members  []string
	metadata Metadata
}

// FromContextResolver returns a Resolver that resolves cached IAM members and metadata from the current context.
func FromContextResolver() Resolver {
	return contextResolver{}
}

type contextResolver struct{}

// ResolveIAMMembers implements Resolver.
func (contextResolver) ResolveIAMMembers(ctx context.Context) ([]string, Metadata, error) {
	members, memberMetadata, ok := FromResolvedContext(ctx)
	if !ok {
		return nil, nil, fmt.Errorf("unresolved IAM member context")
	}
	return members, memberMetadata, nil
}

// ResolveContextUnaryInterceptor returns a gRPC server middleware that resolves IAM members with the provided resolver.
func ResolveContextUnaryInterceptor(resolver Resolver) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		request interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		members, memberMetadata, err := resolver.ResolveIAMMembers(ctx)
		if err != nil {
			return nil, err
		}
		return handler(WithResolvedContext(ctx, members, memberMetadata), request)
	}
}
