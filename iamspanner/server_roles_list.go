package iamspanner

import (
	"context"
	"sort"

	"cloud.google.com/go/iam/admin/apiv1/adminpb"
	"go.einride.tech/aip/pagination"
	"go.einride.tech/aip/validation"
	"google.golang.org/protobuf/proto"
)

// ListRoles implements adminpb.IAMServer.
func (s *IAMServer) ListRoles(
	ctx context.Context,
	request *adminpb.ListRolesRequest,
) (*adminpb.ListRolesResponse, error) {
	var parsedRequest listRolesRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.listRoles(ctx, &parsedRequest)
}

func (s *IAMServer) listRoles(
	_ context.Context,
	request *listRolesRequest,
) (*adminpb.ListRolesResponse, error) {
	roles := make([]*adminpb.Role, 0, s.roles.Count())
	s.roles.RangeRoles(func(role *adminpb.Role) bool {
		roles = append(roles, role)
		return true
	})
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].GetName() < roles[j].GetName()
	})
	response := adminpb.ListRolesResponse{
		Roles: make([]*adminpb.Role, 0, request.pageSize),
	}
	from := int(request.pageToken.Offset)
	to := int(request.pageToken.Offset) + int(request.pageSize)
	if to >= len(roles) {
		to = len(roles)
	} else {
		response.NextPageToken = request.nextPageToken()
	}
	for _, role := range roles[from:to] {
		switch request.view {
		case adminpb.RoleView_FULL:
			response.Roles = append(response.Roles, role)
		default:
			clone := proto.Clone(role).(*adminpb.Role)
			clone.IncludedPermissions = nil
			response.Roles = append(response.Roles, clone)
		}
	}
	return &response, nil
}

type listRolesRequest struct {
	pageSize  int32
	pageToken pagination.PageToken
	view      adminpb.RoleView
	request   *adminpb.ListRolesRequest
}

func (r *listRolesRequest) parse(request *adminpb.ListRolesRequest) error {
	const (
		defaultPageSize = 300
		maxPageSize     = 1_000
	)
	var v validation.MessageValidator
	r.request = request
	// parent = 1
	if request.GetParent() != "" {
		v.AddFieldViolation("parent", "unsupported field")
	}
	// page_size = 2
	switch {
	case request.GetPageSize() < 0:
		v.AddFieldViolation("page_size", "must be >= 0")
	case request.GetPageSize() == 0:
		r.pageSize = defaultPageSize
	case request.GetPageSize() > maxPageSize:
		r.pageSize = maxPageSize
	default:
		r.pageSize = request.GetPageSize()
	}
	// page_token = 3
	pageToken, err := pagination.ParsePageToken(request)
	if err != nil {
		v.AddFieldViolation("page_token", "invalid format")
	}
	r.pageToken = pageToken
	switch request.GetView() {
	case adminpb.RoleView_BASIC, adminpb.RoleView_FULL:
		r.view = request.GetView()
	default:
		v.AddFieldViolation("view", "unsupported value: %d", request.GetView().Number())
	}
	return v.Err()
}

func (r *listRolesRequest) nextPageToken() string {
	return r.pageToken.Next(r.request).String()
}
