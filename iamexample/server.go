package iamexample

import (
	"context"

	"cloud.google.com/go/spanner"
	"go.einride.tech/iam/iamspanner"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements iamexamplev1.FreightServiceServer.
type Server struct {
	IAM     *iamspanner.Server
	Spanner *spanner.Client
	Config  Config
}

// Config for a Server.
type Config struct {
	ErrorHook func(context.Context, error)
}

var _ iamexamplev1.FreightServiceServer = &Server{}

func (s *Server) errorHook(ctx context.Context, err error) {
	if s.Config.ErrorHook != nil {
		s.Config.ErrorHook(ctx, err)
	}
}

func (s *Server) handleStorageError(ctx context.Context, err error) error {
	s.errorHook(ctx, err)
	switch code := status.Code(err); code {
	case codes.Canceled, codes.DeadlineExceeded, codes.Aborted, codes.Unavailable:
		return status.Error(code, "transient storage error")
	default:
		return status.Error(codes.Internal, "internal storage error")
	}
}
