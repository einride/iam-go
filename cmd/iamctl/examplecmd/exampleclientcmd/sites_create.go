package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/protobuf/encoding/protojson"
)

var createSiteCommand = &cobra.Command{
	Use:   "create-site",
	Short: "Create a site",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var cfg createSiteCommandConfig
		if err := viperCfg.Unmarshal(&cfg); err != nil {
			return err
		}
		return runCreateSiteCommand(cmd.Context(), &cfg)
	},
}

type createSiteCommandConfig struct {
	commandConfig `mapstructure:",squash"`
	Parent        string  `mapstructure:"parent"`
	SiteID        string  `mapstructure:"site-id"`
	DisplayName   string  `mapstructure:"display-name"`
	Latitude      float64 `mapstructure:"latitude"`
	Longitude     float64 `mapstructure:"longitude"`
}

func init() {
	createSiteCommand.Flags().String("parent", "", "parent shipper to use for the site")
	createSiteCommand.Flags().String("site-id", "", "ID to use for the Site")
	createSiteCommand.Flags().String("display-name", "", "The display name of the site")
	createSiteCommand.Flags().Float64("latitude", 0, "Latitude of the site")
	createSiteCommand.Flags().Float64("longitude", 0, "Longitude of the site")
	_ = createSiteCommand.MarkFlagRequired("parent")
	_ = createSiteCommand.MarkFlagRequired("display-name")
}

func runCreateSiteCommand(ctx context.Context, config *createSiteCommandConfig) error {
	client, err := config.connect(ctx)
	if err != nil {
		return err
	}
	var latLng *latlng.LatLng
	if config.Latitude != 0 || config.Longitude != 0 {
		latLng = &latlng.LatLng{
			Latitude:  config.Latitude,
			Longitude: config.Longitude,
		}
	}
	site, err := client.CreateSite(ctx, &iamexamplev1.CreateSiteRequest{
		Parent: config.Parent,
		Site: &iamexamplev1.Site{
			DisplayName: config.DisplayName,
			LatLng:      latLng,
		},
		SiteId: config.SiteID,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(site))
	return nil
}
