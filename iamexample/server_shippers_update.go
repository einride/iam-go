package iamexample

import (
	"context"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/fieldbehavior"
	"go.einride.tech/aip/fieldmask"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) UpdateShipper(
	ctx context.Context,
	request *iamexamplev1.UpdateShipperRequest,
) (*iamexamplev1.Shipper, error) {
	if err := s.validateUpdateShipper(ctx, request); err != nil {
		return nil, err
	}
	var shipperID string
	if err := resourcename.Sscan(request.Shipper.Name, "shippers/{shipper}", &shipperID); err != nil {
		s.errorHook(ctx, err)
		return nil, status.Error(codes.Internal, "internal parse error")
	}
	var result *iamexamplev1.Shipper
	commitTime, err := s.Spanner.ReadWriteTransaction(
		ctx,
		func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
			row, err := iamexampledb.Query(tx).GetShippersRow(ctx, iamexampledb.GetShippersRowQuery{
				Key: iamexampledb.ShippersKey{ShipperId: shipperID},
			})
			if err != nil {
				return err
			}
			result, err = convertShipperRowToProto(row)
			if err != nil {
				return err
			}
			fieldbehavior.ClearFields(
				request.Shipper,
				annotations.FieldBehavior_OUTPUT_ONLY,
				annotations.FieldBehavior_IMMUTABLE,
			)
			fieldmask.Update(request.UpdateMask, result, request.Shipper)
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
			return nil, status.Errorf(code, "no such shipper: %s", request.Shipper.Name)
		default:
			return nil, s.storageError(ctx, err)
		}
	}
	result.UpdateTime = timestamppb.New(commitTime)
	return result, nil
}

func (s *Server) validateUpdateShipper(ctx context.Context, request *iamexamplev1.UpdateShipperRequest) error {
	// TODO: Implement me.
	return nil
}
