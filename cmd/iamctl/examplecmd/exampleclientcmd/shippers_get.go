package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var getShipperCommand = &cobra.Command{
	Use:   "get-shipper",
	Short: "Get a shipper",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var cfg getShipperCommandConfig
		if err := viperCfg.Unmarshal(&cfg); err != nil {
			return err
		}
		return runGetShipperCommand(cmd.Context(), &cfg)
	},
}

type getShipperCommandConfig struct {
	commandConfig `mapstructure:",squash"`
	Name          string `mapstructure:"name"`
}

func init() {
	getShipperCommand.Flags().String("name", "", "resource name of the shipper")
	_ = getShipperCommand.MarkFlagRequired("name")
}

func runGetShipperCommand(ctx context.Context, config *getShipperCommandConfig) error {
	client, err := config.connect(ctx)
	if err != nil {
		return err
	}
	shipper, err := client.GetShipper(ctx, &iamexamplev1.GetShipperRequest{
		Name: config.Name,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(shipper))
	return nil
}
