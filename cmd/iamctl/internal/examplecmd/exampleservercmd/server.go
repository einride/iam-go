package exampleservercmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"cloud.google.com/go/spanner"
	firebase "firebase.google.com/go/v4"
	"go.einride.tech/iam/iamauthz"
	"go.einride.tech/iam/iamcaller"
	"go.einride.tech/iam/iamexample"
	"go.einride.tech/iam/iamfirebase"
	"go.einride.tech/iam/iamgoogle"
	"go.einride.tech/iam/iammixin"
	"go.einride.tech/iam/iamspanner"
	"go.einride.tech/iam/iamtoken"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/api/idtoken"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newServer(spannerClient *spanner.Client) (*iamexample.Authorization, error) {
	iamDescriptor, err := iamexamplev1.NewFreightServiceIAMDescriptor()
	if err != nil {
		return nil, err
	}
	iamServer, err := iamspanner.NewIAMServer(
		spannerClient,
		iamDescriptor.PredefinedRoles.Role,
		iamcaller.FromContextResolver(),
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
		freightServer, iamServer, iamcaller.FromContextResolver(),
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
	callerResolver := iamcaller.ChainResolvers(
		// Resolve members from the example members header.
		iamexample.NewMemberHeaderResolver(),
		// Resolve members from ID tokens in the authorization header.
		googleIdentityTokenCallerResolver{},
	)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logRequestUnaryInterceptor,
			iamcaller.ResolveContextUnaryInterceptor(callerResolver),
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
	if err := grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
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

type googleIdentityTokenCallerResolver struct{}

func (googleIdentityTokenCallerResolver) ResolveCaller(ctx context.Context) (*iamv1.Caller, error) {
	const authorizationKey = "authorization"
	var result iamv1.Caller
	if deadline, ok := ctx.Deadline(); ok {
		result.Context = &iamv1.Caller_Context{
			Deadline: timestamppb.New(deadline),
		}
	}
	token, ok := iamtoken.FromIncomingContext(ctx, authorizationKey)
	if !ok {
		return &result, nil
	}
	identityToken, err := iamtoken.ParseIdentityToken(token)
	if err != nil {
		return nil, err
	}
	if err := iamtoken.ValidateIdentityToken(identityToken); err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	metadata := iamv1.Caller_Metadata{
		IdentityToken: identityToken,
	}
	switch {
	case iamgoogle.IsGoogleIdentityToken(identityToken):
		googlePayload, err := idtoken.Validate(ctx, token, "")
		if err != nil {
			return nil, err
		}
		if iamgoogle.IsEmailVerified(googlePayload) {
			if email, ok := iamgoogle.Email(googlePayload); ok {
				metadata.Members = append(metadata.Members, fmt.Sprintf("email:%s", email))
			}
		}
		if hostedDomain, ok := iamgoogle.HostedDomain(googlePayload); ok {
			metadata.Members = append(metadata.Members, fmt.Sprintf("domain:%s", hostedDomain))
		}
		iamcaller.Add(&result, authorizationKey, &metadata)
	case iamfirebase.IsFirebaseIdentityToken(identityToken):
		app, err := firebase.NewApp(ctx, &firebase.Config{
			ProjectID: iamfirebase.ProjectID(identityToken),
		})
		if err != nil {
			return nil, err
		}
		authClient, err := app.Auth(ctx)
		if err != nil {
			return nil, err
		}
		payload, err := authClient.VerifyIDToken(ctx, token)
		if err != nil {
			return nil, err
		}
		metadata.Members = append(metadata.Members, fmt.Sprintf("user:%s", payload.Subject))
		if payload.Firebase.Tenant != "" {
			metadata.Members = append(metadata.Members, fmt.Sprintf("tenant:%s", payload.Firebase.Tenant))
		}
	}
	log.Printf("[IAM]\t%v %v", result.Members, result.Metadata)
	return &result, nil
}
