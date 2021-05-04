package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var createShipperCommand = &cobra.Command{
	Use:   "create-shipper",
	Short: "Create a shipper",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var cfg createShipperCommandConfig
		if err := viperCfg.Unmarshal(&cfg); err != nil {
			return err
		}
		return runCreateShipperCommand(cmd.Context(), &cfg)
	},
}

type createShipperCommandConfig struct {
	commandConfig `mapstructure:",squash"`
	DisplayName   string `mapstructure:"display-name"`
	ShipperID     string `mapstructure:"shipper-id"`
}

func init() {
	createShipperCommand.Flags().String("shipper-id", "", "ID to use for the Shipper")
	createShipperCommand.Flags().String("display-name", "", "The display name of the shipper")
	_ = createShipperCommand.MarkFlagRequired("display-name")
}

func runCreateShipperCommand(ctx context.Context, config *createShipperCommandConfig) error {
	client, err := config.connect(ctx)
	if err != nil {
		return err
	}
	shipper, err := client.CreateShipper(ctx, &iamexamplev1.CreateShipperRequest{
		Shipper: &iamexamplev1.Shipper{
			DisplayName: config.DisplayName,
		},
		ShipperId: config.ShipperID,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(shipper))
	return nil
}
