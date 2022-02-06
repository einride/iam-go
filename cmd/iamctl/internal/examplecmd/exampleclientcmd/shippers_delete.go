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
		var flags deleteShipperFlags
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
		return runDeleteShipperCommand(cmd.Context(), client, &flags)
	},
}

type deleteShipperFlags struct {
	connection.Flags `mapstructure:",squash"`
	Name             string `mapstructure:"name"`
}

// nolint: gochecknoinits
func init() {
	deleteShipperCommand.Flags().String("name", "", "resource name of the shipper")
	_ = deleteShipperCommand.MarkFlagRequired("name")
}

func runDeleteShipperCommand(
	ctx context.Context,
	client iamexamplev1.FreightServiceClient,
	flags *deleteShipperFlags,
) error {
	operation, err := client.DeleteShipper(ctx, &iamexamplev1.DeleteShipperRequest{
		Name: flags.Name,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(operation))
	return nil
}
