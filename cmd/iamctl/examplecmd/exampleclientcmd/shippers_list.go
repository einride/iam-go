package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var listShippersCommand = &cobra.Command{
	Use:   "list-shippers",
	Short: "List shippers",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var cfg listShipperCommandConfig
		if err := viperCfg.Unmarshal(&cfg); err != nil {
			return err
		}
		return runListShipperCommand(cmd.Context(), &cfg)
	},
}

type listShipperCommandConfig struct {
	commandConfig `mapstructure:",squash"`
	PageSize      int32  `mapstructure:"page-size"`
	PageToken     string `mapstructure:"page-token"`
}

func init() {
	listShippersCommand.Flags().Int32("page-size", 0, "page size")
	listShippersCommand.Flags().String("page-token", "", "page token")
}

func runListShipperCommand(ctx context.Context, config *listShipperCommandConfig) error {
	client, err := config.connect(ctx)
	if err != nil {
		return err
	}
	response, err := client.ListShippers(ctx, &iamexamplev1.ListShippersRequest{
		PageSize:  config.PageSize,
		PageToken: config.PageToken,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(response))
	return nil
}
