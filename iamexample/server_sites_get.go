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

// GetSite implements iamexamplev1.FreightServiceServer.
func (s *Server) GetSite(
	ctx context.Context,
	request *iamexamplev1.GetSiteRequest,
) (*iamexamplev1.Site, error) {
	var parsedRequest getSiteRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.getSite(ctx, &parsedRequest)
}

func (s *Server) getSite(
	ctx context.Context,
	request *getSiteRequest,
) (*iamexamplev1.Site, error) {
	tx := s.Spanner.Single()
	defer tx.Close()
	row, err := iamexampledb.Query(tx).GetSitesRow(ctx, iamexampledb.GetSitesRowQuery{
		Key: iamexampledb.SitesKey{
			ShipperId: request.shipperID,
			SiteId:    request.siteID,
		},
	})
	if err != nil {
		if code := status.Code(err); code == codes.NotFound {
			return nil, status.Errorf(code, "not found: %s", request.name)
		}
		return nil, s.handleStorageError(ctx, err)
	}
	msg, err := convertSiteRowToProto(row)
	if err != nil {
		s.errorHook(ctx, err)
		return nil, status.Error(codes.Internal, "internal data conversion error")
	}
	return msg, nil
}

type getSiteRequest struct {
	name      string
	shipperID string
	siteID    string
}

func (r *getSiteRequest) parse(request *iamexamplev1.GetSiteRequest) error {
	var v validation.MessageValidator
	// name = 1
	if request.Name == "" {
		v.AddFieldViolation("name", "required field")
	} else if err := resourcename.Sscan(
		request.Name,
		"shippers/{shipper}/sites/{site}",
		&r.shipperID,
		&r.siteID,
	); err != nil {
		v.AddFieldError("name", err)
	}
	r.name = request.Name
	return v.Err()
}
