package exampleservercmd

import (
	"context"
	"fmt"
	"log"
	"net"

	"cloud.google.com/go/spanner"
	"go.einride.tech/iam/iamauthz"
	"go.einride.tech/iam/iamexample"
	"go.einride.tech/iam/iammember"
	"go.einride.tech/iam/iammember/iamgooglemember"
	"go.einride.tech/iam/iammixin"
	"go.einride.tech/iam/iamregistry"
	"go.einride.tech/iam/iamspanner"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
)

func newServer(spannerClient *spanner.Client) (*iamexample.Authorization, error) {
	iamDescriptor, err := iamexamplev1.NewFreightServiceIAMDescriptor()
	if err != nil {
		return nil, err
	}
	roles, err := iamregistry.NewRoles(iamDescriptor.PredefinedRoles)
	if err != nil {
		return nil, err
	}
	iamServer, err := iamspanner.NewIAMServer(
		spannerClient,
		roles,
		iammember.FromContextResolver(),
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
	authorization, err := iamexamplev1.NewFreightServiceAuthorization(
		freightServer, iamServer, iammember.FromContextResolver(),
	)
	if err != nil {
		return nil, err
	}
	freightServerAuthorization := &iamexample.Authorization{
		FreightServiceAuthorization: authorization,
		Next:                        freightServer,
		IAMServer:                   iamServer,
		IAMDescriptor:               iamDescriptor,
	}
	return freightServerAuthorization, nil
}

func runServer(
	ctx context.Context,
	server *iamexample.Authorization,
	address string,

) error {
	memberResolver := loggingIAMMemberResolver{
		next: iammember.ChainResolvers(
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
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logRequestUnaryInterceptor,
			iammember.ResolveContextUnaryInterceptor(memberResolver),
			iamauthz.RequireAuthorizationUnaryInterceptor,
		),
		grpc.StreamInterceptor(iamauthz.RequireAuthorizationStreamInterceptor),
	)
	iammixin.Register(grpcServer, server)
	longrunning.RegisterOperationsServer(grpcServer, server)
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

func logRequestUnaryInterceptor(
	ctx context.Context,
	request interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("%s\n[REQ]\t%s", info.FullMethod, request)
	response, err := handler(ctx, request)
	if err != nil {
		log.Printf("[ERR]\t%s", err)
	} else {
		log.Printf("[RES]\t%s", response)
	}
	return response, err
}

type googleUserInfoMemberResolver struct{}

func (g googleUserInfoMemberResolver) ResolveIAMMembersFromGoogleUserInfo(
	_ context.Context,
	info *iamgooglemember.UserInfo,
) ([]string, error) {
	members := make([]string, 0, 2)
	if info.Email != "" && info.EmailVerified {
		members = append(members, fmt.Sprintf("email:%s", info.Email))
	}
	if info.HostedDomain != "" {
		members = append(members, fmt.Sprintf("domain:%s", info.HostedDomain))
	}
	return members, nil
}

type loggingIAMMemberResolver struct {
	next iammember.Resolver
}

func (l loggingIAMMemberResolver) ResolveIAMMembers(ctx context.Context) ([]string, iammember.Metadata, error) {
	members, memberMetadata, err := l.next.ResolveIAMMembers(ctx)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("[IAM]\t%v %v", members, memberMetadata)
	return members, memberMetadata, nil
}
