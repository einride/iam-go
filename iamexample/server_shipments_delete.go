package iamexample

import (
	"context"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DeleteShipment implements iamexamplev1.FreightServiceServer.
func (s *Server) DeleteShipment(
	ctx context.Context,
	request *iamexamplev1.DeleteShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	var parsedRequest deleteShipmentRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.deleteShipment(ctx, &parsedRequest)
}

func (s *Server) deleteShipment(
	ctx context.Context,
	request *deleteShipmentRequest,
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
			})
			if err != nil {
				return err
			}
			if row.DeleteTime.Valid {
				return status.Errorf(codes.FailedPrecondition, "site already deleted: %s", request.name)
			}
			row.UpdateTime = spanner.CommitTimestamp
			row.DeleteTime = spanner.NullTime{
				Time:  spanner.CommitTimestamp,
				Valid: true,
			}
			result, err = convertShipmentRowToProto(row)
			if err != nil {
				return err
			}
			return tx.BufferWrite([]*spanner.Mutation{spanner.Update(row.Mutate())})
		},
	)
	if err != nil {
		switch code := status.Code(err); code {
		case codes.FailedPrecondition:
			return nil, err
		case codes.NotFound:
			return nil, status.Errorf(code, "no such site: %s", request.name)
		default:
			return nil, s.handleStorageError(ctx, err)
		}
	}
	result.UpdateTime = timestamppb.New(commitTime)
	result.DeleteTime = result.UpdateTime
	return result, nil
}

type deleteShipmentRequest struct {
	name       string
	shipperID  string
	shipmentID string
}

func (r *deleteShipmentRequest) parse(request *iamexamplev1.DeleteShipmentRequest) error {
	var v validation.MessageValidator
	// name = 1
	if request.Name == "" {
		v.AddFieldViolation("name", "required field")
	} else if err := resourcename.Sscan(
		request.Name,
		"shippers/{shipper}/shipments/{shipment}",
		&r.shipperID,
		&r.shipmentID,
	); err != nil {
		v.AddFieldError("name", err)
	}
	r.name = request.Name
	return v.Err()
}
