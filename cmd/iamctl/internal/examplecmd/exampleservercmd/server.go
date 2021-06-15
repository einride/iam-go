package exampleservercmd

import (
	"context"
	"fmt"
	"log"
	"net"

	"cloud.google.com/go/spanner"
	firebase "firebase.google.com/go/v4"
	"go.einride.tech/iam/iamauthz"
	"go.einride.tech/iam/iamexample"
	"go.einride.tech/iam/iamfirebase"
	"go.einride.tech/iam/iamgoogle"
	"go.einride.tech/iam/iamjwt"
	"go.einride.tech/iam/iammember"
	"go.einride.tech/iam/iammixin"
	"go.einride.tech/iam/iamspanner"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/api/idtoken"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
)

func newServer(spannerClient *spanner.Client) (*iamexample.Authorization, error) {
	iamDescriptor, err := iamexamplev1.NewFreightServiceIAMDescriptor()
	if err != nil {
		return nil, err
	}
	iamServer, err := iamspanner.NewIAMServer(
		spannerClient,
		iamDescriptor.PredefinedRoles.Role,
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
	memberResolver := iammember.ChainResolvers(
		// Resolve members from the example members header.
		iamexample.NewIAMMemberHeaderResolver(),
		// Resolve members from ID tokens in the authorization header.
		googleIDTokenMemberResolver{},
	)
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

type googleIDTokenMemberResolver struct{}

func (googleIDTokenMemberResolver) ResolveIAMMembers(ctx context.Context) (iammember.ResolveResult, error) {
	const authorizationKey = "authorization"
	var result iammember.ResolveResult
	token, ok := iamjwt.FromIncomingContext(ctx, authorizationKey)
	if !ok {
		return result, nil
	}
	var jwt iamjwt.Token
	if err := jwt.UnmarshalString(token); err != nil {
		return iammember.ResolveResult{}, err
	}
	value := iammember.MetadataValue{JWT: &jwt}
	switch {
	case iamgoogle.IsGoogleIDToken(jwt):
		googlePayload, err := idtoken.Validate(ctx, token, "")
		if err != nil {
			return iammember.ResolveResult{}, err
		}
		if iamgoogle.IsEmailVerified(googlePayload) {
			if email, ok := iamgoogle.Email(googlePayload); ok {
				value.Members = append(value.Members, fmt.Sprintf("email:%s", email))
			}
		}
		if hostedDomain, ok := iamgoogle.HostedDomain(googlePayload); ok {
			value.Members = append(value.Members, fmt.Sprintf("domain:%s", hostedDomain))
		}
		result.Add(authorizationKey, iammember.MetadataValue{
			JWT:     &jwt,
			Members: value.Members,
		})
	case iamfirebase.IsFirebaseIDToken(jwt):
		app, err := firebase.NewApp(ctx, &firebase.Config{
			ProjectID: iamfirebase.ProjectID(jwt),
		})
		if err != nil {
			return iammember.ResolveResult{}, err
		}
		authClient, err := app.Auth(ctx)
		if err != nil {
			return iammember.ResolveResult{}, err
		}
		payload, err := authClient.VerifyIDToken(ctx, token)
		if err != nil {
			return iammember.ResolveResult{}, err
		}
		value.Members = append(value.Members, fmt.Sprintf("user:%s", payload.Subject))
		if payload.Firebase.Tenant != "" {
			value.Members = append(value.Members, fmt.Sprintf("tenant:%s", payload.Firebase.Tenant))
		}
	}
	log.Printf("[IAM]\t%v %v", result.Members(), result.Metadata)
	return result, nil
}
