package iamspanner

import (
	"context"

	"go.einride.tech/aip/resourcename"
	"go.einride.tech/aip/validation"
	"google.golang.org/genproto/googleapis/iam/admin/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRole implements admin.IAMServer.
func (s *IAMServer) GetRole(
	ctx context.Context,
	request *admin.GetRoleRequest,
) (*admin.Role, error) {
	var parsedRequest getRoleRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.getRole(ctx, &parsedRequest)
}

// GetRole implements admin.IAMServer.
func (s *IAMServer) getRole(
	_ context.Context,
	request *getRoleRequest,
) (*admin.Role, error) {
	role, ok := s.roles.FindRoleByName(request.name)
	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}
	return role, nil
}

type getRoleRequest struct {
	name string
}

func (r *getRoleRequest) parse(request *admin.GetRoleRequest) error {
	var v validation.MessageValidator
	// name = 1
	switch {
	case request.Name == "":
		v.AddFieldViolation("name", "required field")
	case resourcename.ContainsWildcard(request.Name):
		v.AddFieldViolation("name", "must not contain wildcards")
	case !resourcename.Match("roles/{role}", request.Name):
		v.AddFieldViolation("name", "invalid format")
	default:
		r.name = request.Name
	}
	return v.Err()
}
