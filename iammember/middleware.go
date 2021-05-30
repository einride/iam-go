package iammember

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// FromResolvedContext returns the resolved IAM members and metadata from the provided context.
func FromResolvedContext(ctx context.Context) (ResolveResult, bool) {
	value, ok := ctx.Value(resolvedContextKey{}).(ResolveResult)
	return value, ok
}

// WithResolvedContext returns a new context with cached IAM member resolve result.
func WithResolvedContext(ctx context.Context, resolveResult ResolveResult) context.Context {
	return context.WithValue(ctx, resolvedContextKey{}, resolveResult)
}

type resolvedContextKey struct{}

// FromContextResolver returns a Resolver that resolves cached IAM members and metadata from the current context.
func FromContextResolver() Resolver {
	return contextResolver{}
}

type contextResolver struct{}

// ResolveIAMMembers implements Resolver.
func (contextResolver) ResolveIAMMembers(ctx context.Context) (ResolveResult, error) {
	result, ok := FromResolvedContext(ctx)
	if !ok {
		return ResolveResult{}, fmt.Errorf("unresolved IAM member context")
	}
	return result, nil
}

// ResolveContextUnaryInterceptor returns a gRPC server middleware that resolves IAM members with the provided resolver.
func ResolveContextUnaryInterceptor(resolver Resolver) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		request interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		result, err := resolver.ResolveIAMMembers(ctx)
		if err != nil {
			return nil, err
		}
		return handler(WithResolvedContext(ctx, result), request)
	}
}
