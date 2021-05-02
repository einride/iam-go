package iamspanner

import (
	"context"
	"sort"

	"go.einride.tech/aip/pagination"
	"go.einride.tech/aip/validation"
	"google.golang.org/genproto/googleapis/iam/admin/v1"
	"google.golang.org/protobuf/proto"
)

// ListRoles implements admin.IAMServer.
func (s *Server) ListRoles(
	ctx context.Context,
	request *admin.ListRolesRequest,
) (*admin.ListRolesResponse, error) {
	var parsedRequest listRolesRequest
	if err := parsedRequest.parse(request); err != nil {
		return nil, err
	}
	return s.listRoles(ctx, &parsedRequest)
}

func (s *Server) listRoles(
	_ context.Context,
	request *listRolesRequest,
) (*admin.ListRolesResponse, error) {
	roles := make([]*admin.Role, 0, s.roles.Count())
	s.roles.RangeRoles(func(role *admin.Role) bool {
		roles = append(roles, role)
		return true
	})
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Name < roles[j].Name
	})
	response := admin.ListRolesResponse{
		Roles: make([]*admin.Role, 0, request.pageSize),
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
		case admin.RoleView_FULL:
			response.Roles = append(response.Roles, role)
		default:
			clone := proto.Clone(role).(*admin.Role)
			clone.IncludedPermissions = nil
			response.Roles = append(response.Roles, clone)
		}
	}
	return &response, nil
}

type listRolesRequest struct {
	pageSize  int32
	pageToken pagination.PageToken
	view      admin.RoleView
	request   *admin.ListRolesRequest
}

func (r *listRolesRequest) parse(request *admin.ListRolesRequest) error {
	const (
		defaultPageSize = 300
		maxPageSize     = 1_000
	)
	var v validation.MessageValidator
	r.request = request
	// parent = 1
	if request.Parent != "" {
		v.AddFieldViolation("parent", "unsupported field")
	}
	// page_size = 2
	if request.PageSize < 0 {
		v.AddFieldViolation("page_size", "must be >= 0")
	} else if request.PageSize == 0 {
		r.pageSize = defaultPageSize
	} else if request.PageSize > maxPageSize {
		r.pageSize = maxPageSize
	} else {
		r.pageSize = request.PageSize
	}
	// page_token = 3
	pageToken, err := pagination.ParsePageToken(request)
	if err != nil {
		v.AddFieldViolation("page_token", "invalid format")
	}
	r.pageToken = pageToken
	switch request.View {
	case admin.RoleView_BASIC, admin.RoleView_FULL:
		r.view = request.View
	default:
		v.AddFieldViolation("view", "unsupported value: %d", request.View.Number())
	}
	return v.Err()
}

func (r *listRolesRequest) nextPageToken() string {
	return r.pageToken.Next(r.request).String()
}
