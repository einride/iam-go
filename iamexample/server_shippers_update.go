package iamexample

import (
	"context"

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

// UpdateShipper implements iamexamplev1.FreightServiceServer.
func (s *Server) UpdateShipper(
	ctx context.Context,
	request *iamexamplev1.UpdateShipperRequest,
) (*iamexamplev1.Shipper, error) {
	var parsedRequest updateShipperRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.updateShipper(ctx, &parsedRequest)
}

func (s *Server) updateShipper(
	ctx context.Context,
	request *updateShipperRequest,
) (*iamexamplev1.Shipper, error) {
	var result *iamexamplev1.Shipper
	commitTime, err := s.Spanner.ReadWriteTransaction(
		ctx,
		func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
			row, err := iamexampledb.Query(tx).GetShippersRow(ctx, iamexampledb.GetShippersRowQuery{
				Key: iamexampledb.ShippersKey{ShipperId: request.shipperID},
			})
			if err != nil {
				return err
			}
			result, err = convertShipperRowToProto(row)
			if err != nil {
				return err
			}
			fieldmask.Update(request.updateMask, result, request.shipper)
			resultRow, err := convertShipperProtoToRow(result)
			if err != nil {
				return err
			}
			resultRow.UpdateTime = spanner.CommitTimestamp
			return tx.BufferWrite([]*spanner.Mutation{spanner.Update(resultRow.Mutate())})
		},
	)
	if err != nil {
		switch code := status.Code(err); code {
		case codes.NotFound:
			return nil, status.Errorf(code, "no such shipper: %s", request.shipper.Name)
		default:
			return nil, s.handleStorageError(ctx, err)
		}
	}
	result.UpdateTime = timestamppb.New(commitTime)
	return result, nil
}

type updateShipperRequest struct {
	shipperID  string
	shipper    *iamexamplev1.Shipper
	updateMask *fieldmaskpb.FieldMask
}

func (r *updateShipperRequest) parse(request *iamexamplev1.UpdateShipperRequest) error {
	hasNoMask := len(request.GetUpdateMask().GetPaths()) == 0
	hasWildcardMask := len(request.UpdateMask.GetPaths()) == 1 && request.UpdateMask.Paths[0] == "*"
	has := func(path string) bool {
		if hasWildcardMask {
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
	if request.Shipper == nil {
		v.AddFieldViolation("shipper", "required field")
	} else {
		// name = 1
		if len(request.Shipper.Name) == 0 {
			v.AddFieldViolation("shipper.name", "required field")
		} else if err := resourcename.Sscan(
			request.Shipper.Name,
			"shippers/{shipper}",
			&r.shipperID,
		); err != nil {
			v.AddFieldError("shipper.name", err)
		}
		// create_time = 2
		request.Shipper.CreateTime = nil
		// update_time = 3
		request.Shipper.UpdateTime = nil
		// delete_time = 4
		request.Shipper.DeleteTime = nil
		// display_name = 5
		if has("display_name") || hasNoMask && len(request.Shipper.DisplayName) > 0 {
			if len(request.Shipper.DisplayName) == 0 {
				v.AddFieldViolation("shipper.display_name", "required field")
			} else if len(request.Shipper.DisplayName) >= 64 {
				v.AddFieldViolation("shipper.display_name", "should be <= 63 characters")
			}
		}
		r.shipper = request.Shipper
	}
	// update_mask = 2
	if err := fieldmask.Validate(request.UpdateMask, request.Shipper); err != nil {
		v.AddFieldError("update_mask", err)
	} else {
		r.updateMask = request.UpdateMask
	}
	return v.Err()
}
