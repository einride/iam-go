// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.15.2
// source: einride/iam/v1/annotations.proto

package iamv1

import (
	annotations "google.golang.org/genproto/googleapis/api/annotations"
	v1 "google.golang.org/genproto/googleapis/iam/admin/v1"
	expr "google.golang.org/genproto/googleapis/type/expr"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// A list of predefined roles.
type PredefinedRoles struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Role []*v1.Role `protobuf:"bytes,1,rep,name=role,proto3" json:"role,omitempty"`
}

func (x *PredefinedRoles) Reset() {
	*x = PredefinedRoles{}
	if protoimpl.UnsafeEnabled {
		mi := &file_einride_iam_v1_annotations_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PredefinedRoles) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PredefinedRoles) ProtoMessage() {}

func (x *PredefinedRoles) ProtoReflect() protoreflect.Message {
	mi := &file_einride_iam_v1_annotations_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PredefinedRoles.ProtoReflect.Descriptor instead.
func (*PredefinedRoles) Descriptor() ([]byte, []int) {
	return file_einride_iam_v1_annotations_proto_rawDescGZIP(), []int{0}
}

func (x *PredefinedRoles) GetRole() []*v1.Role {
	if x != nil {
		return x.Role
	}
	return nil
}

// Method authorization options.
type MethodAuthorizationOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Permission to use for authorization.
	//
	// Types that are assignable to Permissions:
	//	*MethodAuthorizationOptions_Permission
	//	*MethodAuthorizationOptions_ResourcePermissions
	Permissions isMethodAuthorizationOptions_Permissions `protobuf_oneof:"permissions"`
	// Strategy that decides if the request is authorized.
	//
	// Types that are assignable to Strategy:
	//	*MethodAuthorizationOptions_Before
	//	*MethodAuthorizationOptions_After
	//	*MethodAuthorizationOptions_Custom
	//	*MethodAuthorizationOptions_None
	Strategy isMethodAuthorizationOptions_Strategy `protobuf_oneof:"strategy"`
}

func (x *MethodAuthorizationOptions) Reset() {
	*x = MethodAuthorizationOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_einride_iam_v1_annotations_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MethodAuthorizationOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MethodAuthorizationOptions) ProtoMessage() {}

