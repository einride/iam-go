package iamexample

import (
	"fmt"
	"strings"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertShipmentProtoToRow(msg *iamexamplev1.Shipment) (*iamexampledb.ShipmentsRow, error) {
	var row iamexampledb.ShipmentsRow
	// name = 1
	if err := resourcename.Sscan(
		msg.GetName(),
		"shippers/{shipper}/shipments/{shipment}",
		&row.ShipperId,
		&row.ShipmentId,
	); err != nil {
		return nil, fmt.Errorf("name: %w", err)
	}
	// create_time = 2
	row.CreateTime = msg.GetCreateTime().AsTime()
	// update_time = 3
	row.UpdateTime = msg.GetUpdateTime().AsTime()
	// delete_time = 4
	if msg.GetDeleteTime() != nil {
		row.DeleteTime = spanner.NullTime{
			Time:  msg.GetDeleteTime().AsTime(),
			Valid: true,
		}
	}
	// origin_site = 5
	var originSiteShipperID string
	if err := resourcename.Sscan(
		msg.OriginSite,
		"shippers/{shipper}/sites/{site}",
		&originSiteShipperID,
		&row.OriginSiteId,
	); err != nil {
		return nil, fmt.Errorf("origin_site: %w", err)
	}
	// destination_site = 6
	var destinationSiteShipperID string
	if err := resourcename.Sscan(
		msg.DestinationSite,
		"shippers/{shipper}/sites/{site}",
		&destinationSiteShipperID,
		&row.DestinationSiteId,
	); err != nil {
		return nil, fmt.Errorf("destination_site: %w", err)
	}
	// pickup_earliest_time = 7
	row.PickupEarliestTime = msg.PickupEarliestTime.AsTime()
	// pickup_latest_time = 8
	row.PickupLatestTime = msg.PickupLatestTime.AsTime()
	// delivery_earliest_time = 9
	row.DeliveryEarliestTime = msg.DeliveryEarliestTime.AsTime()
	// delivery_latest_time = 10
	row.DeliveryLatestTime = msg.DeliveryLatestTime.AsTime()
	// line_items = 11
	row.LineItems = make([]*iamexampledb.LineItemsRow, 0, len(msg.LineItems))
	for i, lineItem := range msg.LineItems {
		var lineItemRow iamexampledb.LineItemsRow
		lineItemRow.ShipperId = row.ShipperId
		lineItemRow.ShipmentId = row.ShipmentId
		lineItemRow.LineNumber = int64(i)
		// title = 1
		lineItemRow.Title = lineItem.Title
		// quantity = 2
		lineItemRow.Quantity = float64(lineItem.Quantity)
		// weight_kg = 3
		lineItemRow.WeightKg = float64(lineItem.WeightKg)
		// volume_m3 = 4
		lineItemRow.VolumeM3 = float64(lineItem.VolumeM3)
		row.LineItems = append(row.LineItems, &lineItemRow)
	}
	// annotations = 12
	row.Annotations = make([]spanner.NullString, 0, len(msg.Annotations))
	for key, value := range msg.Annotations {
		row.Annotations = append(row.Annotations, spanner.NullString{
			StringVal: fmt.Sprintf("%s=%s", key, value),
			Valid:     true,
		})
	}
	return &row, nil
}

func convertShipmentRowToProto(row *iamexampledb.ShipmentsRow) (*iamexamplev1.Shipment, error) {
	var msg iamexamplev1.Shipment
	// name = 1
	msg.Name = resourcename.Sprint("shippers/{shipper}/shipments/{shipment}", row.ShipperId, row.ShipmentId)
	// create_time = 2
	msg.CreateTime = timestamppb.New(row.CreateTime)
	// update_time = 3
	msg.UpdateTime = timestamppb.New(row.UpdateTime)
	// delete_time = 4
	if row.DeleteTime.Valid {
		msg.DeleteTime = timestamppb.New(row.DeleteTime.Time)
	}
	// origin_site = 5
	msg.OriginSite = resourcename.Sprint("shippers/{shipper}/sites/{site}", row.ShipperId, row.OriginSiteId)
	// destination_site = 6
	msg.DestinationSite = resourcename.Sprint("shippers/{shipper}/sites/{site}", row.ShipperId, row.DestinationSiteId)
	// pickup_earliest_time = 7
	msg.PickupEarliestTime = timestamppb.New(row.PickupEarliestTime)
	// pickup_latest_time = 8
	msg.PickupLatestTime = timestamppb.New(row.PickupLatestTime)
	// delivery_earliest_time = 9
	msg.DeliveryEarliestTime = timestamppb.New(row.DeliveryEarliestTime)
	// delivery_latest_time = 10
	msg.DeliveryLatestTime = timestamppb.New(row.DeliveryLatestTime)
	// line_items = 11
	msg.LineItems = make([]*iamexamplev1.LineItem, 0, len(row.LineItems))
	for _, lineItemRow := range row.LineItems {
		var lineItem iamexamplev1.LineItem
		// title = 1
		lineItem.Title = lineItemRow.Title
		// quantity = 2
		lineItem.Quantity = float32(lineItemRow.Quantity)
		// weight_kg = 3
		lineItem.WeightKg = float32(lineItemRow.WeightKg)
		// volume_m3 = 4
		lineItem.VolumeM3 = float32(lineItemRow.VolumeM3)
		msg.LineItems = append(msg.LineItems, &lineItem)
	}
	// annotations = 12
	msg.Annotations = make(map[string]string, len(row.Annotations))
	for i, annotation := range row.Annotations {
		if !annotation.Valid {
			return nil, fmt.Errorf("annotations[%d] is null", i)
		}
		indexOfEquals := strings.IndexByte(annotation.StringVal, '=')
		msg.Annotations[annotation.StringVal[:indexOfEquals]] = annotation.StringVal[indexOfEquals+1:]
	}
	return &msg, nil
}
