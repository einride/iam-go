package iamexample

import (
	"context"

	"go.einride.tech/aip/pagination"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
)

// ListShippers implements iamexamplev1.FreightServiceServer.
func (s *Server) ListShippers(
	ctx context.Context,
	request *iamexamplev1.ListShippersRequest,
) (*iamexamplev1.ListShippersResponse, error) {
	var parsedRequest listShippersRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.listShippers(ctx, &parsedRequest)
}

func (s *Server) listShippers(
	ctx context.Context,
	request *listShippersRequest,
) (*iamexamplev1.ListShippersResponse, error) {
	response := iamexamplev1.ListShippersResponse{
		Shippers: make([]*iamexamplev1.Shipper, 0, request.pageSize+1),
	}
	tx := s.Spanner.Single()
	defer tx.Close()
	if err := iamexampledb.Query(tx).ListShippersRows(ctx, iamexampledb.ListShippersRowsQuery{
		Limit:  request.pageSize + 1,
		Offset: request.pageToken.Offset,
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
	if len(response.Shippers) > int(request.pageSize) {
		response.Shippers = response.Shippers[:request.pageSize]
		response.NextPageToken = request.nextPageToken()
	}
	return &response, nil
}

type listShippersRequest struct {
	protoRequest *iamexamplev1.ListShippersRequest
	pageSize     int32
	pageToken    pagination.PageToken
}

func (r *listShippersRequest) parse(request *iamexamplev1.ListShippersRequest) error {
	var v validation.MessageValidator
	const (
		defaultPageSize = 100
		maxPageSize     = 1_000
	)
	switch {
	case request.PageSize < 0:
		v.AddFieldViolation("page_size", "must be >= 0")
	case request.PageSize == 0:
		r.pageSize = defaultPageSize
	case request.PageSize > maxPageSize:
		r.pageSize = maxPageSize
	default:
		r.pageSize = request.PageSize
	}
	if pageToken, err := pagination.ParsePageToken(request); err != nil {
		v.AddFieldError("page_token", err)
	} else {
		r.pageToken = pageToken
	}
	r.protoRequest = request
	return v.Err()
}

func (r *listShippersRequest) nextPageToken() string {
	return r.pageToken.Next(r.protoRequest).String()
}
