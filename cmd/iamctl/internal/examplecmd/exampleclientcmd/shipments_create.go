package exampleclientcmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.einride.tech/iam/cmd/iamctl/internal/connection"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// nolint: gochecknoglobals
var createShipmentCommand = &cobra.Command{
	Use:   "create-shipment",
	Short: "Create a shipment",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var flags createShipmentFlags
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
		return runCreateShipmentCommand(cmd.Context(), client, &flags)
	},
}

type createShipmentFlags struct {
	connection.Flags     `mapstructure:",squash"`
	Parent               string `mapstructure:"parent"`
	ShipmentID           string `mapstructure:"shipment-id"`
	OriginSite           string `mapstructure:"origin-site"`
	DestinationSite      string `mapstructure:"destination-site"`
	PickupEarliestTime   string `mapstructure:"pickup-earliest-time"`
	PickupLatestTime     string `mapstructure:"pickup-latest-time"`
	DeliveryEarliestTime string `mapstructure:"delivery-earliest-time"`
	DeliveryLatestTime   string `mapstructure:"delivery-latest-time"`
}

// nolint: gochecknoinits
func init() {
	createShipmentCommand.Flags().String("parent", "", "parent shipper to use for the shipment")
	createShipmentCommand.Flags().String("shipment-id", "", "ID to use for the Shipment")
	createShipmentCommand.Flags().String("origin-site", "", "origin site of the shipment")
	createShipmentCommand.Flags().String("destination-site", "", "destination site of the shipment")
	createShipmentCommand.Flags().String("pickup-earliest-time", "", "earliest pickup time of the shipment")
	createShipmentCommand.Flags().String("pickup-latest-time", "", "latest pickup time of the shipment")
	createShipmentCommand.Flags().String("delivery-earliest-time", "", "earliest delivery time of the shipment")
	createShipmentCommand.Flags().String("delivery-latest-time", "", "latest delivery time of the shipment")
	_ = createShipmentCommand.MarkFlagRequired("parent")
	_ = createShipmentCommand.MarkFlagRequired("origin-site")
	_ = createShipmentCommand.MarkFlagRequired("destination-site")
	_ = createShipmentCommand.MarkFlagRequired("pickup-earliest-time")
	_ = createShipmentCommand.MarkFlagRequired("pickup-latest-time")
	_ = createShipmentCommand.MarkFlagRequired("delivery-earliest-time")
	_ = createShipmentCommand.MarkFlagRequired("delivery-latest-time")
}

func runCreateShipmentCommand(
	ctx context.Context,
	client iamexamplev1.FreightServiceClient,
	flags *createShipmentFlags,
) error {
	pickupEarliestTime, err := parseTime("pickup-earliest-time", flags.PickupEarliestTime)
	if err != nil {
		return err
	}
	pickupLatestTime, err := parseTime("pickup-latest-time", flags.PickupLatestTime)
	if err != nil {
		return err
	}
	deliveryEarliestTime, err := parseTime("delivery-earliest-time", flags.DeliveryEarliestTime)
	if err != nil {
		return err
	}
	deliveryLatestTime, err := parseTime("delivery-latest-time", flags.DeliveryLatestTime)
	if err != nil {
		return err
	}
	shipment, err := client.CreateShipment(ctx, &iamexamplev1.CreateShipmentRequest{
		Parent: flags.Parent,
		Shipment: &iamexamplev1.Shipment{
			OriginSite:           flags.OriginSite,
			DestinationSite:      flags.DestinationSite,
			PickupEarliestTime:   pickupEarliestTime,
			PickupLatestTime:     pickupLatestTime,
			DeliveryEarliestTime: deliveryEarliestTime,
			DeliveryLatestTime:   deliveryLatestTime,
		},
		ShipmentId: flags.ShipmentID,
	})
	if err != nil {
		return err
	}
	log.Println(protojson.Format(shipment))
	return nil
}

func parseTime(flagName, s string) (*timestamppb.Timestamp, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", flagName, err)
	}
	return timestamppb.New(t), nil
}
