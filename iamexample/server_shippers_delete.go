package iamexample

import (
	"context"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourceid"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DeleteShipper implements iamexamplev1.FreightServiceServer.
func (s *Server) DeleteShipper(
	ctx context.Context,
	request *iamexamplev1.DeleteShipperRequest,
) (*longrunning.Operation, error) {
	var parsedRequest deleteShipperRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.deleteShipper(ctx, &parsedRequest)
}

func (s *Server) deleteShipper(
	ctx context.Context,
	request *deleteShipperRequest,
) (*longrunning.Operation, error) {
	var result *iamexamplev1.Shipper
	commitTime, err := s.Spanner.ReadWriteTransaction(
		ctx,
		func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
			row, err := iamexampledb.Query(tx).GetShippersRow(ctx, iamexampledb.GetShippersRowQuery{
				Key: iamexampledb.ShippersKey{ShipperId: request.shipperID},
			})
			if err != nil {
				return err
			}
			if row.DeleteTime.Valid {
				return status.Errorf(codes.FailedPrecondition, "shipper already deleted: %s", request.name)
			}
			row.UpdateTime = spanner.CommitTimestamp
			row.DeleteTime = spanner.NullTime{
				Time:  spanner.CommitTimestamp,
				Valid: true,
			}
			result, err = convertShipperRowToProto(row)
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
			return nil, status.Errorf(code, "no such shipper: %s", request.name)
		default:
			return nil, s.handleStorageError(ctx, err)
		}
	}
	result.UpdateTime = timestamppb.New(commitTime)
	result.DeleteTime = result.UpdateTime
	operationResponse, err := anypb.New(result)
	if err != nil {
		s.errorHook(ctx, err)
		return nil, status.Error(codes.Internal, "error marshaling operation response")
	}
	now := timestamppb.Now()
	operationMetadata, err := anypb.New(&iamexamplev1.DeleteShipperOperationMetadata{
		StartTime: now,
		EndTime:   now,
	})
	if err != nil {
		s.errorHook(ctx, err)
		return nil, status.Error(codes.Internal, "error marshaling operation metadata")
	}
	operation := &longrunning.Operation{
		Name:     resourcename.Sprint("operations/{operation}", resourceid.NewSystemGeneratedBase32()),
		Done:     true,
		Metadata: operationMetadata,
		Result: &longrunning.Operation_Response{
			Response: operationResponse,
		},
	}
	return operation, nil
}

type deleteShipperRequest struct {
	name      string
	shipperID string
}

func (r *deleteShipperRequest) parse(request *iamexamplev1.DeleteShipperRequest) error {
	var v validation.MessageValidator
	// name = 1
	if request.Name == "" {
		v.AddFieldViolation("name", "required field")
	} else if err := resourcename.Sscan(request.Name, "shippers/{shipper}", &r.shipperID); err != nil {
		v.AddFieldError("name", err)
	}
	r.name = request.Name
	return v.Err()
}
