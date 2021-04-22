// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.2
// source: einride/iam/v1/annotations.proto

package iamv1

import (
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

// Authorization options for a gRPC method.
type Authorization struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Permission to use for authorization.
	//
	// Types that are assignable to Permissions:
	//	*Authorization_Permission
	//	*Authorization_ResourcePermissions_
	Permissions isAuthorization_Permissions `protobuf_oneof:"permissions"`
	// Strategy that decides if the request is authorized.
	//
	// Types that are assignable to Strategy:
	//	*Authorization_Before
	//	*Authorization_After
	//	*Authorization_Custom
	//	*Authorization_Open
	Strategy isAuthorization_Strategy `protobuf_oneof:"strategy"`
}

func (x *Authorization) Reset() {
	*x = Authorization{}
	if protoimpl.UnsafeEnabled {
		mi := &file_einride_iam_v1_annotations_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Authorization) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Authorization) ProtoMessage() {}

func (x *Authorization) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Authorization.ProtoReflect.Descriptor instead.
func (*Authorization) Descriptor() ([]byte, []int) {
	return file_einride_iam_v1_annotations_proto_rawDescGZIP(), []int{0}
}

func (m *Authorization) GetPermissions() isAuthorization_Permissions {
	if m != nil {
		return m.Permissions
	}
	return nil
}

func (x *Authorization) GetPermission() string {
	if x, ok := x.GetPermissions().(*Authorization_Permission); ok {
		return x.Permission
	}
	return ""
}

func (x *Authorization) GetResourcePermissions() *Authorization_ResourcePermissions {
	if x, ok := x.GetPermissions().(*Authorization_ResourcePermissions_); ok {
		return x.ResourcePermissions
	}
	return nil
}

func (m *Authorization) GetStrategy() isAuthorization_Strategy {
	if m != nil {
		return m.Strategy
	}
	return nil
}

func (x *Authorization) GetBefore() *expr.Expr {
	if x, ok := x.GetStrategy().(*Authorization_Before); ok {
		return x.Before
	}
	return nil
}

func (x *Authorization) GetAfter() *expr.Expr {
	if x, ok := x.GetStrategy().(*Authorization_After); ok {
		return x.After
	}
	return nil
}

func (x *Authorization) GetCustom() string {
	if x, ok := x.GetStrategy().(*Authorization_Custom); ok {
		return x.Custom
	}
	return ""
}

func (x *Authorization) GetOpen() bool {
	if x, ok := x.GetStrategy().(*Authorization_Open); ok {
		return x.Open
	}
	return false
}

type isAuthorization_Permissions interface {
	isAuthorization_Permissions()
}

type Authorization_Permission struct {
	// A single permission used by the method.
	Permission string `protobuf:"bytes,1,opt,name=permission,proto3,oneof"`
}

type Authorization_ResourcePermissions_ struct {
	// Resource permissions used by the method.
	ResourcePermissions *Authorization_ResourcePermissions `protobuf:"bytes,2,opt,name=resource_permissions,json=resourcePermissions,proto3,oneof"`
}

func (*Authorization_Permission) isAuthorization_Permissions() {}

func (*Authorization_ResourcePermissions_) isAuthorization_Permissions() {}

type isAuthorization_Strategy interface {
	isAuthorization_Strategy()
}

type Authorization_Before struct {
	// Expression that decides before the request if the caller is authorized.
	Before *expr.Expr `protobuf:"bytes,3,opt,name=before,proto3,oneof"`
}

type Authorization_After struct {
	// Expression that decides after the request if the caller is authorized.
	After *expr.Expr `protobuf:"bytes,4,opt,name=after,proto3,oneof"`
}

type Authorization_Custom struct {
	// A comment explaining a custom way of determining if the caller is authorized.
	Custom string `protobuf:"bytes,5,opt,name=custom,proto3,oneof"`
}

type Authorization_Open struct {
	// A flag indicating if the method is open.
	Open bool `protobuf:"varint,6,opt,name=open,proto3,oneof"`
}

func (*Authorization_Before) isAuthorization_Strategy() {}

func (*Authorization_After) isAuthorization_Strategy() {}

func (*Authorization_Custom) isAuthorization_Strategy() {}

func (*Authorization_Open) isAuthorization_Strategy() {}

// A list of roles.
type Roles struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Role []*v1.Role `protobuf:"bytes,1,rep,name=role,proto3" json:"role,omitempty"`
}

func (x *Roles) Reset() {
	*x = Roles{}
	if protoimpl.UnsafeEnabled {
		mi := &file_einride_iam_v1_annotations_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Roles) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Roles) ProtoMessage() {}

