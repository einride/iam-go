package exampleclientcmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.einride.tech/iam/cmd/iamctl/internal/connection"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var updateShipmentCommand = &cobra.Command{
	Use:   "update-shipment",
	Short: "Update a shipment",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var flags updateShipmentFlags
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
		return runUpdateShipmentCommand(cmd.Context(), client, &flags)
	},
}

type updateShipmentFlags struct {
	connection.Flags     `mapstructure:",squash"`
	Name                 string   `mapstructure:"name"`
	OriginSite           string   `mapstructure:"origin-site"`
	DestinationSite      string   `mapstructure:"destination-site"`
	PickupEarliestTime   string   `mapstructure:"pickup-earliest-time"`
	PickupLatestTime     string   `mapstructure:"pickup-latest-time"`
	DeliveryEarliestTime string   `mapstructure:"delivery-earliest-time"`
	DeliveryLatestTime   string   `mapstructure:"delivery-latest-time"`
	UpdateMask           []string `mapstructure:"update-mask"`
}

func init() {
	updateShipmentCommand.Flags().String("name", "", "resource name of the shipment")
	updateShipmentCommand.Flags().String("origin-site", "", "origin site of the shipment")
	updateShipmentCommand.Flags().String("destination-site", "", "destination site of the shipment")
	updateShipmentCommand.Flags().String("pickup-earliest-time", "", "earliest pickup time of the shipment")
	updateShipmentCommand.Flags().String("pickup-latest-time", "", "latest pickup time of the shipment")
	updateShipmentCommand.Flags().String("delivery-earliest-time", "", "earliest delivery time of the shipment")
	updateShipmentCommand.Flags().String("delivery-latest-time", "", "latest delivery time of the shipment")
	updateShipmentCommand.Flags().StringSlice("update-mask", nil, "update mask")
	_ = updateShipmentCommand.MarkFlagRequired("name")
}

func runUpdateShipmentCommand(
	ctx context.Context,
	client iamexamplev1.FreightServiceClient,
	config *updateShipmentFlags,
) error {
	var err error
	var pickupEarliestTime *timestamppb.Timestamp
	if config.PickupEarliestTime != "" {
		pickupEarliestTime, err = parseTime("pickup-earliest-time", config.PickupEarliestTime)
		if err != nil {
			return err
		}
	}
	var pickupLatestTime *timestamppb.Timestamp
	if config.PickupLatestTime != "" {
		pickupLatestTime, err = parseTime("pickup-latest-time", config.PickupLatestTime)
		if err != nil {
			return err
		}
	}
	var deliveryEarliestTime *timestamppb.Timestamp
	if config.DeliveryEarliestTime != "" {
		deliveryEarliestTime, err = parseTime("delivery-earliest-time", config.DeliveryEarliestTime)
		if err != nil {
			return err
		}
	}
	var deliveryLatestTime *timestamppb.Timestamp
	if config.DeliveryLatestTime != "" {
		deliveryLatestTime, err = parseTime("delivery-latest-time", config.DeliveryLatestTime)
		if err != nil {
			return err
		}
	}
	var updateMask *fieldmaskpb.FieldMask
	if len(config.UpdateMask) > 0 {
		updateMask = &fieldmaskpb.FieldMask{Paths: config.UpdateMask}
	}
	response, err := client.UpdateShipment(ctx, &iamexamplev1.UpdateShipmentRequest{
		Shipment: &iamexamplev1.Shipment{
			Name:                 config.Name,
			OriginSite:           config.OriginSite,
			DestinationSite:      config.DestinationSite,
			PickupEarliestTime:   pickupEarliestTime,
			PickupLatestTime:     pickupLatestTime,
			DeliveryEarliestTime: deliveryEarliestTime,
			DeliveryLatestTime:   deliveryLatestTime,
		},
		UpdateMask: updateMask,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(response))
	return nil
}
