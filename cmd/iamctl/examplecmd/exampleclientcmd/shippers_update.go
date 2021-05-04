package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var updateShipperCommand = &cobra.Command{
	Use:   "update-shipper",
	Short: "Update a shipper",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var cfg updateShipperCommandConfig
		if err := viperCfg.Unmarshal(&cfg); err != nil {
			return err
		}
		return runUpdateShipperCommand(cmd.Context(), &cfg)
	},
}

type updateShipperCommandConfig struct {
	commandConfig `mapstructure:",squash"`
	Name          string   `mapstructure:"name"`
	DisplayName   string   `mapstructure:"display-name"`
	UpdateMask    []string `mapstructure:"update-mask"`
}

func init() {
	updateShipperCommand.Flags().String("name", "", "resource name of the shipper")
	updateShipperCommand.Flags().String("display-name", "", "page token")
	updateShipperCommand.Flags().StringSlice("update-mask", nil, "update mask")
	_ = updateShipperCommand.MarkFlagRequired("name")
}

func runUpdateShipperCommand(ctx context.Context, config *updateShipperCommandConfig) error {
	client, err := config.connect(ctx)
	if err != nil {
		return err
	}
	var updateMask *fieldmaskpb.FieldMask
	if len(config.UpdateMask) > 0 {
		updateMask = &fieldmaskpb.FieldMask{Paths: config.UpdateMask}
	}
	response, err := client.UpdateShipper(ctx, &iamexamplev1.UpdateShipperRequest{
		Shipper: &iamexamplev1.Shipper{
			Name:        config.Name,
			DisplayName: config.DisplayName,
		},
		UpdateMask: updateMask,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(response))
	return nil
}
