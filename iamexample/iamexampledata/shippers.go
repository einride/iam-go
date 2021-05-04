package iamexampledata

import (
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
)

// Einride returns the Einride example shipper.
func Einride() *iamexamplev1.Shipper {
	return &iamexamplev1.Shipper{
		Name:        "shippers/einride",
		DisplayName: "Einride",
	}
}
