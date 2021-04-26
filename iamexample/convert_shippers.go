package iamexample

import (
	"fmt"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertShipperProtoToRow(msg *iamexamplev1.Shipper) (*iamexampledb.ShippersRow, error) {
	var row iamexampledb.ShippersRow
	// name = 1
	if err := resourcename.Sscan(msg.GetName(), "shippers/{shipper}", &row.ShipperId); err != nil {
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
	return &row, nil
}

func convertShipperRowToProto(row *iamexampledb.ShippersRow) (*iamexamplev1.Shipper, error) {
	var msg iamexamplev1.Shipper
	// name = 1
	msg.Name = resourcename.Sprint("shippers/{shipper}", row.ShipperId)
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
	return &msg, nil
}
