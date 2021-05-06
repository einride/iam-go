package iamexample

import (
	"context"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/fieldmask"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UpdateSite implements iamexamplev1.FreightServiceServer.
func (s *Server) UpdateSite(
	ctx context.Context,
	request *iamexamplev1.UpdateSiteRequest,
) (*iamexamplev1.Site, error) {
	var parsedRequest updateSiteRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.updateSite(ctx, &parsedRequest)
}

func (s *Server) updateSite(
	ctx context.Context,
	request *updateSiteRequest,
) (*iamexamplev1.Site, error) {
	var result *iamexamplev1.Site
	commitTime, err := s.Spanner.ReadWriteTransaction(
		ctx,
		func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
			row, err := iamexampledb.Query(tx).GetSitesRow(ctx, iamexampledb.GetSitesRowQuery{
				Key: iamexampledb.SitesKey{
					ShipperId: request.shipperID,
					SiteId:    request.siteID,
				},
			})
			if err != nil {
				return err
			}
			result, err = convertSiteRowToProto(row)
			if err != nil {
				return err
			}
			fieldmask.Update(request.updateMask, result, request.site)
			resultRow, err := convertSiteProtoToRow(result)
			if err != nil {
				return err
			}
			resultRow.UpdateTime = spanner.CommitTimestamp
			return tx.BufferWrite([]*spanner.Mutation{spanner.Update(resultRow.Mutate())})
		},
	)
	if err != nil {
		switch code := status.Code(err); code {
		case codes.NotFound:
			return nil, status.Errorf(code, "no such site: %s", request.site.Name)
		default:
			return nil, s.handleStorageError(ctx, err)
		}
	}
	result.UpdateTime = timestamppb.New(commitTime)
	return result, nil
}

type updateSiteRequest struct {
	shipperID  string
	siteID     string
	site       *iamexamplev1.Site
	updateMask *fieldmaskpb.FieldMask
}

func (r *updateSiteRequest) parse(request *iamexamplev1.UpdateSiteRequest) error {
	hasNoMask := len(request.GetUpdateMask().GetPaths()) == 0
	hasWildcardMask := len(request.UpdateMask.GetPaths()) == 1 && request.UpdateMask.Paths[0] == "/"
	has := func(path string) bool {
		if hasWildcardMask {
			return true
		}
		for _, maskPath := range request.GetUpdateMask().GetPaths() {
			if path == maskPath {
				return true
			}
		}
		return false
	}
	var v validation.MessageValidator
	// site = 1
	if request.Site == nil {
		v.AddFieldViolation("site", "required field")
	} else {
		r.site = request.Site
		// name = 1
		if len(request.Site.Name) == 0 {
			v.AddFieldViolation("site.name", "required field")
		} else if err := resourcename.Sscan(
			request.Site.Name,
			"shippers/{shipper}/sites/{site}",
			&r.shipperID,
			&r.siteID,
		); err != nil {
			v.AddFieldError("site.name", err)
		}
		// create_time = 2
		request.Site.CreateTime = nil
		// update_time = 3
		request.Site.UpdateTime = nil
		// delete_time = 4
		request.Site.DeleteTime = nil
		// display_name = 5
		if has("display_name") || hasNoMask && len(request.Site.DisplayName) > 0 {
			if len(request.Site.DisplayName) == 0 {
				v.AddFieldViolation("site.display_name", "required field")
			} else if len(request.Site.DisplayName) >= 64 {
				v.AddFieldViolation("site.display_name", "should be <= 63 characters")
			}
		}
		// lat_lng = 6
		if has("lat_lng") || hasNoMask && request.Site.LatLng != nil {
			if !(-90 <= request.Site.LatLng.Latitude && request.Site.LatLng.Latitude <= 90) {
				v.AddFieldViolation("site.lat_lng.latitude", "must be in the range [-90.0, +90.0]")
			}
			if !(-180 <= request.Site.LatLng.Longitude && request.Site.LatLng.Longitude <= 180) {
				v.AddFieldViolation("site.lat_lng.longitude", "must be in the range [-180.0, +180.0]")
			}
		}
	}
	return v.Err()
}