func (x *Roles) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Roles.ProtoReflect.Descriptor instead.
func (*Roles) Descriptor() ([]byte, []int) {
	return file_einride_iam_v1_annotations_proto_rawDescGZIP(), []int{1}
}

func (x *Roles) GetRole() []*v1.Role {
	if x != nil {
		return x.Role
	}
	return nil
}

// Resource permissions.
type Authorization_ResourcePermissions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The resource permissions.
	Resource []*Authorization_ResourcePermission `protobuf:"bytes,1,rep,name=resource,proto3" json:"resource,omitempty"`
}

func (x *Authorization_ResourcePermissions) Reset() {
	*x = Authorization_ResourcePermissions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_einride_iam_v1_annotations_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Authorization_ResourcePermissions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Authorization_ResourcePermissions) ProtoMessage() {}

func (x *Authorization_ResourcePermissions) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Authorization_ResourcePermissions.ProtoReflect.Descriptor instead.
func (*Authorization_ResourcePermissions) Descriptor() ([]byte, []int) {
	return file_einride_iam_v1_annotations_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Authorization_ResourcePermissions) GetResource() []*Authorization_ResourcePermission {
	if x != nil {
		return x.Resource
	}
	return nil
}

// A resource type and a permission.
type Authorization_ResourcePermission struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The resource type.
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// The permission.
	Permission string `protobuf:"bytes,2,opt,name=permission,proto3" json:"permission,omitempty"`
}

func (x *Authorization_ResourcePermission) Reset() {
	*x = Authorization_ResourcePermission{}
	if protoimpl.UnsafeEnabled {
		mi := &file_einride_iam_v1_annotations_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Authorization_ResourcePermission) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Authorization_ResourcePermission) ProtoMessage() {}

func (x *Authorization_ResourcePermission) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Authorization_ResourcePermission.ProtoReflect.Descriptor instead.
func (*Authorization_ResourcePermission) Descriptor() ([]byte, []int) {
	return file_einride_iam_v1_annotations_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Authorization_ResourcePermission) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Authorization_ResourcePermission) GetPermission() string {
	if x != nil {
		return x.Permission
	}
	return ""
}

var file_einride_iam_v1_annotations_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*Authorization)(nil),
		Field:         11000,
		Name:          "einride.iam.v1.authorization",
		Tag:           "bytes,11000,opt,name=authorization",
		Filename:      "einride/iam/v1/annotations.proto",
	},
	{
		ExtendedType:  (*descriptorpb.ServiceOptions)(nil),
		ExtensionType: (*Roles)(nil),
		Field:         11001,
		Name:          "einride.iam.v1.predefined_roles",
		Tag:           "bytes,11001,opt,name=predefined_roles",
		Filename:      "einride/iam/v1/annotations.proto",
	},
}

// Extension fields to descriptorpb.MethodOptions.
var (
	// Authorization for the method.
	//
	// optional einride.iam.v1.Authorization authorization = 11000;
	E_Authorization = &file_einride_iam_v1_annotations_proto_extTypes[0]
)

// Extension fields to descriptorpb.ServiceOptions.
var (
	// Predefined roles for the service.
	//
	// optional einride.iam.v1.Roles predefined_roles = 11001;
	E_PredefinedRoles = &file_einride_iam_v1_annotations_proto_extTypes[1]
)

var File_einride_iam_v1_annotations_proto protoreflect.FileDescriptor

