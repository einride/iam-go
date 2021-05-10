package iamauthz

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Authorize marks the current request as processed by an authorization check.
// WithAuthorization must have been called on the context for the call to be effective.
//
// Authorize should be called at the start of an authorization check, to ensure that any errors resulting from the
// authorization check itself are forwarded to the caller.
func Authorize(ctx context.Context) {
	if value, ok := ctx.Value(contextKey{}).(*contextValue); ok {
		value.mu.Lock()
		value.authorized = true
		value.mu.Unlock()
	}
}

// RequireUnaryAuthorization is a grpc.UnaryServerInterceptor that requires authorization
// to be performed on all incoming requests.
//
// To mark the request as processed by authorization checks, the method implementing authorization should call
// Authorize on the request context as soon as authorization starts.
func RequireUnaryAuthorization(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	ctx = WithAuthorization(ctx)
	resp, err := handler(ctx, req)
	if code := status.Code(err); code == codes.Unauthenticated || code == codes.PermissionDenied {
		return nil, err
	}
	value := ctx.Value(contextKey{}).(*contextValue)
	value.mu.Lock()
	authorized := value.authorized
	value.mu.Unlock()
	if !authorized {
		return nil, status.Error(codes.Internal, "server did not perform authorization")
	}
	return resp, err
}

var _ grpc.UnaryServerInterceptor = RequireUnaryAuthorization

// RequireStreamAuthorization is a grpc.StreamServerInterceptor that aborts all incoming streams, pending implementation
// of stream support in this package.
func RequireStreamAuthorization(
	_ interface{},
	_ grpc.ServerStream,
	_ *grpc.StreamServerInfo,
	_ grpc.StreamHandler,
) error {
	return status.Error(codes.Internal, "server has not implemented stream authorization")
}

var _ grpc.StreamServerInterceptor = RequireStreamAuthorization

// WithAuthorization adds authorization to the current request context.
func WithAuthorization(ctx context.Context) context.Context {
	if _, ok := ctx.Value(contextKey{}).(*contextValue); ok {
		return ctx
	}
	return context.WithValue(ctx, contextKey{}, &contextValue{})
}

type contextKey struct{}

type contextValue struct {
	mu         sync.Mutex
	authorized bool
}
