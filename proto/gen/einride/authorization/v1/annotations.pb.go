// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.2
// source: einride/authorization/v1/annotations.proto

package authorizationv1

import (
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

type PolicyType int32

const (
	PolicyType_POLICY_TYPE_UNSPECIFIED PolicyType = 0
	PolicyType_CUSTOM_MIDDLEWARE       PolicyType = 1
	PolicyType_REQUEST                 PolicyType = 2
	PolicyType_RESPONSE                PolicyType = 3
)

// Enum value maps for PolicyType.
var (
	PolicyType_name = map[int32]string{
		0: "POLICY_TYPE_UNSPECIFIED",
		1: "CUSTOM_MIDDLEWARE",
		2: "REQUEST",
		3: "RESPONSE",
	}
	PolicyType_value = map[string]int32{
		"POLICY_TYPE_UNSPECIFIED": 0,
		"CUSTOM_MIDDLEWARE":       1,
		"REQUEST":                 2,
		"RESPONSE":                3,
	}
)

func (x PolicyType) Enum() *PolicyType {
	p := new(PolicyType)
	*p = x
	return p
}

func (x PolicyType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PolicyType) Descriptor() protoreflect.EnumDescriptor {
	return file_einride_authorization_v1_annotations_proto_enumTypes[0].Descriptor()
}

func (PolicyType) Type() protoreflect.EnumType {
	return &file_einride_authorization_v1_annotations_proto_enumTypes[0]
}

func (x PolicyType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PolicyType.Descriptor instead.
func (PolicyType) EnumDescriptor() ([]byte, []int) {
	return file_einride_authorization_v1_annotations_proto_rawDescGZIP(), []int{0}
}

type Policy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type       PolicyType `protobuf:"varint,1,opt,name=type,proto3,enum=einride.authorization.v1.PolicyType" json:"type,omitempty"`
	Expression string     `protobuf:"bytes,2,opt,name=expression,proto3" json:"expression,omitempty"`
}

func (x *Policy) Reset() {
	*x = Policy{}
	if protoimpl.UnsafeEnabled {
		mi := &file_einride_authorization_v1_annotations_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Policy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Policy) ProtoMessage() {}

func (x *Policy) ProtoReflect() protoreflect.Message {
	mi := &file_einride_authorization_v1_annotations_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Policy.ProtoReflect.Descriptor instead.
func (*Policy) Descriptor() ([]byte, []int) {
	return file_einride_authorization_v1_annotations_proto_rawDescGZIP(), []int{0}
}

func (x *Policy) GetType() PolicyType {
	if x != nil {
		return x.Type
	}
	return PolicyType_POLICY_TYPE_UNSPECIFIED
}

func (x *Policy) GetExpression() string {
	if x != nil {
		return x.Expression
	}
	return ""
}

var file_einride_authorization_v1_annotations_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*Policy)(nil),
		Field:         11000,
		Name:          "einride.authorization.v1.policy",
		Tag:           "bytes,11000,opt,name=policy",
		Filename:      "einride/authorization/v1/annotations.proto",
	},
}

// Extension fields to descriptorpb.MethodOptions.
var (
	// optional einride.authorization.v1.Policy policy = 11000;
	E_Policy = &file_einride_authorization_v1_annotations_proto_extTypes[0]
)

var File_einride_authorization_v1_annotations_proto protoreflect.FileDescriptor

var file_einride_authorization_v1_annotations_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x65, 0x69,
	0x6e, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x62, 0x0a, 0x06, 0x50, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x12, 0x38, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x24, 0x2e, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a, 0x0a,
	0x65, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x65, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2a, 0x5b, 0x0a, 0x0a,
	0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a, 0x17, 0x50, 0x4f,
	0x4c, 0x49, 0x43, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43,
	0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x43, 0x55, 0x53, 0x54, 0x4f,
	0x4d, 0x5f, 0x4d, 0x49, 0x44, 0x44, 0x4c, 0x45, 0x57, 0x41, 0x52, 0x45, 0x10, 0x01, 0x12, 0x0b,
	0x0a, 0x07, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x52,
	0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x10, 0x03, 0x3a, 0x59, 0x0a, 0x06, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0xf8, 0x55, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x65, 0x69, 0x6e,
	0x72, 0x69, 0x64, 0x65, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x06, 0x70, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x42, 0xcc, 0x01, 0x0a, 0x1d, 0x74, 0x65, 0x63, 0x68, 0x2e, 0x65, 0x69,
	0x6e, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x42, 0x10, 0x41, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x66, 0x67, 0x6f, 0x2e, 0x65,
	0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x74, 0x65, 0x63, 0x68, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x67, 0x6f, 0x2d, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2d, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x2f, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64,
	0x65, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f,
	0x76, 0x31, 0x3b, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x58, 0x58, 0xaa, 0x02, 0x15, 0x41, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x56, 0x31, 0x42, 0x65, 0x74, 0x61, 0x31,
	0xca, 0x02, 0x10, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5c, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_einride_authorization_v1_annotations_proto_rawDescOnce sync.Once
	file_einride_authorization_v1_annotations_proto_rawDescData = file_einride_authorization_v1_annotations_proto_rawDesc
)

func file_einride_authorization_v1_annotations_proto_rawDescGZIP() []byte {
	file_einride_authorization_v1_annotations_proto_rawDescOnce.Do(func() {
		file_einride_authorization_v1_annotations_proto_rawDescData = protoimpl.X.CompressGZIP(file_einride_authorization_v1_annotations_proto_rawDescData)
	})
	return file_einride_authorization_v1_annotations_proto_rawDescData
}

var file_einride_authorization_v1_annotations_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_einride_authorization_v1_annotations_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_einride_authorization_v1_annotations_proto_goTypes = []interface{}{
	(PolicyType)(0),                    // 0: einride.authorization.v1.PolicyType
	(*Policy)(nil),                     // 1: einride.authorization.v1.Policy
	(*descriptorpb.MethodOptions)(nil), // 2: google.protobuf.MethodOptions
}
var file_einride_authorization_v1_annotations_proto_depIdxs = []int32{
	0, // 0: einride.authorization.v1.Policy.type:type_name -> einride.authorization.v1.PolicyType
	2, // 1: einride.authorization.v1.policy:extendee -> google.protobuf.MethodOptions
	1, // 2: einride.authorization.v1.policy:type_name -> einride.authorization.v1.Policy
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	2, // [2:3] is the sub-list for extension type_name
	1, // [1:2] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_einride_authorization_v1_annotations_proto_init() }
func file_einride_authorization_v1_annotations_proto_init() {
	if File_einride_authorization_v1_annotations_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_einride_authorization_v1_annotations_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Policy); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_einride_authorization_v1_annotations_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 1,
			NumServices:   0,
		},
		GoTypes:           file_einride_authorization_v1_annotations_proto_goTypes,
		DependencyIndexes: file_einride_authorization_v1_annotations_proto_depIdxs,
		EnumInfos:         file_einride_authorization_v1_annotations_proto_enumTypes,
		MessageInfos:      file_einride_authorization_v1_annotations_proto_msgTypes,
		ExtensionInfos:    file_einride_authorization_v1_annotations_proto_extTypes,
	}.Build()
	File_einride_authorization_v1_annotations_proto = out.File
	file_einride_authorization_v1_annotations_proto_rawDesc = nil
	file_einride_authorization_v1_annotations_proto_goTypes = nil
	file_einride_authorization_v1_annotations_proto_depIdxs = nil
}
