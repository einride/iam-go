package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var deleteShipperCommand = &cobra.Command{
	Use:   "delete-shipper",
	Short: "Delete a shipper",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var cfg deleteShipperCommandConfig
		if err := viperCfg.Unmarshal(&cfg); err != nil {
			return err
		}
		return runDeleteShipperCommand(cmd.Context(), &cfg)
	},
}

type deleteShipperCommandConfig struct {
	commandConfig `mapstructure:",squash"`
	Name          string `mapstructure:"name"`
}

func init() {
	deleteShipperCommand.Flags().String("name", "", "resource name of the shipper")
	_ = deleteShipperCommand.MarkFlagRequired("name")
}

func runDeleteShipperCommand(ctx context.Context, config *deleteShipperCommandConfig) error {
	client, err := config.connect(ctx)
	if err != nil {
		return err
	}
	shipper, err := client.DeleteShipper(ctx, &iamexamplev1.DeleteShipperRequest{
		Name: config.Name,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(shipper))
	return nil
}
