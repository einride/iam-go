package iamexample

import (
	"context"

	"cloud.google.com/go/longrunning/autogen/longrunningpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ListOperations implements longrunningpb.OperationsServer.
func (s *Server) ListOperations(
	ctx context.Context,
	request *longrunningpb.ListOperationsRequest,
) (*longrunningpb.ListOperationsResponse, error) {
	return &longrunningpb.ListOperationsResponse{}, nil
}

// GetOperation implements longrunningpb.OperationsServer.
func (s *Server) GetOperation(
	ctx context.Context,
	request *longrunningpb.GetOperationRequest,
) (*longrunningpb.Operation, error) {
	return nil, status.Error(codes.NotFound, "operation not found")
}

// DeleteOperation implements longrunningpb.OperationsServer.
func (s *Server) DeleteOperation(
	ctx context.Context,
	request *longrunningpb.DeleteOperationRequest,
) (*emptypb.Empty, error) {
	return nil, status.Error(codes.NotFound, "operation not found")
}

// CancelOperation implements longrunningpb.OperationsServer.
func (s *Server) CancelOperation(
	ctx context.Context,
	request *longrunningpb.CancelOperationRequest,
) (*emptypb.Empty, error) {
	return nil, status.Error(codes.NotFound, "operation not found")
}

// WaitOperation implements longrunningpb.OperationsServer.
func (s *Server) WaitOperation(
	ctx context.Context,
	request *longrunningpb.WaitOperationRequest,
) (*longrunningpb.Operation, error) {
	return nil, status.Error(codes.NotFound, "operation not found")
}
