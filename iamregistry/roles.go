package iamregistry

import (
	"fmt"

	"go.einride.tech/iam/iamannotations"
	"go.einride.tech/iam/iamrole"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/genproto/googleapis/iam/admin/v1"
)

// Roles are a set of roles.
type Roles struct {
	roles map[string]*admin.Role
}

// NewRoles creates a set of Roles from a pre-defined roles annotation.
func NewRoles(roles ...*admin.Role) (*Roles, error) {
	if err := iamannotations.ValidatePredefinedRoles(&iamv1.PredefinedRoles{Role: roles}); err != nil {
		return nil, fmt.Errorf("new roles registry: %w", err)
	}
	result := Roles{
		roles: make(map[string]*admin.Role, len(roles)),
	}
	for _, role := range roles {
		result.roles[role.Name] = role
	}
	return &result, nil
}

// Count returns a count of the roles.
func (r *Roles) Count() int {
	return len(r.roles)
}

// RoleHasPermission checks whether the role with the provided name has the provided permission.
func (r *Roles) RoleHasPermission(name, permission string) bool {
	role, ok := r.FindRoleByName(name)
	if !ok {
		return false
	}
	return iamrole.HasPermission(role, permission)
}

// RangeRoles iterates over all registered roles while f returns true.
// The iteration order is undefined.
func (r *Roles) RangeRoles(fn func(*admin.Role) bool) {
	for _, role := range r.roles {
		if !fn(role) {
			break
		}
	}
}

// RangeRolesByPermission iterates over all registered roles with the provided permission while f returns true.
// The iteration order is undefined, and permissions with wildcards are not allowed.
func (r *Roles) RangeRolesByPermission(permission string, fn func(*admin.Role) bool) {
	for _, role := range r.roles {
		if iamrole.HasPermission(role, permission) {
			if !fn(role) {
				break
			}
		}
	}
}

// FindRoleByName looks up a role by resource name.
func (r *Roles) FindRoleByName(name string) (*admin.Role, bool) {
	role, ok := r.roles[name]
	return role, ok
}
