package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var deleteSiteCommand = &cobra.Command{
	Use:   "delete-site",
	Short: "Delete a site",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var cfg deleteSiteCommandConfig
		if err := viperCfg.Unmarshal(&cfg); err != nil {
			return err
		}
		return runDeleteSiteCommand(cmd.Context(), &cfg)
	},
}

type deleteSiteCommandConfig struct {
	commandConfig `mapstructure:",squash"`
	Name          string `mapstructure:"name"`
}

func init() {
	deleteSiteCommand.Flags().String("name", "", "resource name of the site")
	_ = deleteSiteCommand.MarkFlagRequired("name")
}

func runDeleteSiteCommand(ctx context.Context, config *deleteSiteCommandConfig) error {
	client, err := config.connect(ctx)
	if err != nil {
		return err
	}
	site, err := client.DeleteSite(ctx, &iamexamplev1.DeleteSiteRequest{
		Name: config.Name,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(site))
	return nil
}
