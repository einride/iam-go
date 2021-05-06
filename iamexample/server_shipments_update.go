package iamexample

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/fieldmask"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UpdateShipment implements iamexamplev1.FreightServiceServer.
func (s *Server) UpdateShipment(
	ctx context.Context,
	request *iamexamplev1.UpdateShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	var parsedRequest updateShipmentRequest
	if err := parsedRequest.ParseRequest(request); err != nil {
		return nil, err
	}
	return s.updateShipment(ctx, &parsedRequest)
}

func (s *Server) updateShipment(
	ctx context.Context,
	request *updateShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	var result *iamexamplev1.Shipment
	commitTime, err := s.Spanner.ReadWriteTransaction(
		ctx,
		func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
			row, err := iamexampledb.Query(tx).GetShipmentsRow(ctx, iamexampledb.GetShipmentsRowQuery{
				Key: iamexampledb.ShipmentsKey{
					ShipperId:  request.shipperID,
					ShipmentId: request.shipmentID,
				},
				LineItems: true,
			})
			if err != nil {
				return err
			}
			result, err = convertShipmentRowToProto(row)
			if err != nil {
				return err
			}
			fieldmask.Update(request.updateMask, result, request.shipment)
			resultRow, err := convertShipmentProtoToRow(result)
			if err != nil {
				return err
			}
			resultRow.UpdateTime = spanner.CommitTimestamp
			mutations := make([]*spanner.Mutation, 0, 2+len(resultRow.LineItems))
			mutations = append(mutations, spanner.Update(resultRow.Mutate()))
			mutations = append(
				mutations,
				spanner.Delete(
					iamexampledb.Descriptor().LineItems().TableName(),
					resultRow.Key().SpannerKey().AsPrefix(),
				),
			)
			for _, lineItemsRow := range row.LineItems {
				mutations = append(mutations, spanner.Insert(lineItemsRow.Mutate()))
			}
			return tx.BufferWrite(mutations)
		},
	)
	if err != nil {
		switch code := status.Code(err); code {
		case codes.NotFound:
			return nil, status.Errorf(code, "no such shipment: %s", request.shipment.Name)
		default:
			return nil, s.handleStorageError(ctx, err)
		}
	}
	result.UpdateTime = timestamppb.New(commitTime)
	return result, nil
}

type updateShipmentRequest struct {
	shipperID  string
	shipmentID string
	shipment   *iamexamplev1.Shipment
	updateMask *fieldmaskpb.FieldMask
}

