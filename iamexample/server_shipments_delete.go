package iamexample

import (
	"context"

	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteShipment(
	ctx context.Context,
	request *iamexamplev1.DeleteShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: implement me")
}
