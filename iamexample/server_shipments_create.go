package iamexample

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourceid"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateShipment(
	ctx context.Context,
	request *iamexamplev1.CreateShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	var parsedRequest createShipmentRequest
	if err := parsedRequest.ParseRequest(request); err != nil {
		return nil, err
	}
	return s.createShipment(ctx, &parsedRequest)
}

func (s *Server) createShipment(
	ctx context.Context,
	request *createShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	mutations := make([]*spanner.Mutation, 0, 1+len(request.shipment.GetLineItems()))
	row, err := convertShipmentProtoToRow(request.shipment)
	if err != nil {
		s.errorHook(ctx, err)
		return nil, status.Error(codes.Internal, "internal data conversion error")
	}
	row.CreateTime = spanner.CommitTimestamp
	row.UpdateTime = spanner.CommitTimestamp
	mutations = append(mutations, spanner.Insert(row.Mutate()))
	for _, lineItemsRow := range row.LineItems {
		mutations = append(mutations, spanner.Insert(lineItemsRow.Mutate()))
	}
	commitTime, err := s.Spanner.Apply(ctx, mutations)
	if err != nil {
		switch code := status.Code(err); code {
		case codes.AlreadyExists:
			return nil, status.Errorf(code, "shipment %s already exists", request.shipment.GetName())
		default:
			return nil, s.handleStorageError(ctx, err)
		}
	}
	request.shipment.CreateTime = timestamppb.New(commitTime)
	request.shipment.UpdateTime = request.shipment.GetCreateTime()
	return request.shipment, nil
}

type createShipmentRequest struct {
	shipperID  string
	shipmentID string
	shipment   *iamexamplev1.Shipment
}

func (r *createShipmentRequest) ParseRequest(request *iamexamplev1.CreateShipmentRequest) error {
	var v validation.MessageValidator
	// parent = 1
	if request.GetParent() == "" {
		v.AddFieldViolation("parent", "required field")
	} else if resourcename.ContainsWildcard(request.GetParent()) {
		v.AddFieldViolation("parent", "must not contain wildcards")
	} else if err := resourcename.Sscan(request.GetParent(), "shippers/{shipper}", &r.shipperID); err != nil {
		v.AddFieldViolation("parent", "invalid format")
	}
	// shipment_id = 3
	if request.GetShipmentId() != "" {
		if err := resourceid.ValidateUserSettable(request.GetShipmentId()); err != nil {
			v.AddFieldError("shipment_id", err)
		}
		r.shipmentID = request.GetShipmentId()
	} else {
		r.shipmentID = resourceid.NewSystemGeneratedBase32()
	}
	if request.GetShipment() == nil {
		v.AddFieldViolation("shipment", "required field")
	} else {
		// name = 1
		request.Shipment.Name = resourcename.Sprint(
			"shippers/{shipper}/shipments/{shipment}",
			r.shipperID,
			r.shipmentID,
		)
		// create_time = 2
		request.Shipment.CreateTime = nil
		// create_time = 3
		request.Shipment.CreateTime = nil
		// delete_time = 4
		request.Shipment.DeleteTime = nil
		// origin_site = 5
		switch {
		case len(request.GetShipment().GetOriginSite()) == 0:
			v.AddFieldViolation("shipment.origin_site", "required_field")
		case !resourcename.Match("shippers/{shipper}/sites/{site}", request.GetShipment().GetOriginSite()):
			v.AddFieldViolation("shipment.origin_site", "invalid format")
		case !resourcename.HasParent(request.GetShipment().GetOriginSite(), request.GetParent()):
			v.AddFieldViolation("shipment.origin_site", "must have same parent as shipment")
		}
		// destination_site = 6
		switch {
		case len(request.GetShipment().GetDestinationSite()) == 0:
			v.AddFieldViolation("shipment.destination_site", "required field")
		case !resourcename.Match("shippers/{shipper}/sites/{site}", request.GetShipment().GetDestinationSite()):
			v.AddFieldViolation("shipment.destination_site", "invalid format")
		case !resourcename.HasParent(request.GetShipment().GetDestinationSite(), request.GetParent()):
			v.AddFieldViolation("shipment.destination_site", "must have same parent as shipment")
		}
		// pickup_earliest_time = 7
		if request.GetShipment().GetPickupEarliestTime() == nil {
			v.AddFieldViolation("shipment.pickup_earliest_time", "required field")
		} else if err := request.GetShipment().GetPickupEarliestTime().CheckValid(); err != nil {
			v.AddFieldError("shipment.pickup_earliest_time", err)
		}
		// pickup_latest_time = 8
		if request.GetShipment().GetPickupLatestTime() == nil {
			v.AddFieldViolation("shipment.pickup_latest_time", "required field")
		} else if err := request.GetShipment().GetPickupLatestTime().CheckValid(); err != nil {
			v.AddFieldError("shipment.pickup_latest_time", err)
		}
		// delivery_earliest_time = 9
		if request.GetShipment().GetDeliveryEarliestTime() == nil {
			v.AddFieldViolation("shipment.delivery_earliest_time", "required field")
		} else if err := request.GetShipment().GetDeliveryEarliestTime().CheckValid(); err != nil {
			v.AddFieldError("shipment.delivery_earliest_time", err)
		}
		// delivery_latest_time = 10
		if request.GetShipment().GetDeliveryLatestTime() == nil {
			v.AddFieldViolation("shipment.delivery_latest_time", "required field")
		} else if err := request.GetShipment().GetDeliveryLatestTime().CheckValid(); err != nil {
			v.AddFieldError("shipment.delivery_latest_time", err)
		}
		// line_items = 11
		for i, lineItem := range request.GetShipment().GetLineItems() {
			if lineItem.GetTitle() == "" {
				v.AddFieldViolation(fmt.Sprintf("shipment.line_items[%d].title", i), "required field")
			}
			if lineItem.GetQuantity() == 0 {
				v.AddFieldViolation(fmt.Sprintf("shipment.line_items[%d].quantity", i), "required field")
			}
		}
		// annotations = 12
		for key, value := range request.GetShipment().GetAnnotations() {
			if key == "" {
				v.AddFieldViolation(fmt.Sprintf(`shipment.annotations["%s"]`, key), "key must be non-empty")
			}
			if value == "" {
				v.AddFieldViolation(fmt.Sprintf(`shipment.annotations["%s"]`, key), "value must be non-empty")
			}
		}
		r.shipment = request.GetShipment()
	}
	return v.Err()
}
