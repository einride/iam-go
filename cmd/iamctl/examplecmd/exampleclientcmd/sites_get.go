package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var getSiteCommand = &cobra.Command{
	Use:   "get-site",
	Short: "Get a site",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var cfg getSiteCommandConfig
		if err := viperCfg.Unmarshal(&cfg); err != nil {
			return err
		}
		return runGetSiteCommand(cmd.Context(), &cfg)
	},
}

type getSiteCommandConfig struct {
	commandConfig `mapstructure:",squash"`
	Name          string `mapstructure:"name"`
}

func init() {
	getSiteCommand.Flags().String("name", "", "resource name of the site")
	_ = getSiteCommand.MarkFlagRequired("name")
}

func runGetSiteCommand(ctx context.Context, config *getSiteCommandConfig) error {
	client, err := config.connect(ctx)
	if err != nil {
		return err
	}
	site, err := client.GetSite(ctx, &iamexamplev1.GetSiteRequest{
		Name: config.Name,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(site))
	return nil
}
