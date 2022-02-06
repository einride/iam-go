package iamexample

import (
	"context"

	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ListOperations implements longrunning.OperationsServer.
func (s *Server) ListOperations(
	ctx context.Context,
	request *longrunning.ListOperationsRequest,
) (*longrunning.ListOperationsResponse, error) {
	return &longrunning.ListOperationsResponse{}, nil
}

// GetOperation implements longrunning.OperationsServer.
func (s *Server) GetOperation(
	ctx context.Context,
	request *longrunning.GetOperationRequest,
) (*longrunning.Operation, error) {
	return nil, status.Error(codes.NotFound, "operation not found")
}

// DeleteOperation implements longrunning.OperationsServer.
func (s *Server) DeleteOperation(
	ctx context.Context,
	request *longrunning.DeleteOperationRequest,
) (*emptypb.Empty, error) {
	return nil, status.Error(codes.NotFound, "operation not found")
}

// CancelOperation implements longrunning.OperationsServer.
func (s *Server) CancelOperation(
	ctx context.Context,
	request *longrunning.CancelOperationRequest,
) (*emptypb.Empty, error) {
	return nil, status.Error(codes.NotFound, "operation not found")
}

// WaitOperation implements longrunning.OperationsServer.
func (s *Server) WaitOperation(
	ctx context.Context,
	request *longrunning.WaitOperationRequest,
) (*longrunning.Operation, error) {
	return nil, status.Error(codes.NotFound, "operation not found")
}
