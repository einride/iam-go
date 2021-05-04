package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var deleteShipmentCommand = &cobra.Command{
	Use:   "delete-shipment",
	Short: "Delete a shipment",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var cfg deleteShipmentCommandConfig
		if err := viperCfg.Unmarshal(&cfg); err != nil {
			return err
		}
		return runDeleteShipmentCommand(cmd.Context(), &cfg)
	},
}

type deleteShipmentCommandConfig struct {
	commandConfig `mapstructure:",squash"`
	Name          string `mapstructure:"name"`
}

func init() {
	deleteShipmentCommand.Flags().String("name", "", "resource name of the shipment")
	_ = deleteShipmentCommand.MarkFlagRequired("name")
}

func runDeleteShipmentCommand(ctx context.Context, config *deleteShipmentCommandConfig) error {
	client, err := config.connect(ctx)
	if err != nil {
		return err
	}
	shipment, err := client.DeleteShipment(ctx, &iamexamplev1.DeleteShipmentRequest{
		Name: config.Name,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(shipment))
	return nil
}
