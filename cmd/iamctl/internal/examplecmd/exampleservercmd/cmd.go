package exampleservercmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.einride.tech/iam/iamexample/iamexampledata"
	"go.einride.tech/iam/iammember"
	"google.golang.org/grpc"
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
	conn, err := grpc.DialContext(ctx, cfg.SpannerEmulatorHost, grpc.WithInsecure(), grpc.WithBlock())
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
		iammember.WithResolvedContext(
			ctx,
			iammember.ResolveResult{
				Metadata: iammember.Metadata{
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
