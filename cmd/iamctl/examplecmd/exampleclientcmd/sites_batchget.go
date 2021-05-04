package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		var cfg batchGetSitesCommandConfig
		if err := viperCfg.Unmarshal(&cfg); err != nil {
			return err
		}
		return runBatchGetSitesCommand(cmd.Context(), &cfg)
	},
}

type batchGetSitesCommandConfig struct {
	commandConfig `mapstructure:",squash"`
	Parent        string   `mapstructure:"parent"`
	Names         []string `mapstructure:"names"`
}

func init() {
	batchGetSitesCommand.Flags().String("parent", "", "resource name of the parent shipper")
	batchGetSitesCommand.Flags().StringSlice("names", nil, "resource names of the sites")
	_ = batchGetSitesCommand.MarkFlagRequired("names")
}

func runBatchGetSitesCommand(ctx context.Context, config *batchGetSitesCommandConfig) error {
	client, err := config.connect(ctx)
	if err != nil {
		return err
	}
	site, err := client.BatchGetSites(ctx, &iamexamplev1.BatchGetSitesRequest{
		Parent: config.Parent,
		Names:  config.Names,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(site))
	return nil
}
