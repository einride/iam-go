// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.2
// source: einride/iam/v1/caller.proto

package iamv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Caller identity.
type Caller struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The caller's resolved IAM members.
	Members []string `protobuf:"bytes,1,rep,name=members,proto3" json:"members,omitempty"`
	// Caller identity from gRPC metadata key/value pairs.
	Metadata map[string]*Caller_Metadata `protobuf:"bytes,2,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Caller) Reset() {
	*x = Caller{}
	if protoimpl.UnsafeEnabled {
		mi := &file_einride_iam_v1_caller_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Caller) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Caller) ProtoMessage() {}

func (x *Caller) ProtoReflect() protoreflect.Message {
	mi := &file_einride_iam_v1_caller_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Caller.ProtoReflect.Descriptor instead.
func (*Caller) Descriptor() ([]byte, []int) {
	return file_einride_iam_v1_caller_proto_rawDescGZIP(), []int{0}
}

func (x *Caller) GetMembers() []string {
	if x != nil {
		return x.Members
	}
	return nil
}

func (x *Caller) GetMetadata() map[string]*Caller_Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

// Caller identity for a gRPC metadata key/value pair.
type Caller_Metadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The IAM members from the metadata value.
	Members []string `protobuf:"bytes,1,rep,name=members,proto3" json:"members,omitempty"`
	// The identity token from the metadata value.
	IdentityToken *IdentityToken `protobuf:"bytes,2,opt,name=identity_token,json=identityToken,proto3" json:"identity_token,omitempty"`
}

func (x *Caller_Metadata) Reset() {
	*x = Caller_Metadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_einride_iam_v1_caller_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Caller_Metadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Caller_Metadata) ProtoMessage() {}

func (x *Caller_Metadata) ProtoReflect() protoreflect.Message {
	mi := &file_einride_iam_v1_caller_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Caller_Metadata.ProtoReflect.Descriptor instead.
func (*Caller_Metadata) Descriptor() ([]byte, []int) {
	return file_einride_iam_v1_caller_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Caller_Metadata) GetMembers() []string {
	if x != nil {
		return x.Members
	}
	return nil
}

func (x *Caller_Metadata) GetIdentityToken() *IdentityToken {
	if x != nil {
		return x.IdentityToken
	}
	return nil
}

var File_einride_iam_v1_caller_proto protoreflect.FileDescriptor

var file_einride_iam_v1_caller_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x76, 0x31,
	0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x65,
	0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x1a, 0x23, 0x65,
	0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xae, 0x02, 0x0a, 0x06, 0x43, 0x61, 0x6c, 0x6c, 0x65, 0x72, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x12, 0x40, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x65, 0x69, 0x6e,
	0x72, 0x69, 0x64, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x61, 0x6c, 0x6c,
	0x65, 0x72, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x5c, 0x0a, 0x0d, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x35, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x65,
	0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x61,
	0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x6a, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x12, 0x44,
	0x0a, 0x0e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64, 0x65,
	0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x0d, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x42, 0x70, 0x0a, 0x13, 0x74, 0x65, 0x63, 0x68, 0x2e, 0x65, 0x69, 0x6e,
	0x72, 0x69, 0x64, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x42, 0x0b, 0x43, 0x61, 0x6c,
	0x6c, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x32, 0x67, 0x6f, 0x2e, 0x65,
	0x69, 0x6e, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x74, 0x65, 0x63, 0x68, 0x2f, 0x69, 0x61, 0x6d, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x65, 0x69, 0x6e, 0x72, 0x69, 0x64,
	0x65, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x76, 0x31, 0x3b, 0x69, 0x61, 0x6d, 0x76, 0x31, 0xa2, 0x02,
	0x03, 0x49, 0x58, 0x58, 0xaa, 0x02, 0x06, 0x49, 0x61, 0x6d, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x06,
	0x49, 0x61, 0x6d, 0x5c, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_einride_iam_v1_caller_proto_rawDescOnce sync.Once
	file_einride_iam_v1_caller_proto_rawDescData = file_einride_iam_v1_caller_proto_rawDesc
)

func file_einride_iam_v1_caller_proto_rawDescGZIP() []byte {
	file_einride_iam_v1_caller_proto_rawDescOnce.Do(func() {
		file_einride_iam_v1_caller_proto_rawDescData = protoimpl.X.CompressGZIP(file_einride_iam_v1_caller_proto_rawDescData)
	})
	return file_einride_iam_v1_caller_proto_rawDescData
}

var file_einride_iam_v1_caller_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_einride_iam_v1_caller_proto_goTypes = []interface{}{
	(*Caller)(nil),          // 0: einride.iam.v1.Caller
	nil,                     // 1: einride.iam.v1.Caller.MetadataEntry
	(*Caller_Metadata)(nil), // 2: einride.iam.v1.Caller.Metadata
	(*IdentityToken)(nil),   // 3: einride.iam.v1.IdentityToken
}
var file_einride_iam_v1_caller_proto_depIdxs = []int32{
	1, // 0: einride.iam.v1.Caller.metadata:type_name -> einride.iam.v1.Caller.MetadataEntry
	2, // 1: einride.iam.v1.Caller.MetadataEntry.value:type_name -> einride.iam.v1.Caller.Metadata
	3, // 2: einride.iam.v1.Caller.Metadata.identity_token:type_name -> einride.iam.v1.IdentityToken
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_einride_iam_v1_caller_proto_init() }
func file_einride_iam_v1_caller_proto_init() {
	if File_einride_iam_v1_caller_proto != nil {
		return
	}
	file_einride_iam_v1_identity_token_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_einride_iam_v1_caller_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Caller); i {
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
		file_einride_iam_v1_caller_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Caller_Metadata); i {
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
			RawDescriptor: file_einride_iam_v1_caller_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_einride_iam_v1_caller_proto_goTypes,
		DependencyIndexes: file_einride_iam_v1_caller_proto_depIdxs,
		MessageInfos:      file_einride_iam_v1_caller_proto_msgTypes,
	}.Build()
	File_einride_iam_v1_caller_proto = out.File
	file_einride_iam_v1_caller_proto_rawDesc = nil
	file_einride_iam_v1_caller_proto_goTypes = nil
	file_einride_iam_v1_caller_proto_depIdxs = nil
}
