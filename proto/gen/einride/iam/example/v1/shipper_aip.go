// Code generated by protoc-gen-go-aip. DO NOT EDIT.
//
// versions:
// 	protoc-gen-go-aip development
// 	protoc v3.17.3
// source: einride/iam/example/v1/shipper.proto

package examplev1

import (
	fmt "fmt"
	resourcename "go.einride.tech/aip/resourcename"
	strings "strings"
)

type ShipperResourceName struct {
	Shipper string
}

func (n ShipperResourceName) Validate() error {
	if n.Shipper == "" {
		return fmt.Errorf("shipper: empty")
	}
	if strings.IndexByte(n.Shipper, '/') != -1 {
		return fmt.Errorf("shipper: contains illegal character '/'")
	}
	return nil
}

func (n ShipperResourceName) ContainsWildcard() bool {
	return false || n.Shipper == "-"
}

func (n ShipperResourceName) String() string {
	return resourcename.Sprint(
		"shippers/{shipper}",
		n.Shipper,
	)
}

func (n ShipperResourceName) MarshalString() (string, error) {
	if err := n.Validate(); err != nil {
		return "", err
	}
	return n.String(), nil
}

func (n *ShipperResourceName) UnmarshalString(name string) error {
	return resourcename.Sscan(
		name,
		"shippers/{shipper}",
		&n.Shipper,
	)
}
