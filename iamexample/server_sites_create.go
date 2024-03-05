package iamexample

import (
	"context"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourceid"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateSite implements iamexamplev1.FreightServiceServer.
func (s *Server) CreateSite(
	ctx context.Context,
	request *iamexamplev1.CreateSiteRequest,
) (*iamexamplev1.Site, error) {
	var parsedRequest createSiteRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.createSite(ctx, &parsedRequest)
}

func (s *Server) createSite(
	ctx context.Context,
	request *createSiteRequest,
) (*iamexamplev1.Site, error) {
	row, err := convertSiteProtoToRow(request.site)
	if err != nil {
		s.errorHook(ctx, err)
		return nil, status.Error(codes.Internal, "internal data conversion error")
	}
	row.CreateTime = spanner.CommitTimestamp
	row.UpdateTime = spanner.CommitTimestamp
	commitTime, err := s.Spanner.Apply(ctx, []*spanner.Mutation{spanner.Insert(row.Mutate())})
	if err != nil {
		switch code := status.Code(err); code {
		case codes.AlreadyExists:
			return nil, status.Errorf(code, "site %s already exists", request.site.GetName())
		default:
			return nil, s.handleStorageError(ctx, err)
		}
	}
	request.site.CreateTime = timestamppb.New(commitTime)
	request.site.UpdateTime = request.site.GetCreateTime()
	return request.site, nil
}

type createSiteRequest struct {
	shipperID string
	siteID    string
	site      *iamexamplev1.Site
}

func (r *createSiteRequest) parse(request *iamexamplev1.CreateSiteRequest) error {
	var v validation.MessageValidator
	// parent = 1
	if request.GetParent() == "" {
		v.AddFieldViolation("parent", "required field")
	} else if resourcename.ContainsWildcard(request.GetParent()) {
		v.AddFieldViolation("parent", "must not contain wildcards")
	} else if err := resourcename.Sscan(request.GetParent(), "shippers/{shipper}", &r.shipperID); err != nil {
		v.AddFieldViolation("parent", "invalid format")
	}
	// site_id = 3
	if request.GetSiteId() != "" {
		if err := resourceid.ValidateUserSettable(request.GetSiteId()); err != nil {
			v.AddFieldError("site_id", err)
		}
		r.siteID = request.GetSiteId()
	} else {
		r.siteID = resourceid.NewSystemGeneratedBase32()
	}
	// site = 2
	if request.GetSite() == nil {
		v.AddFieldViolation("site", "required field")
	} else {
		// name = 1
		request.Site.Name = resourcename.Sprint(
			"shippers/{shipper}/sites/{site}",
			r.shipperID,
			r.siteID,
		)
		// create_time = 2
		request.Site.CreateTime = nil
		// update_time = 3
		request.Site.UpdateTime = nil
		// delete_time = 4
		request.Site.DeleteTime = nil
		// display_name = 5
		if len(request.GetSite().GetDisplayName()) == 0 {
			v.AddFieldViolation("site.display_name", "required field")
		} else if len(request.GetSite().GetDisplayName()) >= 64 {
			v.AddFieldViolation("site.display_name", "should be <= 63 characters")
		}
		// lat_lng = 6
		if request.GetSite().GetLatLng() != nil {
			if !(-90 <= request.GetSite().GetLatLng().GetLatitude() && request.GetSite().GetLatLng().GetLatitude() <= 90) {
				v.AddFieldViolation("site.lat_lng.latitude", "must be in the range [-90.0, +90.0]")
			}
			if !(-180 <= request.GetSite().GetLatLng().GetLongitude() && request.GetSite().GetLatLng().GetLongitude() <= 180) {
				v.AddFieldViolation("site.lat_lng.longitude", "must be in the range [-180.0, +180.0]")
			}
		}
		r.site = request.GetSite()
	}
	return v.Err()
}
