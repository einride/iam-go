package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var searchSitesCommand = &cobra.Command{
	Use:   "search-sites",
	Short: "Search sites",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var cfg searchSitesCommandConfig
		if err := viperCfg.Unmarshal(&cfg); err != nil {
			return err
		}
		return runSearchSitesCommand(cmd.Context(), &cfg)
	},
}

type searchSitesCommandConfig struct {
	commandConfig `mapstructure:",squash"`
	Parent        string `mapstructure:"parent"`
	PageSize      int32  `mapstructure:"page-size"`
	PageToken     string `mapstructure:"page-token"`
}

func init() {
	searchSitesCommand.Flags().String("parent", "", "parent shipper")
	searchSitesCommand.Flags().Int32("page-size", 0, "page size")
	searchSitesCommand.Flags().String("page-token", "", "page token")
	_ = searchSitesCommand.MarkFlagRequired("parent")
}

func runSearchSitesCommand(ctx context.Context, config *searchSitesCommandConfig) error {
	client, err := config.connect(ctx)
	if err != nil {
		return err
	}
	response, err := client.SearchSites(ctx, &iamexamplev1.SearchSitesRequest{
		Parent:    config.Parent,
		PageSize:  config.PageSize,
		PageToken: config.PageToken,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(response))
	return nil
}