func (r *updateShipmentRequest) ParseRequest(request *iamexamplev1.UpdateShipmentRequest) error {
	hasNoMask := len(request.GetUpdateMask().GetPaths()) == 0
	has := func(path string) bool {
		if fieldmask.IsFullReplacement(request.GetUpdateMask()) {
			return true
		}
		for _, maskPath := range request.GetUpdateMask().GetPaths() {
			if path == maskPath {
				return true
			}
		}
		return false
	}
	var v validation.MessageValidator
	// shipment = 1
	if request.Shipment == nil {
		v.AddFieldViolation("shipment", "required field")
	} else {
		// name = 1
		if len(request.Shipment.Name) == 0 {
			v.AddFieldViolation("shipment.name", "required field")
		} else if err := resourcename.Sscan(
			request.Shipment.Name,
			"shippers/{shipper}/shipments/{shipment}",
			&r.shipperID,
			&r.shipmentID,
		); err != nil {
			v.AddFieldError("shipment.name", err)
		}
		shipper := resourcename.Sprint("shippers/{shipper}", r.shipperID)
		// create_time = 2
		request.Shipment.CreateTime = nil
		// update_time = 3
		request.Shipment.UpdateTime = nil
		// delete_time = 4
		request.Shipment.DeleteTime = nil
		// origin_site = 5
		if has("origin_site") || hasNoMask && len(request.Shipment.OriginSite) > 0 {
			if len(request.Shipment.OriginSite) == 0 {
				v.AddFieldViolation("shipment.origin_site", "required_field")
			} else if !resourcename.Match("shippers/{shipper}/sites/{site}", request.Shipment.OriginSite) {
				v.AddFieldViolation("shipment.origin_site", "invalid format")
			} else if !resourcename.HasParent(request.Shipment.OriginSite, shipper) {
				v.AddFieldViolation("shipment.origin_site", "must have same parent as shipment")
			}
		}
		// destination_site = 6
		if has("destination_site") || hasNoMask && len(request.Shipment.DestinationSite) > 0 {
			if len(request.Shipment.DestinationSite) == 0 {
				v.AddFieldViolation("shipment.destination_site", "required field")
			} else if !resourcename.Match("shippers/{shipper}/sites/{site}", request.Shipment.DestinationSite) {
				v.AddFieldViolation("shipment.destination_site", "invalid format")
			} else if !resourcename.HasParent(request.Shipment.DestinationSite, shipper) {
				v.AddFieldViolation("shipment.destination_site", "must have same parent as shipment")
			}
		}
		// pickup_earliest_time = 7
		if has("pickup_earliest_time") || hasNoMask && request.Shipment.PickupEarliestTime != nil {
			if request.Shipment.PickupEarliestTime == nil {
				v.AddFieldViolation("shipment.pickup_earliest_time", "required field")
			} else if err := request.Shipment.PickupEarliestTime.CheckValid(); err != nil {
				v.AddFieldError("shipment.pickup_earliest_time", err)
			}
		}
		// pickup_latest_time = 8
		if has("pickup_latest_time") || hasNoMask && request.Shipment.PickupLatestTime != nil {
			if request.Shipment.PickupLatestTime == nil {
				v.AddFieldViolation("shipment.pickup_latest_time", "required field")
			} else if err := request.Shipment.PickupLatestTime.CheckValid(); err != nil {
				v.AddFieldError("shipment.pickup_latest_time", err)
			}
		}
		// delivery_earliest_time = 9
		if has("delivery_earliest_time") || hasNoMask && request.Shipment.DeliveryEarliestTime != nil {
			if request.Shipment.DeliveryEarliestTime == nil {
				v.AddFieldViolation("shipment.delivery_earliest_time", "required field")
			} else if err := request.Shipment.DeliveryEarliestTime.CheckValid(); err != nil {
				v.AddFieldError("shipment.delivery_earliest_time", err)
			}
		}
		// delivery_latest_time = 10
		if has("delivery_latest_time") || hasNoMask && request.Shipment.DeliveryLatestTime != nil {
			if request.Shipment.DeliveryLatestTime == nil {
				v.AddFieldViolation("shipment.delivery_latest_time", "required field")
			} else if err := request.Shipment.DeliveryLatestTime.CheckValid(); err != nil {
				v.AddFieldError("shipment.delivery_latest_time", err)
			}
		}
		// line_items = 11
		if has("line_items") || hasNoMask && len(request.Shipment.LineItems) > 0 {
			for i, lineItem := range request.Shipment.LineItems {
				if lineItem.Title == "" {
					v.AddFieldViolation(fmt.Sprintf("shipment.line_items[%d].title", i), "required field")
				}
				if lineItem.Quantity == 0 {
					v.AddFieldViolation(fmt.Sprintf("shipment.line_items[%d].quantity", i), "required field")
				}
			}
		}
		// annotations = 12
		if has("annotations") || hasNoMask && len(request.Shipment.Annotations) > 0 {
			for key, value := range request.Shipment.Annotations {
				if key == "" {
					v.AddFieldViolation(fmt.Sprintf(`shipment.annotations["%s"]`, key), "key must be non-empty")
				}
				if value == "" {
					v.AddFieldViolation(fmt.Sprintf(`shipment.annotations["%s"]`, key), "value must be non-empty")
				}
			}
		}
		r.shipment = request.Shipment
	}
	// update_mask = 2
	if err := fieldmask.Validate(request.UpdateMask, request.Shipment); err != nil {
		v.AddFieldError("update_mask", err)
	} else {
		r.updateMask = request.UpdateMask
	}
	return v.Err()
}
