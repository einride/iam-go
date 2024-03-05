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

// GetShipper implements iamexamplev1.FreightServiceServer.
func (s *Server) GetShipper(
	ctx context.Context,
	request *iamexamplev1.GetShipperRequest,
) (*iamexamplev1.Shipper, error) {
	var parsedRequest getShipperRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.getShipper(ctx, &parsedRequest)
}

func (s *Server) getShipper(
	ctx context.Context,
	request *getShipperRequest,
) (*iamexamplev1.Shipper, error) {
	tx := s.Spanner.Single()
	defer tx.Close()
	row, err := iamexampledb.Query(tx).GetShippersRow(ctx, iamexampledb.GetShippersRowQuery{
		Key: iamexampledb.ShippersKey{
			ShipperId: request.shipperID,
		},
	})
	if err != nil {
		if code := status.Code(err); code == codes.NotFound {
			return nil, status.Errorf(code, "not found: %s", request.name)
		}
		return nil, s.handleStorageError(ctx, err)
	}
	msg, err := convertShipperRowToProto(row)
	if err != nil {
		s.errorHook(ctx, err)
		return nil, status.Error(codes.Internal, "internal data conversion error")
	}
	return msg, nil
}

type getShipperRequest struct {
	name      string
	shipperID string
}

func (r *getShipperRequest) parse(request *iamexamplev1.GetShipperRequest) error {
	var v validation.MessageValidator
	// name = 1
	if request.GetName() == "" {
		v.AddFieldViolation("name", "required field")
	} else if err := resourcename.Sscan(
		request.GetName(),
		"shippers/{shipper}",
		&r.shipperID,
	); err != nil {
		v.AddFieldError("name", err)
	}
	r.name = request.GetName()
	return v.Err()
}
