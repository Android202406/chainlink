// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: tokendata.proto

package ccippb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// TokenDataRequest is a gRPC adapter for the input arguments of
// [github.com/smartcontractkit/chainlink-common/chainlink-common/pkg/types/ccip/TokenDataReader.ReadTokenData]
type TokenDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg        *EVM2EVMOnRampCCIPSendRequestedWithMeta `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	TokenIndex uint64                                  `protobuf:"varint,2,opt,name=token_index,json=tokenIndex,proto3" json:"token_index,omitempty"`
}

func (x *TokenDataRequest) Reset() {
	*x = TokenDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tokendata_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenDataRequest) ProtoMessage() {}

func (x *TokenDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tokendata_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenDataRequest.ProtoReflect.Descriptor instead.
func (*TokenDataRequest) Descriptor() ([]byte, []int) {
	return file_tokendata_proto_rawDescGZIP(), []int{0}
}

func (x *TokenDataRequest) GetMsg() *EVM2EVMOnRampCCIPSendRequestedWithMeta {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *TokenDataRequest) GetTokenIndex() uint64 {
	if x != nil {
		return x.TokenIndex
	}
	return 0
}

// TokenDataResponse is a gRPC adapter for the return value of
// [github.com/smartcontractkit/chainlink-common/chainlink-common/pkg/types/ccip/TokenDataReader.ReadTokenData]
type TokenDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TokenData []byte `protobuf:"bytes,1,opt,name=token_data,json=tokenData,proto3" json:"token_data,omitempty"`
}

func (x *TokenDataResponse) Reset() {
	*x = TokenDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tokendata_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenDataResponse) ProtoMessage() {}

func (x *TokenDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tokendata_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenDataResponse.ProtoReflect.Descriptor instead.
func (*TokenDataResponse) Descriptor() ([]byte, []int) {
	return file_tokendata_proto_rawDescGZIP(), []int{1}
}

func (x *TokenDataResponse) GetTokenData() []byte {
	if x != nil {
		return x.TokenData
	}
	return nil
}

var File_tokendata_proto protoreflect.FileDescriptor

var file_tokendata_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x15, 0x6c, 0x6f, 0x6f, 0x70, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2e, 0x70, 0x62, 0x2e, 0x63, 0x63, 0x69, 0x70, 0x1a, 0x0c, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x84, 0x01, 0x0a, 0x10, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x4f, 0x0a, 0x03, 0x6d,
	0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x3d, 0x2e, 0x6c, 0x6f, 0x6f, 0x70, 0x2e,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x70, 0x62, 0x2e, 0x63, 0x63, 0x69, 0x70,
	0x2e, 0x45, 0x56, 0x4d, 0x32, 0x45, 0x56, 0x4d, 0x4f, 0x6e, 0x52, 0x61, 0x6d, 0x70, 0x43, 0x43,
	0x49, 0x50, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64, 0x57,
	0x69, 0x74, 0x68, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x1f, 0x0a, 0x0b,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0a, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x32, 0x0a,
	0x11, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x44, 0x61, 0x74,
	0x61, 0x32, 0x77, 0x0a, 0x0f, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x12, 0x64, 0x0a, 0x0d, 0x52, 0x65, 0x61, 0x64, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x44, 0x61, 0x74, 0x61, 0x12, 0x27, 0x2e, 0x6c, 0x6f, 0x6f, 0x70, 0x2e, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x70, 0x62, 0x2e, 0x63, 0x63, 0x69, 0x70, 0x2e, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28,
	0x2e, 0x6c, 0x6f, 0x6f, 0x70, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x70,
	0x62, 0x2e, 0x63, 0x63, 0x69, 0x70, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x44, 0x61, 0x74, 0x61,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x4f, 0x5a, 0x4d, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6b, 0x69, 0x74, 0x2f, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x6c,
	0x69, 0x6e, 0x6b, 0x2d, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6c,
	0x6f, 0x6f, 0x70, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x62, 0x2f,
	0x63, 0x63, 0x69, 0x70, 0x3b, 0x63, 0x63, 0x69, 0x70, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_tokendata_proto_rawDescOnce sync.Once
	file_tokendata_proto_rawDescData = file_tokendata_proto_rawDesc
)

func file_tokendata_proto_rawDescGZIP() []byte {
	file_tokendata_proto_rawDescOnce.Do(func() {
		file_tokendata_proto_rawDescData = protoimpl.X.CompressGZIP(file_tokendata_proto_rawDescData)
	})
	return file_tokendata_proto_rawDescData
}

var file_tokendata_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_tokendata_proto_goTypes = []interface{}{
	(*TokenDataRequest)(nil),                       // 0: loop.internal.pb.ccip.TokenDataRequest
	(*TokenDataResponse)(nil),                      // 1: loop.internal.pb.ccip.TokenDataResponse
	(*EVM2EVMOnRampCCIPSendRequestedWithMeta)(nil), // 2: loop.internal.pb.ccip.EVM2EVMOnRampCCIPSendRequestedWithMeta
}
var file_tokendata_proto_depIdxs = []int32{
	2, // 0: loop.internal.pb.ccip.TokenDataRequest.msg:type_name -> loop.internal.pb.ccip.EVM2EVMOnRampCCIPSendRequestedWithMeta
	0, // 1: loop.internal.pb.ccip.TokenDataReader.ReadTokenData:input_type -> loop.internal.pb.ccip.TokenDataRequest
	1, // 2: loop.internal.pb.ccip.TokenDataReader.ReadTokenData:output_type -> loop.internal.pb.ccip.TokenDataResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_tokendata_proto_init() }
func file_tokendata_proto_init() {
	if File_tokendata_proto != nil {
		return
	}
	file_models_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_tokendata_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenDataRequest); i {
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
		file_tokendata_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenDataResponse); i {
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
			RawDescriptor: file_tokendata_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tokendata_proto_goTypes,
		DependencyIndexes: file_tokendata_proto_depIdxs,
		MessageInfos:      file_tokendata_proto_msgTypes,
	}.Build()
	File_tokendata_proto = out.File
	file_tokendata_proto_rawDesc = nil
	file_tokendata_proto_goTypes = nil
	file_tokendata_proto_depIdxs = nil
}
