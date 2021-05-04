package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var listSitesCommand = &cobra.Command{
	Use:   "list-sites",
	Short: "List sites",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var cfg listSiteCommandConfig
		if err := viperCfg.Unmarshal(&cfg); err != nil {
			return err
		}
		return runListSiteCommand(cmd.Context(), &cfg)
	},
}

type listSiteCommandConfig struct {
	commandConfig `mapstructure:",squash"`
	Parent        string `mapstructure:"parent"`
	PageSize      int32  `mapstructure:"page-size"`
	PageToken     string `mapstructure:"page-token"`
}

func init() {
	listSitesCommand.Flags().String("parent", "", "parent shipper")
	listSitesCommand.Flags().Int32("page-size", 0, "page size")
	listSitesCommand.Flags().String("page-token", "", "page token")
	_ = listSitesCommand.MarkFlagRequired("parent")
}

func runListSiteCommand(ctx context.Context, config *listSiteCommandConfig) error {
	client, err := config.connect(ctx)
	if err != nil {
		return err
	}
	response, err := client.ListSites(ctx, &iamexamplev1.ListSitesRequest{
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
