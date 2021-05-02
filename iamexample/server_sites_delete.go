package iamexample

import (
	"context"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DeleteSite implements iamexamplev1.FreightServiceServer.
func (s *Server) DeleteSite(
	ctx context.Context,
	request *iamexamplev1.DeleteSiteRequest,
) (*iamexamplev1.Site, error) {
	var parsedRequest deleteSiteRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.deleteSite(ctx, &parsedRequest)
}

func (s *Server) deleteSite(
	ctx context.Context,
	request *deleteSiteRequest,
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
			if row.DeleteTime.Valid {
				return status.Errorf(codes.FailedPrecondition, "site already deleted: %s", request.name)
			}
			row.UpdateTime = spanner.CommitTimestamp
			row.DeleteTime = spanner.NullTime{
				Time:  spanner.CommitTimestamp,
				Valid: true,
			}
			result, err = convertSiteRowToProto(row)
			if err != nil {
				return err
			}
			return tx.BufferWrite([]*spanner.Mutation{spanner.Update(row.Mutate())})
		},
	)
	if err != nil {
		switch code := status.Code(err); code {
		case codes.FailedPrecondition:
			return nil, err
		case codes.NotFound:
			return nil, status.Errorf(code, "no such site: %s", request.name)
		default:
			return nil, s.handleStorageError(ctx, err)
		}
	}
	result.UpdateTime = timestamppb.New(commitTime)
	result.DeleteTime = result.UpdateTime
	return result, nil
}

type deleteSiteRequest struct {
	name      string
	shipperID string
	siteID    string
}

func (r *deleteSiteRequest) parse(request *iamexamplev1.DeleteSiteRequest) error {
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
