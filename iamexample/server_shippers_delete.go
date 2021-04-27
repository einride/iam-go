package iamexample

import (
	"context"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) DeleteShipper(
	ctx context.Context,
	request *iamexamplev1.DeleteShipperRequest,
) (*iamexamplev1.Shipper, error) {
	if err := s.validateDeleteShipper(ctx, request); err != nil {
		return nil, err
	}
	var shipperID string
	if err := resourcename.Sscan(request.Name, "shippers/{shipper}", &shipperID); err != nil {
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
			if row.DeleteTime.Valid {
				return status.Errorf(codes.FailedPrecondition, "shipper already deleted: %s", request.Name)
			}
			row.UpdateTime = spanner.CommitTimestamp
			row.DeleteTime = spanner.NullTime{
				Time:  spanner.CommitTimestamp,
				Valid: true,
			}
			result, err = convertShipperRowToProto(row)
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
			return nil, status.Errorf(code, "no such shipper: %s", request.Name)
		default:
			return nil, s.handleStorageError(ctx, err)
		}
	}
	result.UpdateTime = timestamppb.New(commitTime)
	result.DeleteTime = result.UpdateTime
	return result, nil
}

func (s *Server) validateDeleteShipper(ctx context.Context, request *iamexamplev1.DeleteShipperRequest) error {
	// TODO: Implement me.
	return nil
}
