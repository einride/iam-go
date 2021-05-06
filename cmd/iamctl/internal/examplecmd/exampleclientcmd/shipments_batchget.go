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

var batchGetShipmentsCommand = &cobra.Command{
	Use:   "batch-get-shipments",
	Short: "Batch get shipments",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var flags batchGetShipmentsFlags
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
		return runBatchGetShipmentsCommand(cmd.Context(), client, &flags)
	},
}

type batchGetShipmentsFlags struct {
	connection.Flags `mapstructure:",squash"`
	Parent           string   `mapstructure:"parent"`
	Names            []string `mapstructure:"names"`
}

func init() {
	batchGetShipmentsCommand.Flags().String("parent", "", "resource name of the parent shipper")
	batchGetShipmentsCommand.Flags().StringSlice("names", nil, "resource names of the shipments")
	_ = batchGetShipmentsCommand.MarkFlagRequired("names")
}

func runBatchGetShipmentsCommand(
	ctx context.Context,
	client iamexamplev1.FreightServiceClient,
	flags *batchGetShipmentsFlags,
) error {
	shipment, err := client.BatchGetShipments(ctx, &iamexamplev1.BatchGetShipmentsRequest{
		Parent: flags.Parent,
		Names:  flags.Names,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(shipment))
	return nil
}
