package exampleservercmd

import (
	"context"
	"fmt"
	"log"
	"net"

	"cloud.google.com/go/spanner"
	"go.einride.tech/iam/iamauthz"
	"go.einride.tech/iam/iamexample"
	"go.einride.tech/iam/iamgoogle"
	"go.einride.tech/iam/iammember"
	"go.einride.tech/iam/iammixin"
	"go.einride.tech/iam/iamregistry"
	"go.einride.tech/iam/iamspanner"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
			// Resolve members from ID tokens in the authorization header.
			authorizationIDTokenMemberResolver{},
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

type authorizationIDTokenMemberResolver struct{}

func (authorizationIDTokenMemberResolver) ResolveIAMMembers(ctx context.Context) (iammember.ResolveResult, error) {
	const authorizationKey = "authorization"
	var result iammember.ResolveResult
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return result, nil
	}
	authorizationValues := md.Get(authorizationKey)
	if len(authorizationValues) == 0 {
		return result, nil
	}
	authorization := authorizationValues[0]
	var idToken iamgoogle.IDToken
	if err := idToken.UnmarshalAuthorization(authorization); err != nil {
		return iammember.ResolveResult{}, err
	}
	if err := idToken.Validate(); err != nil {
		return iammember.ResolveResult{}, err
	}
	if idToken.EmailVerified && idToken.Email != "" {
		result.Add(authorizationKey, fmt.Sprintf("email:%s", idToken.Email))
	}
	if idToken.HostedDomain != "" {
		result.Add(authorizationKey, fmt.Sprintf("domain:%s", idToken.HostedDomain))
	}
	return result, nil
}

type loggingIAMMemberResolver struct {
	next iammember.Resolver
}

func (l loggingIAMMemberResolver) ResolveIAMMembers(ctx context.Context) (iammember.ResolveResult, error) {
	result, err := l.next.ResolveIAMMembers(ctx)
	if err == nil {
		log.Printf("[IAM]\t%v %v", result.Members, result.Metadata)
	}
	return result, err
}
