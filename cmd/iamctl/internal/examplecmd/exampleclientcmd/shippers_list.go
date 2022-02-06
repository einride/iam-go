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
var listShippersCommand = &cobra.Command{
	Use:   "list-shippers",
	Short: "List shippers",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var flags listShipperFlags
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
		return runListShipperCommand(cmd.Context(), client, &flags)
	},
}

type listShipperFlags struct {
	connection.Flags `mapstructure:",squash"`
	PageSize         int32  `mapstructure:"page-size"`
	PageToken        string `mapstructure:"page-token"`
}

// nolint: gochecknoinits
func init() {
	listShippersCommand.Flags().Int32("page-size", 0, "page size")
	listShippersCommand.Flags().String("page-token", "", "page token")
}

func runListShipperCommand(
	ctx context.Context,
	client iamexamplev1.FreightServiceClient,
	flags *listShipperFlags,
) error {
	response, err := client.ListShippers(ctx, &iamexamplev1.ListShippersRequest{
		PageSize:  flags.PageSize,
		PageToken: flags.PageToken,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(response))
	return nil
}
