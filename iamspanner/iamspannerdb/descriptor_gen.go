package iamspannerdb

// Code generated by spanner-aip-go. DO NOT EDIT.

import (
	"cloud.google.com/go/spanner/spansql"
)

func Descriptor() DatabaseDescriptor {
	return &descriptor
}

var descriptor = databaseDescriptor{
	iamPolicyBindings: iamPolicyBindingsTableDescriptor{
		tableID: "iam_policy_bindings",
		resource: columnDescriptor{
			columnID:             "resource",
			columnType:           spansql.Type{Array: false, Base: 4, Len: 9223372036854775807},
			notNull:              true,
			allowCommitTimestamp: false,
		},
		bindingIndex: columnDescriptor{
			columnID:             "binding_index",
			columnType:           spansql.Type{Array: false, Base: 1, Len: 0},
			notNull:              true,
			allowCommitTimestamp: false,
		},
		role: columnDescriptor{
			columnID:             "role",
			columnType:           spansql.Type{Array: false, Base: 4, Len: 9223372036854775807},
			notNull:              true,
			allowCommitTimestamp: false,
		},
		memberIndex: columnDescriptor{
			columnID:             "member_index",
			columnType:           spansql.Type{Array: false, Base: 1, Len: 0},
			notNull:              true,
			allowCommitTimestamp: false,
		},
		member: columnDescriptor{
			columnID:             "member",
			columnType:           spansql.Type{Array: false, Base: 4, Len: 9223372036854775807},
			notNull:              true,
			allowCommitTimestamp: false,
		},
	},
	iamPolicyBindingsByMemberAndResource: iamPolicyBindingsByMemberAndResourceIndexDescriptor{
		indexID: "iam_policy_bindings_by_member_and_resource",
		member: columnDescriptor{
			columnID: "member",
		},
		resource: columnDescriptor{
			columnID: "resource",
		},
		role: columnDescriptor{
			columnID: "role",
		},
	},
	iamPolicyBindingsByMemberAndRole: iamPolicyBindingsByMemberAndRoleIndexDescriptor{
		indexID: "iam_policy_bindings_by_member_and_role",
		member: columnDescriptor{
			columnID: "member",
		},
		role: columnDescriptor{
			columnID: "role",
		},
		resource: columnDescriptor{
			columnID: "resource",
		},
	},
}

type DatabaseDescriptor interface {
	IamPolicyBindings() IamPolicyBindingsTableDescriptor
	IamPolicyBindingsByMemberAndResource() IamPolicyBindingsByMemberAndResourceIndexDescriptor
	IamPolicyBindingsByMemberAndRole() IamPolicyBindingsByMemberAndRoleIndexDescriptor
}

type databaseDescriptor struct {
	iamPolicyBindings                    iamPolicyBindingsTableDescriptor
	iamPolicyBindingsByMemberAndResource iamPolicyBindingsByMemberAndResourceIndexDescriptor
	iamPolicyBindingsByMemberAndRole     iamPolicyBindingsByMemberAndRoleIndexDescriptor
}

func (d *databaseDescriptor) IamPolicyBindings() IamPolicyBindingsTableDescriptor {
	return &d.iamPolicyBindings
}

func (d *databaseDescriptor) IamPolicyBindingsByMemberAndResource() IamPolicyBindingsByMemberAndResourceIndexDescriptor {
	return &d.iamPolicyBindingsByMemberAndResource
}

func (d *databaseDescriptor) IamPolicyBindingsByMemberAndRole() IamPolicyBindingsByMemberAndRoleIndexDescriptor {
	return &d.iamPolicyBindingsByMemberAndRole
}

type IamPolicyBindingsTableDescriptor interface {
	TableName() string
	TableID() spansql.ID
	ColumnNames() []string
	ColumnIDs() []spansql.ID
	ColumnExprs() []spansql.Expr
	Resource() ColumnDescriptor
	BindingIndex() ColumnDescriptor
	Role() ColumnDescriptor
	MemberIndex() ColumnDescriptor
	Member() ColumnDescriptor
}

type iamPolicyBindingsTableDescriptor struct {
	tableID      spansql.ID
	resource     columnDescriptor
	bindingIndex columnDescriptor
	role         columnDescriptor
	memberIndex  columnDescriptor
	member       columnDescriptor
}

func (d *iamPolicyBindingsTableDescriptor) TableName() string {
	return string(d.tableID)
}

func (d *iamPolicyBindingsTableDescriptor) TableID() spansql.ID {
	return d.tableID
}

func (d *iamPolicyBindingsTableDescriptor) ColumnNames() []string {
	return []string{
		"resource",
		"binding_index",
		"role",
		"member_index",
		"member",
	}
}

func (d *iamPolicyBindingsTableDescriptor) ColumnIDs() []spansql.ID {
	return []spansql.ID{
		"resource",
		"binding_index",
		"role",
		"member_index",
		"member",
	}
}

func (d *iamPolicyBindingsTableDescriptor) ColumnExprs() []spansql.Expr {
	return []spansql.Expr{
		spansql.ID("resource"),
		spansql.ID("binding_index"),
		spansql.ID("role"),
		spansql.ID("member_index"),
		spansql.ID("member"),
	}
}

