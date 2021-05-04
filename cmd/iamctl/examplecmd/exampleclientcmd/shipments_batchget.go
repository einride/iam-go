package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var batchGetShipmentsCommand = &cobra.Command{
	Use:   "batch-get-shipments",
	Short: "Batch get shipments",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var cfg batchGetShipmentsCommandConfig
		if err := viperCfg.Unmarshal(&cfg); err != nil {
			return err
		}
		return runBatchGetShipmentsCommand(cmd.Context(), &cfg)
	},
}

type batchGetShipmentsCommandConfig struct {
	commandConfig `mapstructure:",squash"`
	Parent        string   `mapstructure:"parent"`
	Names         []string `mapstructure:"names"`
}

func init() {
	batchGetShipmentsCommand.Flags().String("parent", "", "resource name of the parent shipper")
	batchGetShipmentsCommand.Flags().StringSlice("names", nil, "resource names of the shipments")
	_ = batchGetShipmentsCommand.MarkFlagRequired("names")
}

func runBatchGetShipmentsCommand(ctx context.Context, config *batchGetShipmentsCommandConfig) error {
	client, err := config.connect(ctx)
	if err != nil {
		return err
	}
	shipment, err := client.BatchGetShipments(ctx, &iamexamplev1.BatchGetShipmentsRequest{
		Parent: config.Parent,
		Names:  config.Names,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(shipment))
	return nil
}
