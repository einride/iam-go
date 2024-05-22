package main

import (
	"context"
	"flag"
	"log"

	"go.einride.tech/iam/iamcaller"
	"go.einride.tech/iam/iamexample/iamexampledata"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := runMain(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func runMain(ctx context.Context) error {
	address := flag.String("address", "localhost:8080", "listen to address")
	spannerEmulatorHost := flag.String("spanner-emulator-host", "localhost:9010", "connect to emulator host")
	log.Printf("connecting to Spanner emulator on address %s...", *spannerEmulatorHost)
	conn, err := grpc.NewClient(
		*spannerEmulatorHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	log.Printf("initializing Spanner database...")
	spannerClient, err := newSpannerDatabase(ctx, conn)
	if err != nil {
		return err
	}
	if err := iamexampledata.BootstrapRootAdmin(ctx, spannerClient); err != nil {
		return err
	}
	server, err := newServer(spannerClient)
	if err != nil {
		return err
	}
	log.Println("initializing example resources...")
	if err := iamexampledata.InitializeResources(
		iamcaller.WithResolvedContext(
			ctx,
			&iamv1.Caller{
				Members: []string{iamexampledata.RootAdminMember},
				Metadata: map[string]*iamv1.Caller_Metadata{
					"x-example-data-init": {
						Members: []string{iamexampledata.RootAdminMember},
					},
				},
			},
		),
		server,
	); err != nil {
		return err
	}
	return runServer(ctx, server, *address)
}
