package iamexample

import (
	"context"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourceid"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateShipper implements iamexamplev1.FreightServiceServer.
func (s *Server) CreateShipper(
	ctx context.Context,
	request *iamexamplev1.CreateShipperRequest,
) (*iamexamplev1.Shipper, error) {
	var parsedRequest createShipperRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.createShipper(ctx, &parsedRequest)
}

func (s *Server) createShipper(
	ctx context.Context,
	request *createShipperRequest,
) (*iamexamplev1.Shipper, error) {
	row, err := convertShipperProtoToRow(request.shipper)
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
			return nil, status.Errorf(code, "shipper %s already exists", request.shipper.GetName())
		default:
			return nil, s.handleStorageError(ctx, err)
		}
	}
	request.shipper.CreateTime = timestamppb.New(commitTime)
	request.shipper.UpdateTime = request.shipper.GetCreateTime()
	return request.shipper, nil
}

type createShipperRequest struct {
	shipperID string
	shipper   *iamexamplev1.Shipper
}

func (r *createShipperRequest) parse(request *iamexamplev1.CreateShipperRequest) error {
	var v validation.MessageValidator
	// shipper_id = 3
	if request.GetShipperId() != "" {
		if err := resourceid.ValidateUserSettable(request.GetShipperId()); err != nil {
			v.AddFieldError("shipper_id", err)
		}
		r.shipperID = request.GetShipperId()
	} else {
		r.shipperID = resourceid.NewSystemGeneratedBase32()
	}
	if request.GetShipper() == nil {
		v.AddFieldViolation("shipment", "required field")
	} else {
		// name = 1
		request.Shipper.Name = resourcename.Sprint("shippers/{shipper}", r.shipperID)
		// create_time = 2
		request.Shipper.CreateTime = nil
		// create_time = 3
		request.Shipper.CreateTime = nil
		// delete_time = 4
		request.Shipper.DeleteTime = nil
		// display_name = 5
		if len(request.GetShipper().GetDisplayName()) == 0 {
			v.AddFieldViolation("shipper.display_name", "required field")
		} else if len(request.GetShipper().GetDisplayName()) >= 64 {
			v.AddFieldViolation("shipper.display_name", "should be <= 63 characters")
		}
		r.shipper = request.GetShipper()
	}
	return v.Err()
}
