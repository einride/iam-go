package iamexample

import (
	"context"
	"fmt"

	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BatchGetSites implements iamexamplev1.FreightServiceServer.
func (s *Server) BatchGetSites(
	ctx context.Context,
	request *iamexamplev1.BatchGetSitesRequest,
) (*iamexamplev1.BatchGetSitesResponse, error) {
	var parsedRequest batchGetSitesRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.batchGetSites(ctx, &parsedRequest)
}

func (s *Server) batchGetSites(
	ctx context.Context,
	request *batchGetSitesRequest,
) (*iamexamplev1.BatchGetSitesResponse, error) {
	tx := s.Spanner.Single()
	defer tx.Close()
	rows, err := iamexampledb.Query(tx).BatchGetSitesRows(ctx, iamexampledb.BatchGetSitesRowsQuery{
		Keys: request.keys,
	})
	if err != nil {
		return nil, s.handleStorageError(ctx, err)
	}
	response := iamexamplev1.BatchGetSitesResponse{
		Sites: make([]*iamexamplev1.Site, 0, len(request.keys)),
	}
	for i, key := range request.keys {
		row, ok := rows[key]
		if !ok {
			return nil, status.Errorf(codes.NotFound, "not found: %s", request.names[i])
		}
		msg, err := convertSiteRowToProto(row)
		if err != nil {
			s.errorHook(ctx, err)
			return nil, status.Errorf(codes.Internal, "internal data conversion error")
		}
		response.Sites = append(response.Sites, msg)
	}
	return &response, nil
}

type batchGetSitesRequest struct {
	parent iamexampledb.ShippersKey
	names  []string
	keys   []iamexampledb.SitesKey
}

func (r *batchGetSitesRequest) parse(request *iamexamplev1.BatchGetSitesRequest) error {
	var v validation.MessageValidator
	// parent = 1
	if request.Parent != "" {
		if err := resourcename.Sscan(request.Parent, "shippers/{shipper}", &r.parent.ShipperId); err != nil {
			v.AddFieldError("parent", err)
		}
	}
	// names = 2
	if len(request.Names) == 0 {
		v.AddFieldViolation("names", "required field")
	}
	r.keys = make([]iamexampledb.SitesKey, 0, len(request.Names))
	r.names = request.Names
	for i, name := range request.Names {
		if resourcename.ContainsWildcard(name) {
			v.AddFieldViolation(fmt.Sprintf("names[%d]", i), "wildcard not supported")
		}
		if request.Parent != "" && !resourcename.HasParent(name, request.Parent) {
			v.AddFieldViolation(fmt.Sprintf("names[%d]", i), "must be a child of parent %s", request.Parent)
		}
		var key iamexampledb.SitesKey
		if err := resourcename.Sscan(
			name,
			"shippers/{shipper}/sites/{site}",
			&key.ShipperId,
			&key.SiteId,
		); err != nil {
			v.AddFieldError(fmt.Sprintf("names[%d]", i), err)
			continue
		}
		r.keys = append(r.keys, key)
	}
	return v.Err()
}
