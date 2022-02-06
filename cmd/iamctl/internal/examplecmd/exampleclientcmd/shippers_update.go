package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.einride.tech/iam/cmd/iamctl/internal/connection"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// nolint: gochecknoglobals
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
		var flags updateShipperFlags
		if err := viperCfg.Unmarshal(&flags); err != nil {
			return err
		}
		conn, err := flags.Connect(cmd.Context())
		if err != nil {
			return err
		}
		defer func() {
			_ = conn.Close()
		}()
		client := iamexamplev1.NewFreightServiceClient(conn)
		return runUpdateShipperCommand(cmd.Context(), client, &flags)
	},
}

type updateShipperFlags struct {
	connection.Flags `mapstructure:",squash"`
	Name             string   `mapstructure:"name"`
	DisplayName      string   `mapstructure:"display-name"`
	UpdateMask       []string `mapstructure:"update-mask"`
}

// nolint: gochecknoinits
func init() {
	updateShipperCommand.Flags().String("name", "", "resource name of the shipper")
	updateShipperCommand.Flags().String("display-name", "", "page token")
	updateShipperCommand.Flags().StringSlice("update-mask", nil, "update mask")
	_ = updateShipperCommand.MarkFlagRequired("name")
}

func runUpdateShipperCommand(
	ctx context.Context,
	client iamexamplev1.FreightServiceClient,
	flags *updateShipperFlags,
) error {
	var updateMask *fieldmaskpb.FieldMask
	if len(flags.UpdateMask) > 0 {
		updateMask = &fieldmaskpb.FieldMask{Paths: flags.UpdateMask}
	}
	response, err := client.UpdateShipper(ctx, &iamexamplev1.UpdateShipperRequest{
		Shipper: &iamexamplev1.Shipper{
			Name:        flags.Name,
			DisplayName: flags.DisplayName,
		},
		UpdateMask: updateMask,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(response))
	return nil
}
