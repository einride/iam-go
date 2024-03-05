package iamexample

import (
	"context"

	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetShipment implements iamexamplev1.FreightServiceServer.
func (s *Server) GetShipment(
	ctx context.Context,
	request *iamexamplev1.GetShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	var parsedRequest getShipmentRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.getShipment(ctx, &parsedRequest)
}

func (s *Server) getShipment(
	ctx context.Context,
	request *getShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	tx := s.Spanner.Single()
	defer tx.Close()
	row, err := iamexampledb.Query(tx).GetShipmentsRow(ctx, iamexampledb.GetShipmentsRowQuery{
		Key: iamexampledb.ShipmentsKey{
			ShipperId:  request.shipperID,
			ShipmentId: request.shipmentID,
		},
	})
	if err != nil {
		if code := status.Code(err); code == codes.NotFound {
			return nil, status.Errorf(code, "not found: %s", request.name)
		}
		return nil, s.handleStorageError(ctx, err)
	}
	msg, err := convertShipmentRowToProto(row)
	if err != nil {
		s.errorHook(ctx, err)
		return nil, status.Error(codes.Internal, "internal data conversion error")
	}
	return msg, nil
}

type getShipmentRequest struct {
	name       string
	shipperID  string
	shipmentID string
}

func (r *getShipmentRequest) parse(request *iamexamplev1.GetShipmentRequest) error {
	var v validation.MessageValidator
	// name = 1
	if request.GetName() == "" {
		v.AddFieldViolation("name", "required field")
	} else if err := resourcename.Sscan(
		request.GetName(),
		"shippers/{shipper}/shipments/{shipment}",
		&r.shipperID,
		&r.shipmentID,
	); err != nil {
		v.AddFieldError("name", err)
	}
	r.name = request.GetName()
	return v.Err()
}