var file_einride_iam_v1_annotations_proto_rawDesc = []byte{
	0x0a, 0x20, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x76, 0x31,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0e, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e,
	0x76, 0x31, 0x1a, 0x1d, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x74, 0x79, 0x70, 0x65,
	0x2f, 0x65, 0x78, 0x70, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xeb, 0x03, 0x0a, 0x0d,
	0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a,
	0x0a, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x0a, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x66, 0x0a, 0x14, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x70, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x31, 0x2e,
	0x65, 0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x48, 0x00, 0x52, 0x13, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x2b, 0x0a, 0x06, 0x62, 0x65, 0x66, 0x6f, 0x72,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x45, 0x78, 0x70, 0x72, 0x48, 0x01, 0x52, 0x06, 0x62, 0x65,
	0x66, 0x6f, 0x72, 0x65, 0x12, 0x29, 0x0a, 0x05, 0x61, 0x66, 0x74, 0x65, 0x72, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x74, 0x79, 0x70,
	0x65, 0x2e, 0x45, 0x78, 0x70, 0x72, 0x48, 0x01, 0x52, 0x05, 0x61, 0x66, 0x74, 0x65, 0x72, 0x12,
	0x18, 0x0a, 0x06, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x01, 0x52, 0x06, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x12, 0x14, 0x0a, 0x04, 0x6f, 0x70, 0x65,
	0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x48, 0x01, 0x52, 0x04, 0x6f, 0x70, 0x65, 0x6e, 0x1a,
	0x63, 0x0a, 0x13, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x4c, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x65, 0x69, 0x6e, 0x72, 0x69,
	0x64, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x1a, 0x48, 0x0a, 0x12, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1e,
	0x0a, 0x0a, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x0d,
	0x0a, 0x0b, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x0a, 0x0a,
	0x08, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x22, 0x36, 0x0a, 0x05, 0x52, 0x6f, 0x6c,
	0x65, 0x73, 0x12, 0x2d, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x6f, 0x6c,
	0x65, 0x3a, 0x64, 0x0a, 0x0d, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0xf8, 0x55, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x65, 0x69, 0x6e, 0x72,
	0x69, 0x64, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x62, 0x0a, 0x10, 0x70, 0x72, 0x65, 0x64, 0x65,
	0x66, 0x69, 0x6e, 0x65, 0x64, 0x5f, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x12, 0x1f, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xf9, 0x55, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x69, 0x61,
	0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x52, 0x0f, 0x70, 0x72, 0x65, 0x64,
	0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x42, 0x75, 0x0a, 0x13, 0x74,
	0x65, 0x63, 0x68, 0x2e, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e,
	0x76, 0x31, 0x42, 0x10, 0x41, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x32, 0x67, 0x6f, 0x2e, 0x65, 0x69, 0x6e, 0x72, 0x69,
	0x64, 0x65, 0x2e, 0x74, 0x65, 0x63, 0x68, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2f, 0x69, 0x61,
	0x6d, 0x2f, 0x76, 0x31, 0x3b, 0x69, 0x61, 0x6d, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x49, 0x58, 0x58,
	0xaa, 0x02, 0x06, 0x49, 0x61, 0x6d, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x06, 0x49, 0x61, 0x6d, 0x5c,
	0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_einride_iam_v1_annotations_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_einride_iam_v1_annotations_proto_goTypes = []interface{}{
	(*Authorization)(nil),                     // 0: einride.iam.v1.Authorization
	(*Roles)(nil),                             // 1: einride.iam.v1.Roles
	(*Authorization_ResourcePermissions)(nil), // 2: einride.iam.v1.Authorization.ResourcePermissions
	(*Authorization_ResourcePermission)(nil),  // 3: einride.iam.v1.Authorization.ResourcePermission
	(*expr.Expr)(nil),                         // 4: google.type.Expr
	(*v1.Role)(nil),                           // 5: google.iam.admin.v1.Role
	(*descriptorpb.MethodOptions)(nil),        // 6: google.protobuf.MethodOptions
	(*descriptorpb.ServiceOptions)(nil),       // 7: google.protobuf.ServiceOptions
}
var file_einride_iam_v1_annotations_proto_depIdxs = []int32{
	2, // 0: einride.iam.v1.Authorization.resource_permissions:type_name -> einride.iam.v1.Authorization.ResourcePermissions
	4, // 1: einride.iam.v1.Authorization.before:type_name -> google.type.Expr
	4, // 2: einride.iam.v1.Authorization.after:type_name -> google.type.Expr
	5, // 3: einride.iam.v1.Roles.role:type_name -> google.iam.admin.v1.Role
	3, // 4: einride.iam.v1.Authorization.ResourcePermissions.resource:type_name -> einride.iam.v1.Authorization.ResourcePermission
	6, // 5: einride.iam.v1.authorization:extendee -> google.protobuf.MethodOptions
	7, // 6: einride.iam.v1.predefined_roles:extendee -> google.protobuf.ServiceOptions
	0, // 7: einride.iam.v1.authorization:type_name -> einride.iam.v1.Authorization
	1, // 8: einride.iam.v1.predefined_roles:type_name -> einride.iam.v1.Roles
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	7, // [7:9] is the sub-list for extension type_name
	5, // [5:7] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_einride_iam_v1_annotations_proto_init() }
func file_einride_iam_v1_annotations_proto_init() {
	if File_einride_iam_v1_annotations_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_einride_iam_v1_annotations_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Authorization); i {
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
			switch v := v.(*Roles); i {
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
			switch v := v.(*Authorization_ResourcePermissions); i {
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
			switch v := v.(*Authorization_ResourcePermission); i {
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
	file_einride_iam_v1_annotations_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Authorization_Permission)(nil),
		(*Authorization_ResourcePermissions_)(nil),
		(*Authorization_Before)(nil),
		(*Authorization_After)(nil),
		(*Authorization_Custom)(nil),
		(*Authorization_Open)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_einride_iam_v1_annotations_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 2,
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