func (x *MethodAuthorizationOptions) ProtoReflect() protoreflect.Message {
	mi := &file_einride_iam_v1_annotations_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MethodAuthorizationOptions.ProtoReflect.Descriptor instead.
func (*MethodAuthorizationOptions) Descriptor() ([]byte, []int) {
	return file_einride_iam_v1_annotations_proto_rawDescGZIP(), []int{1}
}

func (m *MethodAuthorizationOptions) GetPermissions() isMethodAuthorizationOptions_Permissions {
	if m != nil {
		return m.Permissions
	}
	return nil
}

func (x *MethodAuthorizationOptions) GetPermission() string {
	if x, ok := x.GetPermissions().(*MethodAuthorizationOptions_Permission); ok {
		return x.Permission
	}
	return ""
}

func (x *MethodAuthorizationOptions) GetResourcePermissions() *ResourcePermissions {
	if x, ok := x.GetPermissions().(*MethodAuthorizationOptions_ResourcePermissions); ok {
		return x.ResourcePermissions
	}
	return nil
}

func (m *MethodAuthorizationOptions) GetStrategy() isMethodAuthorizationOptions_Strategy {
	if m != nil {
		return m.Strategy
	}
	return nil
}

func (x *MethodAuthorizationOptions) GetBefore() *expr.Expr {
	if x, ok := x.GetStrategy().(*MethodAuthorizationOptions_Before); ok {
		return x.Before
	}
	return nil
}

func (x *MethodAuthorizationOptions) GetAfter() *expr.Expr {
	if x, ok := x.GetStrategy().(*MethodAuthorizationOptions_After); ok {
		return x.After
	}
	return nil
}

func (x *MethodAuthorizationOptions) GetCustom() bool {
	if x, ok := x.GetStrategy().(*MethodAuthorizationOptions_Custom); ok {
		return x.Custom
	}
	return false
}

func (x *MethodAuthorizationOptions) GetNone() bool {
	if x, ok := x.GetStrategy().(*MethodAuthorizationOptions_None); ok {
		return x.None
	}
	return false
}

type isMethodAuthorizationOptions_Permissions interface {
	isMethodAuthorizationOptions_Permissions()
}

type MethodAuthorizationOptions_Permission struct {
	// A single permission used by the method.
	Permission string `protobuf:"bytes,1,opt,name=permission,proto3,oneof"`
}

type MethodAuthorizationOptions_ResourcePermissions struct {
	// Resource permissions used by the method.
	ResourcePermissions *ResourcePermissions `protobuf:"bytes,2,opt,name=resource_permissions,json=resourcePermissions,proto3,oneof"`
}

func (*MethodAuthorizationOptions_Permission) isMethodAuthorizationOptions_Permissions() {}

func (*MethodAuthorizationOptions_ResourcePermissions) isMethodAuthorizationOptions_Permissions() {}

type isMethodAuthorizationOptions_Strategy interface {
	isMethodAuthorizationOptions_Strategy()
}

type MethodAuthorizationOptions_Before struct {
	// Expression that decides before the request if the caller is authorized.
	Before *expr.Expr `protobuf:"bytes,3,opt,name=before,proto3,oneof"`
}

type MethodAuthorizationOptions_After struct {
	// Expression that decides after the request if the caller is authorized.
	After *expr.Expr `protobuf:"bytes,4,opt,name=after,proto3,oneof"`
}

type MethodAuthorizationOptions_Custom struct {
	// A flag indicating if the method requires custom-implemented authorization.
	Custom bool `protobuf:"varint,5,opt,name=custom,proto3,oneof"`
}

type MethodAuthorizationOptions_None struct {
	// A flag indicating if the method requires no authorization.
	None bool `protobuf:"varint,6,opt,name=none,proto3,oneof"`
}

func (*MethodAuthorizationOptions_Before) isMethodAuthorizationOptions_Strategy() {}

func (*MethodAuthorizationOptions_After) isMethodAuthorizationOptions_Strategy() {}

func (*MethodAuthorizationOptions_Custom) isMethodAuthorizationOptions_Strategy() {}

func (*MethodAuthorizationOptions_None) isMethodAuthorizationOptions_Strategy() {}

// Resource permissions.
type ResourcePermissions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The resource permissions.
	ResourcePermission []*ResourcePermission `protobuf:"bytes,1,rep,name=resource_permission,json=resourcePermission,proto3" json:"resource_permission,omitempty"`
}

func (x *ResourcePermissions) Reset() {
	*x = ResourcePermissions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_einride_iam_v1_annotations_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourcePermissions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourcePermissions) ProtoMessage() {}

func (x *ResourcePermissions) ProtoReflect() protoreflect.Message {
	mi := &file_einride_iam_v1_annotations_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourcePermissions.ProtoReflect.Descriptor instead.
func (*ResourcePermissions) Descriptor() ([]byte, []int) {
	return file_einride_iam_v1_annotations_proto_rawDescGZIP(), []int{2}
}

func (x *ResourcePermissions) GetResourcePermission() []*ResourcePermission {
	if x != nil {
		return x.ResourcePermission
	}
	return nil
}

// A resource type and a permission.
type ResourcePermission struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The resource.
	// When used for authorization method options, only the type must be provided.
	Resource *annotations.ResourceDescriptor `protobuf:"bytes,1,opt,name=resource,proto3" json:"resource,omitempty"`
	// The permission.
	Permission string `protobuf:"bytes,2,opt,name=permission,proto3" json:"permission,omitempty"`
}

func (x *ResourcePermission) Reset() {
	*x = ResourcePermission{}
	if protoimpl.UnsafeEnabled {
		mi := &file_einride_iam_v1_annotations_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourcePermission) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourcePermission) ProtoMessage() {}

