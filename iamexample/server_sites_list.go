package iamexample

import (
	"context"

	"go.einride.tech/aip/pagination"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
)

// ListSites implements iamexamplev1.FreightServiceServer.
func (s *Server) ListSites(
	ctx context.Context,
	request *iamexamplev1.ListSitesRequest,
) (*iamexamplev1.ListSitesResponse, error) {
	var parsedRequest listSitesRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.listSites(ctx, &parsedRequest)
}

func (s *Server) listSites(
	ctx context.Context,
	request *listSitesRequest,
) (*iamexamplev1.ListSitesResponse, error) {
	response := iamexamplev1.ListSitesResponse{
		Sites: make([]*iamexamplev1.Site, 0, request.pageSize+1),
	}
	tx := s.Spanner.Single()
	defer tx.Close()
	if err := iamexampledb.Query(tx).ListSitesRows(ctx, iamexampledb.ListSitesRowsQuery{
		Where:  iamexampledb.ShippersKey{ShipperId: request.shipperID}.BoolExpr(),
		Limit:  request.pageSize + 1,
		Offset: request.pageToken.Offset,
	}).Do(func(row *iamexampledb.SitesRow) error {
		msg, err := convertSiteRowToProto(row)
		if err != nil {
			return err
		}
		response.Sites = append(response.Sites, msg)
		return nil
	}); err != nil {
		return nil, s.handleStorageError(ctx, err)
	}
	if len(response.Sites) > int(request.pageSize) {
		response.Sites = response.Sites[:request.pageSize]
		response.NextPageToken = request.nextPageToken()
	}
	return &response, nil
}

type listSitesRequest struct {
	protoRequest *iamexamplev1.ListSitesRequest
	shipperID    string
	pageSize     int32
	pageToken    pagination.PageToken
}

func (r *listSitesRequest) parse(request *iamexamplev1.ListSitesRequest) error {
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

func (r *listSitesRequest) nextPageToken() string {
	return r.pageToken.Next(r.protoRequest).String()
}
