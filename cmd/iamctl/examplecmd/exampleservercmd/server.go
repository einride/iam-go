package exampleservercmd

import (
	"context"
	"log"
	"net"

	"cloud.google.com/go/spanner"
	"go.einride.tech/iam/iamexample"
	"go.einride.tech/iam/iamexample/iamexampledata"
	"go.einride.tech/iam/iamregistry"
	"go.einride.tech/iam/iamspanner"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/genproto/googleapis/iam/admin/v1"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/grpc"
)

func newServer(spannerClient *spanner.Client) (iamexamplev1.FreightServiceServer, error) {
	roles, err := iamregistry.NewRoles(iamexampledata.PredefinedRoles())
	if err != nil {
		return nil, err
	}
	iamServer, err := iamspanner.NewServer(
		spannerClient,
		roles,
		iamexample.NewIAMMemberHeaderResolver(),
		iamspanner.ServerConfig{
			ErrorHook: func(ctx context.Context, err error) {
				log.Println(err)
			},
		},
	)
	if err != nil {
		return nil, err
	}
	freightServer := &iamexample.Server{
		IAM:     iamServer,
		Spanner: spannerClient,
		Config: iamexample.Config{
			ErrorHook: func(ctx context.Context, err error) {
				log.Println(err)
			},
		},
	}
	freightServerAuthorization := &iamexample.Authorization{
		Next: freightServer,
		IAM:  iamServer,
	}
	return freightServerAuthorization, nil
}

func runServer(ctx context.Context, server iamexamplev1.FreightServiceServer, address string) error {
	grpcServer := grpc.NewServer()
	iam.RegisterIAMPolicyServer(grpcServer, server)
	if adminServer, ok := server.(admin.IAMServer); ok {
		admin.RegisterIAMServer(grpcServer, adminServer)
	}
	iamexamplev1.RegisterFreightServiceServer(grpcServer, server)
	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	log.Printf("example server listening on %s", address)
	if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		return err
	}
	return nil
}
