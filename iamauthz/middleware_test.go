package iamauthz

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"gotest.tools/v3/assert"
)

func TestRequireUnaryAuthorization(t *testing.T) {
	t.Run("authorized", func(t *testing.T) {
		lis, err := net.Listen("tcp", "localhost:0")
		assert.NilError(t, err)
		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(RequireAuthorizationUnaryInterceptor))
		healthpb.RegisterHealthServer(grpcServer, &authorizedHealthServer{})
		errChan := make(chan error)
		go func() {
			if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
				errChan <- err
				return
			}
			errChan <- nil
		}()
		t.Cleanup(func() {
			assert.NilError(t, <-errChan)
		})
		t.Cleanup(func() {
			grpcServer.GracefulStop()
		})
		ctx := withTestDeadline(context.Background(), t)
		conn, err := grpc.DialContext(ctx, lis.Addr().String(), grpc.WithInsecure(), grpc.WithBlock())
		assert.NilError(t, err)
		client := healthpb.NewHealthClient(conn)
		response, err := client.Check(ctx, &healthpb.HealthCheckRequest{})
		assert.NilError(t, err)
		assert.Equal(t, healthpb.HealthCheckResponse_SERVING, response.GetStatus())
	})

	t.Run("not authorized", func(t *testing.T) {
		lis, err := net.Listen("tcp", "localhost:0")
		assert.NilError(t, err)
		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(RequireAuthorizationUnaryInterceptor))
		healthpb.RegisterHealthServer(grpcServer, &healthServer{})
		errChan := make(chan error)
		go func() {
			if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
				errChan <- err
				return
			}
			errChan <- nil
		}()
		t.Cleanup(func() {
			assert.NilError(t, <-errChan)
		})
		t.Cleanup(func() {
			grpcServer.GracefulStop()
		})
		ctx := withTestDeadline(context.Background(), t)
		conn, err := grpc.DialContext(ctx, lis.Addr().String(), grpc.WithInsecure(), grpc.WithBlock())
		assert.NilError(t, err)
		client := healthpb.NewHealthClient(conn)
		response, err := client.Check(ctx, &healthpb.HealthCheckRequest{})
		assert.Assert(t, response == nil)
		assert.Equal(t, codes.Internal, status.Code(err))
	})
}

type authorizedHealthServer struct {
	healthServer
}

func (s *authorizedHealthServer) Check(
	ctx context.Context,
	request *healthpb.HealthCheckRequest,
) (*healthpb.HealthCheckResponse, error) {
	Authorize(ctx)
	return s.healthServer.Check(ctx, request)
}

func (s *authorizedHealthServer) Watch(
	request *healthpb.HealthCheckRequest,
	server healthpb.Health_WatchServer,
) error {
	Authorize(server.Context())
	return s.healthServer.Watch(request, server)
}

type healthServer struct{}

func (s *healthServer) Check(
	_ context.Context,
	_ *healthpb.HealthCheckRequest,
) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *healthServer) Watch(
	_ *healthpb.HealthCheckRequest,
	_ healthpb.Health_WatchServer,
) error {
	return nil
}

func withTestDeadline(ctx context.Context, t *testing.T) context.Context {
	deadline, ok := t.Deadline()
	if !ok {
		return ctx
	}
	ctx, cancel := context.WithDeadline(ctx, deadline)
	t.Cleanup(cancel)
	return ctx
}