func (x *ResourcePermission) ProtoReflect() protoreflect.Message {
	mi := &file_einride_iam_v1_annotations_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourcePermission.ProtoReflect.Descriptor instead.
func (*ResourcePermission) Descriptor() ([]byte, []int) {
	return file_einride_iam_v1_annotations_proto_rawDescGZIP(), []int{3}
}

func (x *ResourcePermission) GetResource() *annotations.ResourceDescriptor {
	if x != nil {
		return x.Resource
	}
	return nil
}

func (x *ResourcePermission) GetPermission() string {
	if x != nil {
		return x.Permission
	}
	return ""
}

// Long-running operations permissions.
type LongRunningOperationsAuthorizationOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The long-running operation permissions.
	OperationPermissions []*LongRunningOperationPermissions `protobuf:"bytes,1,rep,name=operation_permissions,json=operationPermissions,proto3" json:"operation_permissions,omitempty"`
	// Strategy that decides if the request is authorized.
	//
	// Types that are assignable to Strategy:
	//	*LongRunningOperationsAuthorizationOptions_Before
	//	*LongRunningOperationsAuthorizationOptions_Custom
	//	*LongRunningOperationsAuthorizationOptions_None
	Strategy isLongRunningOperationsAuthorizationOptions_Strategy `protobuf_oneof:"strategy"`
}

func (x *LongRunningOperationsAuthorizationOptions) Reset() {
	*x = LongRunningOperationsAuthorizationOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_einride_iam_v1_annotations_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LongRunningOperationsAuthorizationOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LongRunningOperationsAuthorizationOptions) ProtoMessage() {}

func (x *LongRunningOperationsAuthorizationOptions) ProtoReflect() protoreflect.Message {
	mi := &file_einride_iam_v1_annotations_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LongRunningOperationsAuthorizationOptions.ProtoReflect.Descriptor instead.
func (*LongRunningOperationsAuthorizationOptions) Descriptor() ([]byte, []int) {
	return file_einride_iam_v1_annotations_proto_rawDescGZIP(), []int{4}
}

func (x *LongRunningOperationsAuthorizationOptions) GetOperationPermissions() []*LongRunningOperationPermissions {
	if x != nil {
		return x.OperationPermissions
	}
	return nil
}

func (m *LongRunningOperationsAuthorizationOptions) GetStrategy() isLongRunningOperationsAuthorizationOptions_Strategy {
	if m != nil {
		return m.Strategy
	}
	return nil
}

func (x *LongRunningOperationsAuthorizationOptions) GetBefore() bool {
	if x, ok := x.GetStrategy().(*LongRunningOperationsAuthorizationOptions_Before); ok {
		return x.Before
	}
	return false
}

func (x *LongRunningOperationsAuthorizationOptions) GetCustom() bool {
	if x, ok := x.GetStrategy().(*LongRunningOperationsAuthorizationOptions_Custom); ok {
		return x.Custom
	}
	return false
}

func (x *LongRunningOperationsAuthorizationOptions) GetNone() bool {
	if x, ok := x.GetStrategy().(*LongRunningOperationsAuthorizationOptions_None); ok {
		return x.None
	}
	return false
}

type isLongRunningOperationsAuthorizationOptions_Strategy interface {
	isLongRunningOperationsAuthorizationOptions_Strategy()
}

type LongRunningOperationsAuthorizationOptions_Before struct {
	// A flag indicating if a standard authorization checked is performed before the request.
	Before bool `protobuf:"varint,3,opt,name=before,proto3,oneof"`
}

type LongRunningOperationsAuthorizationOptions_Custom struct {
	// A flag indicating if custom-implemented authorization is required.
	Custom bool `protobuf:"varint,4,opt,name=custom,proto3,oneof"`
}

type LongRunningOperationsAuthorizationOptions_None struct {
	// A flag indicating if no authorization is required.
	None bool `protobuf:"varint,5,opt,name=none,proto3,oneof"`
}

func (*LongRunningOperationsAuthorizationOptions_Before) isLongRunningOperationsAuthorizationOptions_Strategy() {
}

func (*LongRunningOperationsAuthorizationOptions_Custom) isLongRunningOperationsAuthorizationOptions_Strategy() {
}

func (*LongRunningOperationsAuthorizationOptions_None) isLongRunningOperationsAuthorizationOptions_Strategy() {
}

// Permissions for a long-running operation.
type LongRunningOperationPermissions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The long-running operation resource. The type field is required.
	Operation *annotations.ResourceDescriptor `protobuf:"bytes,1,opt,name=operation,proto3" json:"operation,omitempty"`
	// Permission for listing operations.
	List string `protobuf:"bytes,2,opt,name=list,proto3" json:"list,omitempty"`
	// Permission for getting an operation.
	Get string `protobuf:"bytes,3,opt,name=get,proto3" json:"get,omitempty"`
	// Permission for cancelling an operation.
	Cancel string `protobuf:"bytes,4,opt,name=cancel,proto3" json:"cancel,omitempty"`
	// Permission for deleting an operation.
	Delete string `protobuf:"bytes,5,opt,name=delete,proto3" json:"delete,omitempty"`
	// Permission for waiting on an operation.
	Wait string `protobuf:"bytes,6,opt,name=wait,proto3" json:"wait,omitempty"`
}

