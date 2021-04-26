package iamexample

import (
	"context"

	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetShipment(
	ctx context.Context,
	request *iamexamplev1.GetShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: implement me")
}
