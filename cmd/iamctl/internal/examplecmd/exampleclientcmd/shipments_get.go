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

// nolint: gochecknoglobals
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
		var flags getShipmentFlags
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
		return runGetShipmentCommand(cmd.Context(), client, &flags)
	},
}

type getShipmentFlags struct {
	connection.Flags `mapstructure:",squash"`
	Name             string `mapstructure:"name"`
}

// nolint: gochecknoinits
func init() {
	getShipmentCommand.Flags().String("name", "", "resource name of the shipment")
	_ = getShipmentCommand.MarkFlagRequired("name")
}

func runGetShipmentCommand(
	ctx context.Context,
	client iamexamplev1.FreightServiceClient,
	flags *getShipmentFlags,
) error {
	shipment, err := client.GetShipment(ctx, &iamexamplev1.GetShipmentRequest{
		Name: flags.Name,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(shipment))
	return nil
}