func (x *LongRunningOperationPermissions) Reset() {
	*x = LongRunningOperationPermissions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_einride_iam_v1_annotations_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LongRunningOperationPermissions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LongRunningOperationPermissions) ProtoMessage() {}

func (x *LongRunningOperationPermissions) ProtoReflect() protoreflect.Message {
	mi := &file_einride_iam_v1_annotations_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LongRunningOperationPermissions.ProtoReflect.Descriptor instead.
func (*LongRunningOperationPermissions) Descriptor() ([]byte, []int) {
	return file_einride_iam_v1_annotations_proto_rawDescGZIP(), []int{5}
}

func (x *LongRunningOperationPermissions) GetOperation() *annotations.ResourceDescriptor {
	if x != nil {
		return x.Operation
	}
	return nil
}

func (x *LongRunningOperationPermissions) GetList() string {
	if x != nil {
		return x.List
	}
	return ""
}

func (x *LongRunningOperationPermissions) GetGet() string {
	if x != nil {
		return x.Get
	}
	return ""
}

func (x *LongRunningOperationPermissions) GetCancel() string {
	if x != nil {
		return x.Cancel
	}
	return ""
}

func (x *LongRunningOperationPermissions) GetDelete() string {
	if x != nil {
		return x.Delete
	}
	return ""
}

func (x *LongRunningOperationPermissions) GetWait() string {
	if x != nil {
		return x.Wait
	}
	return ""
}

var file_einride_iam_v1_annotations_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.ServiceOptions)(nil),
		ExtensionType: (*PredefinedRoles)(nil),
		Field:         201601,
		Name:          "einride.iam.v1.predefined_roles",
		Tag:           "bytes,201601,opt,name=predefined_roles",
		Filename:      "einride/iam/v1/annotations.proto",
	},
	{
		ExtendedType:  (*descriptorpb.ServiceOptions)(nil),
		ExtensionType: (*LongRunningOperationsAuthorizationOptions)(nil),
		Field:         201602,
		Name:          "einride.iam.v1.long_running_operations_authorization",
		Tag:           "bytes,201602,opt,name=long_running_operations_authorization",
		Filename:      "einride/iam/v1/annotations.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*MethodAuthorizationOptions)(nil),
		Field:         201600,
		Name:          "einride.iam.v1.method_authorization",
		Tag:           "bytes,201600,opt,name=method_authorization",
		Filename:      "einride/iam/v1/annotations.proto",
	},
}

