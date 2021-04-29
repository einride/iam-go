package iamexample

import (
	"context"

	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteSite(
	ctx context.Context,
	request *iamexamplev1.DeleteSiteRequest,
) (*iamexamplev1.Site, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: implement me")
}