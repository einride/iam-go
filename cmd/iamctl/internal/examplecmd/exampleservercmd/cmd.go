package exampleservercmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.einride.tech/iam/iamcaller"
	"go.einride.tech/iam/iamexample/iamexampledata"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Command contains sub-commands for the IAM example server.
var Command = &cobra.Command{
	Use:   "server",
	Short: "IAM example server",
}

var startCommand = &cobra.Command{
	Use:   "start",
	Short: "Start example server",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		var cfg startFlags
		if err := viperCfg.Unmarshal(&cfg); err != nil {
			return err
		}
		return runStartCommand(cmd.Context(), &cfg)
	},
}

type startFlags struct {
	Address             string `mapstructure:"address"`
	SpannerEmulatorHost string `mapstructure:"spanner-emulator-host"`
}

func init() {
	Command.AddCommand(startCommand)
	startCommand.Flags().String("address", "localhost:8080", "address to listen on")
	startCommand.Flags().String("spanner-emulator-host", "localhost:9010", "emulator host to connect to")
}

func runStartCommand(ctx context.Context, cfg *startFlags) error {
	log.Printf("connecting to Spanner emulator on address %s...", cfg.SpannerEmulatorHost)
	conn, err := grpc.DialContext(
		ctx,
		cfg.SpannerEmulatorHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
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
	return runServer(ctx, server, cfg.Address)
}
