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
		var flags deleteShipmentFlags
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
		return runDeleteShipmentCommand(cmd.Context(), client, &flags)
	},
}

type deleteShipmentFlags struct {
	connection.Flags `mapstructure:",squash"`
	Name             string `mapstructure:"name"`
}

func init() {
	deleteShipmentCommand.Flags().String("name", "", "resource name of the shipment")
	_ = deleteShipmentCommand.MarkFlagRequired("name")
}

func runDeleteShipmentCommand(
	ctx context.Context,
	client iamexamplev1.FreightServiceClient,
	flags *deleteShipmentFlags,
) error {
	shipment, err := client.DeleteShipment(ctx, &iamexamplev1.DeleteShipmentRequest{
		Name: flags.Name,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(shipment))
	return nil
}
