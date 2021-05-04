package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var listShipmentsCommand = &cobra.Command{
	Use:   "list-shipments",
	Short: "List shipments",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var cfg listShipmentCommandConfig
		if err := viperCfg.Unmarshal(&cfg); err != nil {
			return err
		}
		return runListShipmentCommand(cmd.Context(), &cfg)
	},
}

type listShipmentCommandConfig struct {
	commandConfig `mapstructure:",squash"`
	Parent        string `mapstructure:"parent"`
	PageSize      int32  `mapstructure:"page-size"`
	PageToken     string `mapstructure:"page-token"`
}

func init() {
	listShipmentsCommand.Flags().String("parent", "", "parent shipper")
	listShipmentsCommand.Flags().Int32("page-size", 0, "page size")
	listShipmentsCommand.Flags().String("page-token", "", "page token")
	_ = listShipmentsCommand.MarkFlagRequired("parent")
}

func runListShipmentCommand(ctx context.Context, config *listShipmentCommandConfig) error {
	client, err := config.connect(ctx)
	if err != nil {
		return err
	}
	response, err := client.ListShipments(ctx, &iamexamplev1.ListShipmentsRequest{
		Parent:    config.Parent,
		PageSize:  config.PageSize,
		PageToken: config.PageToken,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(response))
	return nil
}