func (d *iamPolicyBindingsTableDescriptor) Resource() ColumnDescriptor {
	return &d.resource
}

func (d *iamPolicyBindingsTableDescriptor) BindingIndex() ColumnDescriptor {
	return &d.bindingIndex
}

func (d *iamPolicyBindingsTableDescriptor) Role() ColumnDescriptor {
	return &d.role
}

func (d *iamPolicyBindingsTableDescriptor) MemberIndex() ColumnDescriptor {
	return &d.memberIndex
}

func (d *iamPolicyBindingsTableDescriptor) Member() ColumnDescriptor {
	return &d.member
}

type IamPolicyBindingsByMemberAndResourceIndexDescriptor interface {
	IndexName() string
	IndexID() spansql.ID
	ColumnNames() []string
	ColumnIDs() []spansql.ID
	ColumnExprs() []spansql.Expr
	Member() ColumnDescriptor
	Resource() ColumnDescriptor
	Role() ColumnDescriptor
}

type iamPolicyBindingsByMemberAndResourceIndexDescriptor struct {
	indexID  spansql.ID
	member   columnDescriptor
	resource columnDescriptor
	role     columnDescriptor
}

func (d *iamPolicyBindingsByMemberAndResourceIndexDescriptor) IndexName() string {
	return string(d.indexID)
}

func (d *iamPolicyBindingsByMemberAndResourceIndexDescriptor) IndexID() spansql.ID {
	return d.indexID
}

func (d *iamPolicyBindingsByMemberAndResourceIndexDescriptor) ColumnNames() []string {
	return []string{
		"member",
		"resource",
		"role",
	}
}

func (d *iamPolicyBindingsByMemberAndResourceIndexDescriptor) ColumnIDs() []spansql.ID {
	return []spansql.ID{
		"member",
		"resource",
		"role",
	}
}

func (d *iamPolicyBindingsByMemberAndResourceIndexDescriptor) ColumnExprs() []spansql.Expr {
	return []spansql.Expr{
		spansql.ID("member"),
		spansql.ID("resource"),
		spansql.ID("role"),
	}
}

func (d *iamPolicyBindingsByMemberAndResourceIndexDescriptor) Member() ColumnDescriptor {
	return &d.member
}

func (d *iamPolicyBindingsByMemberAndResourceIndexDescriptor) Resource() ColumnDescriptor {
	return &d.resource
}

func (d *iamPolicyBindingsByMemberAndResourceIndexDescriptor) Role() ColumnDescriptor {
	return &d.role
}

type IamPolicyBindingsByMemberAndRoleIndexDescriptor interface {
	IndexName() string
	IndexID() spansql.ID
	ColumnNames() []string
	ColumnIDs() []spansql.ID
	ColumnExprs() []spansql.Expr
	Member() ColumnDescriptor
	Role() ColumnDescriptor
	Resource() ColumnDescriptor
}

type iamPolicyBindingsByMemberAndRoleIndexDescriptor struct {
	indexID  spansql.ID
	member   columnDescriptor
	role     columnDescriptor
	resource columnDescriptor
}

func (d *iamPolicyBindingsByMemberAndRoleIndexDescriptor) IndexName() string {
	return string(d.indexID)
}

func (d *iamPolicyBindingsByMemberAndRoleIndexDescriptor) IndexID() spansql.ID {
	return d.indexID
}

func (d *iamPolicyBindingsByMemberAndRoleIndexDescriptor) ColumnNames() []string {
	return []string{
		"member",
		"role",
		"resource",
	}
}

func (d *iamPolicyBindingsByMemberAndRoleIndexDescriptor) ColumnIDs() []spansql.ID {
	return []spansql.ID{
		"member",
		"role",
		"resource",
	}
}

func (d *iamPolicyBindingsByMemberAndRoleIndexDescriptor) ColumnExprs() []spansql.Expr {
	return []spansql.Expr{
		spansql.ID("member"),
		spansql.ID("role"),
		spansql.ID("resource"),
	}
}

func (d *iamPolicyBindingsByMemberAndRoleIndexDescriptor) Member() ColumnDescriptor {
	return &d.member
}

func (d *iamPolicyBindingsByMemberAndRoleIndexDescriptor) Role() ColumnDescriptor {
	return &d.role
}

func (d *iamPolicyBindingsByMemberAndRoleIndexDescriptor) Resource() ColumnDescriptor {
	return &d.resource
}

type ColumnDescriptor interface {
	ColumnID() spansql.ID
	ColumnName() string
	ColumnType() spansql.Type
	NotNull() bool
	AllowCommitTimestamp() bool
}

type columnDescriptor struct {
	columnID             spansql.ID
	columnType           spansql.Type
	notNull              bool
	allowCommitTimestamp bool
}

func (d *columnDescriptor) ColumnName() string {
	return string(d.columnID)
}

func (d *columnDescriptor) ColumnID() spansql.ID {
	return d.columnID
}

func (d *columnDescriptor) ColumnType() spansql.Type {
	return d.columnType
}

func (d *columnDescriptor) ColumnExpr() spansql.Expr {
	return d.columnID
}

func (d *columnDescriptor) NotNull() bool {
	return d.notNull
}

func (d *columnDescriptor) AllowCommitTimestamp() bool {
	return d.allowCommitTimestamp
}
