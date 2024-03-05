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

// BatchGetShipments implements iamexamplev1.FreightServiceServer.
func (s *Server) BatchGetShipments(
	ctx context.Context,
	request *iamexamplev1.BatchGetShipmentsRequest,
) (*iamexamplev1.BatchGetShipmentsResponse, error) {
	var parsedRequest batchGetShipmentsRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.batchGetShipments(ctx, &parsedRequest)
}

func (s *Server) batchGetShipments(
	ctx context.Context,
	request *batchGetShipmentsRequest,
) (*iamexamplev1.BatchGetShipmentsResponse, error) {
	tx := s.Spanner.ReadOnlyTransaction()
	defer tx.Close()
	rows, err := iamexampledb.Query(tx).BatchGetShipmentsRows(ctx, iamexampledb.BatchGetShipmentsRowsQuery{
		Keys:      request.keys,
		LineItems: true,
	})
	if err != nil {
		return nil, s.handleStorageError(ctx, err)
	}
	response := iamexamplev1.BatchGetShipmentsResponse{
		Shipments: make([]*iamexamplev1.Shipment, 0, len(request.keys)),
	}
	for i, key := range request.keys {
		row, ok := rows[key]
		if !ok {
			return nil, status.Errorf(codes.NotFound, "not found: %s", request.names[i])
		}
		msg, err := convertShipmentRowToProto(row)
		if err != nil {
			s.errorHook(ctx, err)
			return nil, status.Errorf(codes.Internal, "internal data conversion error")
		}
		response.Shipments = append(response.Shipments, msg)
	}
	return &response, nil
}

type batchGetShipmentsRequest struct {
	parent iamexampledb.ShippersKey
	names  []string
	keys   []iamexampledb.ShipmentsKey
}

func (r *batchGetShipmentsRequest) parse(request *iamexamplev1.BatchGetShipmentsRequest) error {
	var v validation.MessageValidator
	// parent = 1
	if request.GetParent() != "" {
		if err := resourcename.Sscan(request.GetParent(), "shippers/{shipper}", &r.parent.ShipperId); err != nil {
			v.AddFieldError("parent", err)
		}
	}
	// names = 2
	if len(request.GetNames()) == 0 {
		v.AddFieldViolation("names", "required field")
	}
	r.keys = make([]iamexampledb.ShipmentsKey, 0, len(request.GetNames()))
	r.names = request.GetNames()
	for i, name := range request.GetNames() {
		if resourcename.ContainsWildcard(name) {
			v.AddFieldViolation(fmt.Sprintf("names[%d]", i), "wildcard not supported")
		}
		if request.GetParent() != "" && !resourcename.HasParent(name, request.GetParent()) {
			v.AddFieldViolation(fmt.Sprintf("names[%d]", i), "%s is not a child of parent %s", name, request.GetParent())
		}
		var key iamexampledb.ShipmentsKey
		if err := resourcename.Sscan(
			name,
			"shippers/{shipper}/shipments/{shipment}",
			&key.ShipperId,
			&key.ShipmentId,
		); err != nil {
			v.AddFieldError(fmt.Sprintf("names[%d]", i), err)
			continue
		}
		r.keys = append(r.keys, key)
	}
	return v.Err()
}
