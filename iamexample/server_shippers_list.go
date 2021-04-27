package iamexample

import (
	"context"

	"go.einride.tech/aip/pagination"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ListShippers(
	ctx context.Context,
	request *iamexamplev1.ListShippersRequest,
) (*iamexamplev1.ListShippersResponse, error) {
	pageToken, err := pagination.ParsePageToken(request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid page token")
	}
	response := iamexamplev1.ListShippersResponse{
		Shippers: make([]*iamexamplev1.Shipper, 0, request.PageSize+1),
	}
	tx := s.Spanner.Single()
	defer tx.Close()
	if err := iamexampledb.Query(tx).ListShippersRows(ctx, iamexampledb.ListShippersRowsQuery{
		Limit:  request.PageSize + 1,
		Offset: pageToken.Offset,
	}).Do(func(row *iamexampledb.ShippersRow) error {
		msg, err := convertShipperRowToProto(row)
		if err != nil {
			return err
		}
		response.Shippers = append(response.Shippers, msg)
		return nil
	}); err != nil {
		return nil, s.handleStorageError(ctx, err)
	}
	if len(response.Shippers) > int(request.PageSize) {
		response.Shippers = response.Shippers[:request.PageSize]
		response.NextPageToken = pageToken.Next(request).String()
	}
	return &response, nil
}
