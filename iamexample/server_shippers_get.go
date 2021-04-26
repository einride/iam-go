package iamexample

import (
	"context"
	"fmt"

	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetShipper(
	ctx context.Context,
	request *iamexamplev1.GetShipperRequest,
) (*iamexamplev1.Shipper, error) {
	var shipperID string
	if err := resourcename.Sscan(request.GetName(), "shippers/{shipper}", &shipperID); err != nil {
		s.errorHook(ctx, err)
		return nil, s.badRequestError(ctx, &errdetails.BadRequest_FieldViolation{
			Field:       "name",
			Description: fmt.Sprintf("invalid format: %v", err),
		})
	}
	tx := s.Spanner.Single()
	defer tx.Close()
	row, err := iamexampledb.Query(tx).GetShippersRow(ctx, iamexampledb.GetShippersRowQuery{
		Key: iamexampledb.ShippersKey{
			ShipperId: shipperID,
		},
	})
	if err != nil {
		switch code := status.Code(err); code {
		case codes.NotFound:
			return nil, status.Errorf(code, "not found: %s", request.Name)
		}
		return nil, s.storageError(ctx, err)
	}
	msg, err := convertShipperRowToProto(row)
	if err != nil {
		s.errorHook(ctx, err)
		return nil, status.Error(codes.Internal, "internal data conversion error")
	}
	return msg, nil
}
