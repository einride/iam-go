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
		var flags deleteSiteFlags
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
		return runDeleteSiteCommand(cmd.Context(), client, &flags)
	},
}

type deleteSiteFlags struct {
	connection.Flags `mapstructure:",squash"`
	Name             string `mapstructure:"name"`
}

// nolint: gochecknoinits
func init() {
	deleteSiteCommand.Flags().String("name", "", "resource name of the site")
	_ = deleteSiteCommand.MarkFlagRequired("name")
}

func runDeleteSiteCommand(
	ctx context.Context,
	client iamexamplev1.FreightServiceClient,
	config *deleteSiteFlags,
) error {
	site, err := client.DeleteSite(ctx, &iamexamplev1.DeleteSiteRequest{
		Name: config.Name,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(site))
	return nil
}