// Extension fields to descriptorpb.ServiceOptions.
var (
	// Predefined roles for the service.
	//
	// optional einride.iam.v1.PredefinedRoles predefined_roles = 201601;
	E_PredefinedRoles = &file_einride_iam_v1_annotations_proto_extTypes[0]
	// Long-running operations authorization for the service.
	//
	// optional einride.iam.v1.LongRunningOperationsAuthorizationOptions long_running_operations_authorization = 201602;
	E_LongRunningOperationsAuthorization = &file_einride_iam_v1_annotations_proto_extTypes[1]
)

// Extension fields to descriptorpb.MethodOptions.
var (
	// Method authorization options.
	//
	// optional einride.iam.v1.MethodAuthorizationOptions method_authorization = 201600;
	E_MethodAuthorization = &file_einride_iam_v1_annotations_proto_extTypes[2]
)

var File_einride_iam_v1_annotations_proto protoreflect.FileDescriptor

var file_einride_iam_v1_annotations_proto_rawDesc = []byte{
	0x0a, 0x20, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x76, 0x31,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0e, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e,
	0x76, 0x31, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f,
	0x76, 0x31, 0x2f, 0x69, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x65, 0x78, 0x70, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x40, 0x0a, 0x0f, 0x50, 0x72, 0x65, 0x64, 0x65, 0x66,
	0x69, 0x6e, 0x65, 0x64, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x12, 0x2d, 0x0a, 0x04, 0x72, 0x6f, 0x6c,
	0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x6f,
	0x6c, 0x65, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x22, 0xbb, 0x02, 0x0a, 0x1a, 0x4d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x20, 0x0a, 0x0a, 0x70, 0x65, 0x72, 0x6d, 0x69,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0a, 0x70,
	0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x58, 0x0a, 0x14, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64,
	0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x48, 0x00, 0x52, 0x13,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x2b, 0x0a, 0x06, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x74, 0x79, 0x70,
	0x65, 0x2e, 0x45, 0x78, 0x70, 0x72, 0x48, 0x01, 0x52, 0x06, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65,
	0x12, 0x29, 0x0a, 0x05, 0x61, 0x66, 0x74, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x11, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x45, 0x78,
	0x70, 0x72, 0x48, 0x01, 0x52, 0x05, 0x61, 0x66, 0x74, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x06, 0x63,
	0x75, 0x73, 0x74, 0x6f, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x48, 0x01, 0x52, 0x06, 0x63,
	0x75, 0x73, 0x74, 0x6f, 0x6d, 0x12, 0x14, 0x0a, 0x04, 0x6e, 0x6f, 0x6e, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x08, 0x48, 0x01, 0x52, 0x04, 0x6e, 0x6f, 0x6e, 0x65, 0x42, 0x0d, 0x0a, 0x0b, 0x70,
	0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x0a, 0x0a, 0x08, 0x73, 0x74,
	0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x22, 0x6a, 0x0a, 0x13, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x53, 0x0a,
	0x13, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x65, 0x69, 0x6e,
	0x72, 0x69, 0x64, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x12,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x22, 0x70, 0x0a, 0x12, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x65,
	0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x22, 0xe7, 0x01, 0x0a, 0x29, 0x4c, 0x6f, 0x6e, 0x67, 0x52, 0x75, 0x6e,
	0x6e, 0x69, 0x6e, 0x67, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x41, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x64, 0x0a, 0x15, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x2f, 0x2e, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x6e, 0x67, 0x52, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x4f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x14, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x65, 0x72,
	0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x18, 0x0a, 0x06, 0x62, 0x65, 0x66, 0x6f,
	0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x06, 0x62, 0x65, 0x66, 0x6f,
	0x72, 0x65, 0x12, 0x18, 0x0a, 0x06, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x08, 0x48, 0x00, 0x52, 0x06, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x12, 0x14, 0x0a, 0x04,
	0x6e, 0x6f, 0x6e, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x04, 0x6e, 0x6f,
	0x6e, 0x65, 0x42, 0x0a, 0x0a, 0x08, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x22, 0xc9,
	0x01, 0x0a, 0x1f, 0x4c, 0x6f, 0x6e, 0x67, 0x52, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x4f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x3c, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x44, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x6f, 0x72, 0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x12, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6c, 0x69, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x67, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x67, 0x65, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x61, 0x6e, 0x63, 0x65, 0x6c,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x12, 0x16,
	0x0a, 0x06, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x77, 0x61, 0x69, 0x74, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x77, 0x61, 0x69, 0x74, 0x3a, 0x6d, 0x0a, 0x10, 0x70, 0x72,
	0x65, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x5f, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x12, 0x1f,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x81, 0xa7, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64,
	0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x65, 0x64, 0x65, 0x66, 0x69,
	0x6e, 0x65, 0x64, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x52, 0x0f, 0x70, 0x72, 0x65, 0x64, 0x65, 0x66,
	0x69, 0x6e, 0x65, 0x64, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x3a, 0xaf, 0x01, 0x0a, 0x25, 0x6c, 0x6f,
	0x6e, 0x67, 0x5f, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x6f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x5f, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0x82, 0xa7, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x39, 0x2e, 0x65,
	0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f,
	0x6e, 0x67, 0x52, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x22, 0x6c, 0x6f, 0x6e, 0x67, 0x52, 0x75, 0x6e,
	0x6e, 0x69, 0x6e, 0x67, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x41, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x7f, 0x0a, 0x14, 0x6d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x5f, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x80, 0xa7, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x65, 0x69,
	0x6e, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x13, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x41,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x75, 0x0a, 0x13,
	0x74, 0x65, 0x63, 0x68, 0x2e, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x69, 0x61, 0x6d,
	0x2e, 0x76, 0x31, 0x42, 0x10, 0x41, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x32, 0x67, 0x6f, 0x2e, 0x65, 0x69, 0x6e, 0x72,
	0x69, 0x64, 0x65, 0x2e, 0x74, 0x65, 0x63, 0x68, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2f, 0x69,
	0x61, 0x6d, 0x2f, 0x76, 0x31, 0x3b, 0x69, 0x61, 0x6d, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x49, 0x58,
	0x58, 0xaa, 0x02, 0x06, 0x49, 0x61, 0x6d, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x06, 0x49, 0x61, 0x6d,
	0x5c, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_einride_iam_v1_annotations_proto_rawDescOnce sync.Once
	file_einride_iam_v1_annotations_proto_rawDescData = file_einride_iam_v1_annotations_proto_rawDesc
)

