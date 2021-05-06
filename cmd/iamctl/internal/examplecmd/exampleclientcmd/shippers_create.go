package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.einride.tech/iam/cmd/iamctl/internal/connection"
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
		var flags createShipperFlags
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
		return runCreateShipperCommand(cmd.Context(), client, &flags)
	},
}

type createShipperFlags struct {
	connection.Flags `mapstructure:",squash"`
	DisplayName      string `mapstructure:"display-name"`
	ShipperID        string `mapstructure:"shipper-id"`
}

func init() {
	createShipperCommand.Flags().String("shipper-id", "", "ID to use for the Shipper")
	createShipperCommand.Flags().String("display-name", "", "The display name of the shipper")
	_ = createShipperCommand.MarkFlagRequired("display-name")
}

func runCreateShipperCommand(
	ctx context.Context,
	client iamexamplev1.FreightServiceClient,
	flags *createShipperFlags,
) error {
	shipper, err := client.CreateShipper(ctx, &iamexamplev1.CreateShipperRequest{
		Shipper: &iamexamplev1.Shipper{
			DisplayName: flags.DisplayName,
		},
		ShipperId: flags.ShipperID,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(shipper))
	return nil
}
