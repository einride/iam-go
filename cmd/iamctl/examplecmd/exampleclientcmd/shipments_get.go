package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var getShipmentCommand = &cobra.Command{
	Use:   "get-shipment",
	Short: "Get a shipment",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var cfg getShipmentCommandConfig
		if err := viperCfg.Unmarshal(&cfg); err != nil {
			return err
		}
		return runGetShipmentCommand(cmd.Context(), &cfg)
	},
}

type getShipmentCommandConfig struct {
	commandConfig `mapstructure:",squash"`
	Name          string `mapstructure:"name"`
}

func init() {
	getShipmentCommand.Flags().String("name", "", "resource name of the shipment")
	_ = getShipmentCommand.MarkFlagRequired("name")
}

func runGetShipmentCommand(ctx context.Context, config *getShipmentCommandConfig) error {
	client, err := config.connect(ctx)
	if err != nil {
		return err
	}
	shipment, err := client.GetShipment(ctx, &iamexamplev1.GetShipmentRequest{
		Name: config.Name,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(shipment))
	return nil
}
