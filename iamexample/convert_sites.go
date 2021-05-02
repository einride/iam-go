package iamexample

import (
	"fmt"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertSiteProtoToRow(msg *iamexamplev1.Site) (*iamexampledb.SitesRow, error) {
	var row iamexampledb.SitesRow
	// name = 1
	if err := resourcename.Sscan(
		msg.GetName(),
		"shippers/{shipper}/sites/{site}",
		&row.ShipperId,
		&row.SiteId,
	); err != nil {
		return nil, fmt.Errorf("name: %w", err)
	}
	// create_time = 2
	row.CreateTime = msg.GetCreateTime().AsTime()
	// update_time = 3
	row.UpdateTime = msg.GetUpdateTime().AsTime()
	// delete_time = 4
	if msg.GetDeleteTime() != nil {
		row.DeleteTime = spanner.NullTime{
			Time:  msg.GetDeleteTime().AsTime(),
			Valid: true,
		}
	}
	// display_name = 5
	row.DisplayName = msg.GetDisplayName()
	if err := row.Validate(); err != nil {
		return nil, fmt.Errorf("validate row: %w", err)
	}
	// lat_lng = 6
	if msg.GetLatLng() != nil {
		row.Latitude = spanner.NullFloat64{
			Float64: msg.GetLatLng().GetLatitude(),
			Valid:   true,
		}
		row.Longitude = spanner.NullFloat64{
			Float64: msg.GetLatLng().GetLongitude(),
			Valid:   true,
		}
	}
	return &row, nil
}

func convertSiteRowToProto(row *iamexampledb.SitesRow) (*iamexamplev1.Site, error) {
	var msg iamexamplev1.Site
	// name = 1
	msg.Name = resourcename.Sprint("shippers/{shipper}/sites/{site}", row.ShipperId, row.SiteId)
	// create_time = 2
	msg.CreateTime = timestamppb.New(row.CreateTime)
	// update_time = 3
	msg.UpdateTime = timestamppb.New(row.UpdateTime)
	// delete_time = 4
	if row.DeleteTime.Valid {
		msg.DeleteTime = timestamppb.New(row.DeleteTime.Time)
	}
	// display_name = 5
	msg.DisplayName = row.DisplayName
	// lat_lng = 6
	if row.Latitude.Valid && row.Longitude.Valid {
		msg.LatLng = &latlng.LatLng{
			Latitude:  row.Latitude.Float64,
			Longitude: row.Longitude.Float64,
		}
	}
	return &msg, nil
}
