package iamexample

import (
	"context"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourceid"
	"go.einride.tech/aip/resourcename"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateShipper(
	ctx context.Context,
	request *iamexamplev1.CreateShipperRequest,
) (*iamexamplev1.Shipper, error) {
	if err := s.validateCreateShipper(ctx, request); err != nil {
		return nil, err
	}
	if request.ShipperId == "" {
		request.ShipperId = resourceid.NewSystemGeneratedBase32()
	}
	request.Shipper.Name = resourcename.Sprint("shippers/{shipper}", request.ShipperId)
	row, err := convertShipperProtoToRow(request.Shipper)
	if err != nil {
		s.errorHook(ctx, err)
		return nil, status.Error(codes.Internal, "internal data conversion error")
	}
	row.CreateTime = spanner.CommitTimestamp
	row.UpdateTime = spanner.CommitTimestamp
	commitTime, err := s.Spanner.Apply(ctx, []*spanner.Mutation{spanner.Insert(row.Mutate())})
	if err != nil {
		switch code := status.Code(err); code {
		case codes.AlreadyExists:
			return nil, status.Errorf(code, "shipper %s already exists", request.Shipper.Name)
		default:
			return nil, s.storageError(ctx, err)
		}
	}
	request.Shipper.CreateTime = timestamppb.New(commitTime)
	request.Shipper.UpdateTime = request.Shipper.CreateTime
	return request.Shipper, nil
}

func (s *Server) validateCreateShipper(ctx context.Context, request *iamexamplev1.CreateShipperRequest) error {
	// TODO: Implement me.
	return nil
}
