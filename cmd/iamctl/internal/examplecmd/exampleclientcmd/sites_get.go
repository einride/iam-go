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
		var flags getSiteFlags
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
		return runGetSiteCommand(cmd.Context(), client, &flags)
	},
}

type getSiteFlags struct {
	connection.Flags `mapstructure:",squash"`
	Name             string `mapstructure:"name"`
}

func init() {
	getSiteCommand.Flags().String("name", "", "resource name of the site")
	_ = getSiteCommand.MarkFlagRequired("name")
}

func runGetSiteCommand(
	ctx context.Context,
	client iamexamplev1.FreightServiceClient,
	config *getSiteFlags,
) error {
	site, err := client.GetSite(ctx, &iamexamplev1.GetSiteRequest{
		Name: config.Name,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(site))
	return nil
}
