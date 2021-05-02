package iamexample

import (
	"context"

	"go.einride.tech/aip/pagination"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
)

// ListShipments implements iamexamplev1.FreightServiceServer.
func (s *Server) ListShipments(
	ctx context.Context,
	request *iamexamplev1.ListShipmentsRequest,
) (*iamexamplev1.ListShipmentsResponse, error) {
	var parsedRequest listShipmentsRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.listShipments(ctx, &parsedRequest)
}

func (s *Server) listShipments(
	ctx context.Context,
	request *listShipmentsRequest,
) (*iamexamplev1.ListShipmentsResponse, error) {
	response := iamexamplev1.ListShipmentsResponse{
		Shipments: make([]*iamexamplev1.Shipment, 0, request.pageSize+1),
	}
	tx := s.Spanner.Single()
	defer tx.Close()
	if err := iamexampledb.Query(tx).ListShipmentsRows(ctx, iamexampledb.ListShipmentsRowsQuery{
		Where:     iamexampledb.ShippersKey{ShipperId: request.shipperID}.BoolExpr(),
		Limit:     request.pageSize + 1,
		Offset:    request.pageToken.Offset,
		LineItems: true,
	}).Do(func(row *iamexampledb.ShipmentsRow) error {
		msg, err := convertShipmentRowToProto(row)
		if err != nil {
			return err
		}
		response.Shipments = append(response.Shipments, msg)
		return nil
	}); err != nil {
		return nil, s.handleStorageError(ctx, err)
	}
	if len(response.Shipments) > int(request.pageSize) {
		response.Shipments = response.Shipments[:request.pageSize]
		response.NextPageToken = request.nextPageToken()
	}
	return &response, nil
}

type listShipmentsRequest struct {
	protoRequest *iamexamplev1.ListShipmentsRequest
	shipperID    string
	pageSize     int32
	pageToken    pagination.PageToken
}

func (r *listShipmentsRequest) parse(request *iamexamplev1.ListShipmentsRequest) error {
	var v validation.MessageValidator
	const (
		defaultPageSize = 100
		maxPageSize     = 1_000
	)
	if request.PageSize < 0 {
		v.AddFieldViolation("page_size", "must be >= 0")
	} else if request.PageSize == 0 {
		r.pageSize = defaultPageSize
	} else if request.PageSize > maxPageSize {
		r.pageSize = maxPageSize
	} else {
		r.pageSize = request.PageSize
	}
	if request.Parent == "" {
		v.AddFieldViolation("parent", "missing required field")
	} else if resourcename.ContainsWildcard(request.Parent) {
		v.AddFieldViolation("parent", "wildcard not allowed")
	} else if err := resourcename.Sscan(request.Parent, "shippers/{shipper}", &r.shipperID); err != nil {
		v.AddFieldError("parent", err)
	}
	if pageToken, err := pagination.ParsePageToken(request); err != nil {
		v.AddFieldError("page_token", err)
	} else {
		r.pageToken = pageToken
	}
	r.protoRequest = request
	return v.Err()
}

func (r *listShipmentsRequest) nextPageToken() string {
	return r.pageToken.Next(r.protoRequest).String()
}
