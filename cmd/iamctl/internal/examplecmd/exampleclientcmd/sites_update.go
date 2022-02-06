package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.einride.tech/iam/cmd/iamctl/internal/connection"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// nolint: gochecknoglobals
var updateSiteCommand = &cobra.Command{
	Use:   "update-site",
	Short: "Update a site",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var flags updateSiteFlags
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
		return runUpdateSiteCommand(cmd.Context(), client, &flags)
	},
}

type updateSiteFlags struct {
	connection.Flags `mapstructure:",squash"`
	Name             string   `mapstructure:"name"`
	DisplayName      string   `mapstructure:"display-name"`
	Latitude         float64  `mapstructure:"latitude"`
	Longitude        float64  `mapstructure:"longitude"`
	UpdateMask       []string `mapstructure:"update-mask"`
}

// nolint: gochecknoinits
func init() {
	updateSiteCommand.Flags().String("name", "", "resource name of the site")
	updateSiteCommand.Flags().String("display-name", "", "page token")
	updateSiteCommand.Flags().Float64("latitude", 0, "latitude of the site")
	updateSiteCommand.Flags().Float64("longitude", 0, "latitude of the site")
	updateSiteCommand.Flags().StringSlice("update-mask", nil, "update mask")
	_ = updateSiteCommand.MarkFlagRequired("name")
}

func runUpdateSiteCommand(
	ctx context.Context,
	client iamexamplev1.FreightServiceClient,
	config *updateSiteFlags,
) error {
	var latLng *latlng.LatLng
	if config.Latitude != 0 || config.Longitude != 0 {
		latLng = &latlng.LatLng{
			Latitude:  config.Latitude,
			Longitude: config.Longitude,
		}
	}
	var updateMask *fieldmaskpb.FieldMask
	if len(config.UpdateMask) > 0 {
		updateMask = &fieldmaskpb.FieldMask{Paths: config.UpdateMask}
	}
	response, err := client.UpdateSite(ctx, &iamexamplev1.UpdateSiteRequest{
		Site: &iamexamplev1.Site{
			Name:        config.Name,
			DisplayName: config.DisplayName,
			LatLng:      latLng,
		},
		UpdateMask: updateMask,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(response))
	return nil
}