func file_einride_iam_v1_annotations_proto_rawDescGZIP() []byte {
	file_einride_iam_v1_annotations_proto_rawDescOnce.Do(func() {
		file_einride_iam_v1_annotations_proto_rawDescData = protoimpl.X.CompressGZIP(file_einride_iam_v1_annotations_proto_rawDescData)
	})
	return file_einride_iam_v1_annotations_proto_rawDescData
}

var file_einride_iam_v1_annotations_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_einride_iam_v1_annotations_proto_goTypes = []interface{}{
	(*PredefinedRoles)(nil),                           // 0: einride.iam.v1.PredefinedRoles
	(*MethodAuthorizationOptions)(nil),                // 1: einride.iam.v1.MethodAuthorizationOptions
	(*ResourcePermissions)(nil),                       // 2: einride.iam.v1.ResourcePermissions
	(*ResourcePermission)(nil),                        // 3: einride.iam.v1.ResourcePermission
	(*LongRunningOperationsAuthorizationOptions)(nil), // 4: einride.iam.v1.LongRunningOperationsAuthorizationOptions
	(*LongRunningOperationPermissions)(nil),           // 5: einride.iam.v1.LongRunningOperationPermissions
	(*v1.Role)(nil),                                   // 6: google.iam.admin.v1.Role
	(*expr.Expr)(nil),                                 // 7: google.type.Expr
	(*annotations.ResourceDescriptor)(nil),            // 8: google.api.ResourceDescriptor
	(*descriptorpb.ServiceOptions)(nil),               // 9: google.protobuf.ServiceOptions
	(*descriptorpb.MethodOptions)(nil),                // 10: google.protobuf.MethodOptions
}
var file_einride_iam_v1_annotations_proto_depIdxs = []int32{
	6,  // 0: einride.iam.v1.PredefinedRoles.role:type_name -> google.iam.admin.v1.Role
	2,  // 1: einride.iam.v1.MethodAuthorizationOptions.resource_permissions:type_name -> einride.iam.v1.ResourcePermissions
	7,  // 2: einride.iam.v1.MethodAuthorizationOptions.before:type_name -> google.type.Expr
	7,  // 3: einride.iam.v1.MethodAuthorizationOptions.after:type_name -> google.type.Expr
	3,  // 4: einride.iam.v1.ResourcePermissions.resource_permission:type_name -> einride.iam.v1.ResourcePermission
	8,  // 5: einride.iam.v1.ResourcePermission.resource:type_name -> google.api.ResourceDescriptor
	5,  // 6: einride.iam.v1.LongRunningOperationsAuthorizationOptions.operation_permissions:type_name -> einride.iam.v1.LongRunningOperationPermissions
	8,  // 7: einride.iam.v1.LongRunningOperationPermissions.operation:type_name -> google.api.ResourceDescriptor
	9,  // 8: einride.iam.v1.predefined_roles:extendee -> google.protobuf.ServiceOptions
	9,  // 9: einride.iam.v1.long_running_operations_authorization:extendee -> google.protobuf.ServiceOptions
	10, // 10: einride.iam.v1.method_authorization:extendee -> google.protobuf.MethodOptions
	0,  // 11: einride.iam.v1.predefined_roles:type_name -> einride.iam.v1.PredefinedRoles
	4,  // 12: einride.iam.v1.long_running_operations_authorization:type_name -> einride.iam.v1.LongRunningOperationsAuthorizationOptions
	1,  // 13: einride.iam.v1.method_authorization:type_name -> einride.iam.v1.MethodAuthorizationOptions
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	11, // [11:14] is the sub-list for extension type_name
	8,  // [8:11] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_einride_iam_v1_annotations_proto_init() }
func file_einride_iam_v1_annotations_proto_init() {
	if File_einride_iam_v1_annotations_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_einride_iam_v1_annotations_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PredefinedRoles); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_einride_iam_v1_annotations_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MethodAuthorizationOptions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_einride_iam_v1_annotations_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourcePermissions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_einride_iam_v1_annotations_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourcePermission); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_einride_iam_v1_annotations_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LongRunningOperationsAuthorizationOptions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_einride_iam_v1_annotations_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LongRunningOperationPermissions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_einride_iam_v1_annotations_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*MethodAuthorizationOptions_Permission)(nil),
		(*MethodAuthorizationOptions_ResourcePermissions)(nil),
		(*MethodAuthorizationOptions_Before)(nil),
		(*MethodAuthorizationOptions_After)(nil),
		(*MethodAuthorizationOptions_Custom)(nil),
		(*MethodAuthorizationOptions_None)(nil),
	}
	file_einride_iam_v1_annotations_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*LongRunningOperationsAuthorizationOptions_Before)(nil),
		(*LongRunningOperationsAuthorizationOptions_Custom)(nil),
		(*LongRunningOperationsAuthorizationOptions_None)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_einride_iam_v1_annotations_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 3,
			NumServices:   0,
		},
		GoTypes:           file_einride_iam_v1_annotations_proto_goTypes,
		DependencyIndexes: file_einride_iam_v1_annotations_proto_depIdxs,
		MessageInfos:      file_einride_iam_v1_annotations_proto_msgTypes,
		ExtensionInfos:    file_einride_iam_v1_annotations_proto_extTypes,
	}.Build()
	File_einride_iam_v1_annotations_proto = out.File
	file_einride_iam_v1_annotations_proto_rawDesc = nil
	file_einride_iam_v1_annotations_proto_goTypes = nil
	file_einride_iam_v1_annotations_proto_depIdxs = nil
}
