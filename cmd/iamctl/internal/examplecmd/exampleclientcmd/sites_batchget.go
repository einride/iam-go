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

var batchGetSitesCommand = &cobra.Command{
	Use:   "batch-get-sites",
	Short: "Batch get sites",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var flags batchGetSitesFlags
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
		return runBatchGetSitesCommand(cmd.Context(), client, &flags)
	},
}

type batchGetSitesFlags struct {
	connection.Flags `mapstructure:",squash"`
	Parent           string   `mapstructure:"parent"`
	Names            []string `mapstructure:"names"`
}

func init() {
	batchGetSitesCommand.Flags().String("parent", "", "resource name of the parent shipper")
	batchGetSitesCommand.Flags().StringSlice("names", nil, "resource names of the sites")
	_ = batchGetSitesCommand.MarkFlagRequired("names")
}

func runBatchGetSitesCommand(
	ctx context.Context,
	client iamexamplev1.FreightServiceClient,
	flags *batchGetSitesFlags,
) error {
	site, err := client.BatchGetSites(ctx, &iamexamplev1.BatchGetSitesRequest{
		Parent: flags.Parent,
		Names:  flags.Names,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(site))
	return nil
}
