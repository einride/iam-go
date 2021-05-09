package exampleservercmd

import (
	"context"
	"fmt"
	"log"
	"net"

	"cloud.google.com/go/spanner"
	"go.einride.tech/iam/iamexample"
	"go.einride.tech/iam/iamexample/iamexampledata"
	"go.einride.tech/iam/iammember"
	"go.einride.tech/iam/iammember/iamgooglemember"
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
		iammember.ChainResolvers(
			// Resolve members from the example members header.
			iamexample.NewIAMMemberHeaderResolver(),
			// Resolve members from the authorization header.
			iamgooglemember.ResolveAuthorizationHeader(googleUserInfoMemberResolver{}),
			// Resolve members from the Cloud Endpoint UserInfo header.
			iamgooglemember.ResolveUserInfoHeader(
				iamgooglemember.GoogleCloudEndpointUserInfoHeader,
				googleUserInfoMemberResolver{},
			),
			// Resolve members from the API Gateway UserInfo header.
			iamgooglemember.ResolveUserInfoHeader(
				iamgooglemember.GoogleCloudAPIGatewayUserInfoHeader,
				googleUserInfoMemberResolver{},
			),
		),
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

type googleUserInfoMemberResolver struct{}

func (g googleUserInfoMemberResolver) ResolveIAMMembersFromGoogleUserInfo(
	ctx context.Context,
	info *iamgooglemember.UserInfo,
) (context.Context, []string, error) {
	if info.HostedDomain != "" && info.Email != "" {
		return ctx, []string{fmt.Sprintf("%s:%s", info.HostedDomain, info.Email)}, nil
	}
	return ctx, nil, nil
}

func runServer(ctx context.Context, server iamexamplev1.FreightServiceServer, address string) error {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(logUnary))
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

func logUnary(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("%s\n->\t%s\n<-\t%s", info.FullMethod, req, err)
	} else {
		log.Printf("%s\n->\t%s\n<-\t%s", info.FullMethod, req, resp)
	}
	return resp, err
}
