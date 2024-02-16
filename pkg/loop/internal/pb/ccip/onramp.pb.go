// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: onramp.proto

package ccippb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// GetSendRequestBetweenSeqNumsRequest is a gRPC adapter for the input arguments of
// [github.com/smartcontractkit/chainlink-common/pkg/types/OnRampReader.GetSendRequestBetweenSeqNums]
type GetSendRequestBetweenSeqNumsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SeqNumMin uint64 `protobuf:"varint,1,opt,name=seq_num_min,json=seqNumMin,proto3" json:"seq_num_min,omitempty"`
	SeqNumMax uint64 `protobuf:"varint,2,opt,name=seq_num_max,json=seqNumMax,proto3" json:"seq_num_max,omitempty"`
	Finalized bool   `protobuf:"varint,3,opt,name=finalized,proto3" json:"finalized,omitempty"`
}

func (x *GetSendRequestBetweenSeqNumsRequest) Reset() {
	*x = GetSendRequestBetweenSeqNumsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_onramp_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSendRequestBetweenSeqNumsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSendRequestBetweenSeqNumsRequest) ProtoMessage() {}

func (x *GetSendRequestBetweenSeqNumsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_onramp_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSendRequestBetweenSeqNumsRequest.ProtoReflect.Descriptor instead.
func (*GetSendRequestBetweenSeqNumsRequest) Descriptor() ([]byte, []int) {
	return file_onramp_proto_rawDescGZIP(), []int{0}
}

func (x *GetSendRequestBetweenSeqNumsRequest) GetSeqNumMin() uint64 {
	if x != nil {
		return x.SeqNumMin
	}
	return 0
}

func (x *GetSendRequestBetweenSeqNumsRequest) GetSeqNumMax() uint64 {
	if x != nil {
		return x.SeqNumMax
	}
	return 0
}

func (x *GetSendRequestBetweenSeqNumsRequest) GetFinalized() bool {
	if x != nil {
		return x.Finalized
	}
	return false
}

// GetSendRequestBetweenSeqNumsResponse is a gRPC adapter for the output arguments of
// [github.com/smartcontractkit/chainlink-common/pkg/types/OnRampReader.GetSendRequestBetweenSeqNums]
type GetSendRequestBetweenSeqNumsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Note, the content here maybe better modeled as a oneof when CCIP supports
	// multiple types of messages/chains
	SendRequests []*EVM2EVMMessageWithTxMeta `protobuf:"bytes,1,rep,name=send_requests,json=sendRequests,proto3" json:"send_requests,omitempty"`
}

func (x *GetSendRequestBetweenSeqNumsResponse) Reset() {
	*x = GetSendRequestBetweenSeqNumsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_onramp_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSendRequestBetweenSeqNumsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSendRequestBetweenSeqNumsResponse) ProtoMessage() {}

func (x *GetSendRequestBetweenSeqNumsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_onramp_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSendRequestBetweenSeqNumsResponse.ProtoReflect.Descriptor instead.
func (*GetSendRequestBetweenSeqNumsResponse) Descriptor() ([]byte, []int) {
	return file_onramp_proto_rawDescGZIP(), []int{1}
}

func (x *GetSendRequestBetweenSeqNumsResponse) GetSendRequests() []*EVM2EVMMessageWithTxMeta {
	if x != nil {
		return x.SendRequests
	}
	return nil
}

// RouterAddressResponse is a gRPC adapter for the output arguments of
// [github.com/smartcontractkit/chainlink-common/pkg/types/OnRampReader.RouterAddress]
type RouterAddressResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RouterAddress string `protobuf:"bytes,1,opt,name=router_address,json=routerAddress,proto3" json:"router_address,omitempty"`
}

func (x *RouterAddressResponse) Reset() {
	*x = RouterAddressResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_onramp_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouterAddressResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouterAddressResponse) ProtoMessage() {}

func (x *RouterAddressResponse) ProtoReflect() protoreflect.Message {
	mi := &file_onramp_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouterAddressResponse.ProtoReflect.Descriptor instead.
func (*RouterAddressResponse) Descriptor() ([]byte, []int) {
	return file_onramp_proto_rawDescGZIP(), []int{2}
}

func (x *RouterAddressResponse) GetRouterAddress() string {
	if x != nil {
		return x.RouterAddress
	}
	return ""
}

// OnrampAddressResponse is a gRPC adapter for the output arguments of
// [github.com/smartcontractkit/chainlink-common/pkg/types/OnRampReader.Address]
type OnrampAddressResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *OnrampAddressResponse) Reset() {
	*x = OnrampAddressResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_onramp_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnrampAddressResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnrampAddressResponse) ProtoMessage() {}

func (x *OnrampAddressResponse) ProtoReflect() protoreflect.Message {
	mi := &file_onramp_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnrampAddressResponse.ProtoReflect.Descriptor instead.
func (*OnrampAddressResponse) Descriptor() ([]byte, []int) {
	return file_onramp_proto_rawDescGZIP(), []int{3}
}

func (x *OnrampAddressResponse) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

// GetDynamicConfigResponse is a gRPC adapter for the output arguments of
// [github.com/smartcontractkit/chainlink-common/pkg/types/OnRampReader.GetDynamicConfig]
type GetDynamicConfigResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DynamicConfig *OnRampDynamicConfig `protobuf:"bytes,1,opt,name=dynamic_config,json=dynamicConfig,proto3" json:"dynamic_config,omitempty"`
}

func (x *GetDynamicConfigResponse) Reset() {
	*x = GetDynamicConfigResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_onramp_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDynamicConfigResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDynamicConfigResponse) ProtoMessage() {}

func (x *GetDynamicConfigResponse) ProtoReflect() protoreflect.Message {
	mi := &file_onramp_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDynamicConfigResponse.ProtoReflect.Descriptor instead.
func (*GetDynamicConfigResponse) Descriptor() ([]byte, []int) {
	return file_onramp_proto_rawDescGZIP(), []int{4}
}

func (x *GetDynamicConfigResponse) GetDynamicConfig() *OnRampDynamicConfig {
	if x != nil {
		return x.DynamicConfig
	}
	return nil
}

// OnRampDynamicConfig is a gRPC adapter for the struct
// [github.com/smartcontractkit/chainlink-common/pkg/types/OnRampDynamicConfig]
type OnRampDynamicConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Router                            string `protobuf:"bytes,1,opt,name=router,proto3" json:"router,omitempty"` // Address
	MaxNumberOfTokenPerMsg            uint32 `protobuf:"varint,2,opt,name=max_number_of_token_per_msg,json=maxNumberOfTokenPerMsg,proto3" json:"max_number_of_token_per_msg,omitempty"`
	DestGasOverhead                   uint32 `protobuf:"varint,3,opt,name=dest_gas_overhead,json=destGasOverhead,proto3" json:"dest_gas_overhead,omitempty"`
	DestGasPerByte                    uint32 `protobuf:"varint,4,opt,name=dest_gas_per_byte,json=destGasPerByte,proto3" json:"dest_gas_per_byte,omitempty"`
	DestDataAvailabilityOverheadGas   uint32 `protobuf:"varint,5,opt,name=dest_data_availability_overhead_gas,json=destDataAvailabilityOverheadGas,proto3" json:"dest_data_availability_overhead_gas,omitempty"`
	DestDataAvailabilityGasPerByte    uint32 `protobuf:"varint,6,opt,name=dest_data_availability_gas_per_byte,json=destDataAvailabilityGasPerByte,proto3" json:"dest_data_availability_gas_per_byte,omitempty"`
	DestDataAvailabilityMultiplierBps uint32 `protobuf:"varint,7,opt,name=dest_data_availability_multiplier_bps,json=destDataAvailabilityMultiplierBps,proto3" json:"dest_data_availability_multiplier_bps,omitempty"`
	PriceRegistry                     string `protobuf:"bytes,8,opt,name=price_registry,json=priceRegistry,proto3" json:"price_registry,omitempty"` // Address
	MaxDataBytes                      uint32 `protobuf:"varint,9,opt,name=max_data_bytes,json=maxDataBytes,proto3" json:"max_data_bytes,omitempty"`
	MaxPerMsgGasLimite                uint32 `protobuf:"varint,10,opt,name=max_per_msg_gas_limite,json=maxPerMsgGasLimite,proto3" json:"max_per_msg_gas_limite,omitempty"`
}

func (x *OnRampDynamicConfig) Reset() {
	*x = OnRampDynamicConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_onramp_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnRampDynamicConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnRampDynamicConfig) ProtoMessage() {}

func (x *OnRampDynamicConfig) ProtoReflect() protoreflect.Message {
	mi := &file_onramp_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnRampDynamicConfig.ProtoReflect.Descriptor instead.
func (*OnRampDynamicConfig) Descriptor() ([]byte, []int) {
	return file_onramp_proto_rawDescGZIP(), []int{5}
}

func (x *OnRampDynamicConfig) GetRouter() string {
	if x != nil {
		return x.Router
	}
	return ""
}

func (x *OnRampDynamicConfig) GetMaxNumberOfTokenPerMsg() uint32 {
	if x != nil {
		return x.MaxNumberOfTokenPerMsg
	}
	return 0
}

func (x *OnRampDynamicConfig) GetDestGasOverhead() uint32 {
	if x != nil {
		return x.DestGasOverhead
	}
	return 0
}

func (x *OnRampDynamicConfig) GetDestGasPerByte() uint32 {
	if x != nil {
		return x.DestGasPerByte
	}
	return 0
}

func (x *OnRampDynamicConfig) GetDestDataAvailabilityOverheadGas() uint32 {
	if x != nil {
		return x.DestDataAvailabilityOverheadGas
	}
	return 0
}

func (x *OnRampDynamicConfig) GetDestDataAvailabilityGasPerByte() uint32 {
	if x != nil {
		return x.DestDataAvailabilityGasPerByte
	}
	return 0
}

func (x *OnRampDynamicConfig) GetDestDataAvailabilityMultiplierBps() uint32 {
	if x != nil {
		return x.DestDataAvailabilityMultiplierBps
	}
	return 0
}

func (x *OnRampDynamicConfig) GetPriceRegistry() string {
	if x != nil {
		return x.PriceRegistry
	}
	return ""
}

func (x *OnRampDynamicConfig) GetMaxDataBytes() uint32 {
	if x != nil {
		return x.MaxDataBytes
	}
	return 0
}

func (x *OnRampDynamicConfig) GetMaxPerMsgGasLimite() uint32 {
	if x != nil {
		return x.MaxPerMsgGasLimite
	}
	return 0
}

// EVM2EVMMessageWithTxMeta is a gRPC adapter for the struct
// [github.com/smartcontractkit/chainlink-common/pkg/types/EVM2EVMMessageWithTxMeta]
type EVM2EVMMessageWithTxMeta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message *EVM2EVMMessage `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	TxMeta  *TxMeta         `protobuf:"bytes,2,opt,name=tx_meta,json=txMeta,proto3" json:"tx_meta,omitempty"`
}

func (x *EVM2EVMMessageWithTxMeta) Reset() {
	*x = EVM2EVMMessageWithTxMeta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_onramp_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EVM2EVMMessageWithTxMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EVM2EVMMessageWithTxMeta) ProtoMessage() {}

func (x *EVM2EVMMessageWithTxMeta) ProtoReflect() protoreflect.Message {
	mi := &file_onramp_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EVM2EVMMessageWithTxMeta.ProtoReflect.Descriptor instead.
func (*EVM2EVMMessageWithTxMeta) Descriptor() ([]byte, []int) {
	return file_onramp_proto_rawDescGZIP(), []int{6}
}

func (x *EVM2EVMMessageWithTxMeta) GetMessage() *EVM2EVMMessage {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *EVM2EVMMessageWithTxMeta) GetTxMeta() *TxMeta {
	if x != nil {
		return x.TxMeta
	}
	return nil
}

var File_onramp_proto protoreflect.FileDescriptor

var file_onramp_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6f, 0x6e, 0x72, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15,
	0x6c, 0x6f, 0x6f, 0x70, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x70, 0x62,
	0x2e, 0x63, 0x63, 0x69, 0x70, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x0c, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x83, 0x01, 0x0a, 0x23, 0x47, 0x65, 0x74, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x42, 0x65, 0x74, 0x77, 0x65, 0x65, 0x6e, 0x53, 0x65, 0x71, 0x4e, 0x75, 0x6d,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0b, 0x73, 0x65, 0x71, 0x5f,
	0x6e, 0x75, 0x6d, 0x5f, 0x6d, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x73,
	0x65, 0x71, 0x4e, 0x75, 0x6d, 0x4d, 0x69, 0x6e, 0x12, 0x1e, 0x0a, 0x0b, 0x73, 0x65, 0x71, 0x5f,
	0x6e, 0x75, 0x6d, 0x5f, 0x6d, 0x61, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x73,
	0x65, 0x71, 0x4e, 0x75, 0x6d, 0x4d, 0x61, 0x78, 0x12, 0x1c, 0x0a, 0x09, 0x66, 0x69, 0x6e, 0x61,
	0x6c, 0x69, 0x7a, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x66, 0x69, 0x6e,
	0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x22, 0x7c, 0x0a, 0x24, 0x47, 0x65, 0x74, 0x53, 0x65, 0x6e,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x65, 0x74, 0x77, 0x65, 0x65, 0x6e, 0x53,
	0x65, 0x71, 0x4e, 0x75, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x54,
	0x0a, 0x0d, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x6c, 0x6f, 0x6f, 0x70, 0x2e, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x70, 0x62, 0x2e, 0x63, 0x63, 0x69, 0x70, 0x2e, 0x45, 0x56,
	0x4d, 0x32, 0x45, 0x56, 0x4d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x69, 0x74, 0x68,
	0x54, 0x78, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x0c, 0x73, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x73, 0x22, 0x3e, 0x0a, 0x15, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a,
	0x0e, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x22, 0x31, 0x0a, 0x15, 0x4f, 0x6e, 0x72, 0x61, 0x6d, 0x70, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x6d, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x44, 0x79,
	0x6e, 0x61, 0x6d, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x51, 0x0a, 0x0e, 0x64, 0x79, 0x6e, 0x61, 0x6d, 0x69, 0x63, 0x5f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x6c, 0x6f,
	0x6f, 0x70, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x70, 0x62, 0x2e, 0x63,
	0x63, 0x69, 0x70, 0x2e, 0x4f, 0x6e, 0x52, 0x61, 0x6d, 0x70, 0x44, 0x79, 0x6e, 0x61, 0x6d, 0x69,
	0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x0d, 0x64, 0x79, 0x6e, 0x61, 0x6d, 0x69, 0x63,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0xaf, 0x04, 0x0a, 0x13, 0x4f, 0x6e, 0x52, 0x61, 0x6d,
	0x70, 0x44, 0x79, 0x6e, 0x61, 0x6d, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x16,
	0x0a, 0x06, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x12, 0x3b, 0x0a, 0x1b, 0x6d, 0x61, 0x78, 0x5f, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x5f, 0x6f, 0x66, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x70, 0x65,
	0x72, 0x5f, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x16, 0x6d, 0x61, 0x78,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x4f, 0x66, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x50, 0x65, 0x72,
	0x4d, 0x73, 0x67, 0x12, 0x2a, 0x0a, 0x11, 0x64, 0x65, 0x73, 0x74, 0x5f, 0x67, 0x61, 0x73, 0x5f,
	0x6f, 0x76, 0x65, 0x72, 0x68, 0x65, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0f,
	0x64, 0x65, 0x73, 0x74, 0x47, 0x61, 0x73, 0x4f, 0x76, 0x65, 0x72, 0x68, 0x65, 0x61, 0x64, 0x12,
	0x29, 0x0a, 0x11, 0x64, 0x65, 0x73, 0x74, 0x5f, 0x67, 0x61, 0x73, 0x5f, 0x70, 0x65, 0x72, 0x5f,
	0x62, 0x79, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0e, 0x64, 0x65, 0x73, 0x74,
	0x47, 0x61, 0x73, 0x50, 0x65, 0x72, 0x42, 0x79, 0x74, 0x65, 0x12, 0x4c, 0x0a, 0x23, 0x64, 0x65,
	0x73, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69,
	0x6c, 0x69, 0x74, 0x79, 0x5f, 0x6f, 0x76, 0x65, 0x72, 0x68, 0x65, 0x61, 0x64, 0x5f, 0x67, 0x61,
	0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x1f, 0x64, 0x65, 0x73, 0x74, 0x44, 0x61, 0x74,
	0x61, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x4f, 0x76, 0x65,
	0x72, 0x68, 0x65, 0x61, 0x64, 0x47, 0x61, 0x73, 0x12, 0x4b, 0x0a, 0x23, 0x64, 0x65, 0x73, 0x74,
	0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69,
	0x74, 0x79, 0x5f, 0x67, 0x61, 0x73, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x1e, 0x64, 0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x41,
	0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x47, 0x61, 0x73, 0x50, 0x65,
	0x72, 0x42, 0x79, 0x74, 0x65, 0x12, 0x50, 0x0a, 0x25, 0x64, 0x65, 0x73, 0x74, 0x5f, 0x64, 0x61,
	0x74, 0x61, 0x5f, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f,
	0x6d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x69, 0x65, 0x72, 0x5f, 0x62, 0x70, 0x73, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x21, 0x64, 0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x41, 0x76,
	0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70,
	0x6c, 0x69, 0x65, 0x72, 0x42, 0x70, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x70, 0x72, 0x69, 0x63, 0x65,
	0x5f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x70, 0x72, 0x69, 0x63, 0x65, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x12, 0x24,
	0x0a, 0x0e, 0x6d, 0x61, 0x78, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x6d, 0x61, 0x78, 0x44, 0x61, 0x74, 0x61, 0x42,
	0x79, 0x74, 0x65, 0x73, 0x12, 0x32, 0x0a, 0x16, 0x6d, 0x61, 0x78, 0x5f, 0x70, 0x65, 0x72, 0x5f,
	0x6d, 0x73, 0x67, 0x5f, 0x67, 0x61, 0x73, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x12, 0x6d, 0x61, 0x78, 0x50, 0x65, 0x72, 0x4d, 0x73, 0x67, 0x47,
	0x61, 0x73, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x22, 0x93, 0x01, 0x0a, 0x18, 0x45, 0x56, 0x4d,
	0x32, 0x45, 0x56, 0x4d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x69, 0x74, 0x68, 0x54,
	0x78, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x3f, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x6c, 0x6f, 0x6f, 0x70, 0x2e, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x70, 0x62, 0x2e, 0x63, 0x63, 0x69, 0x70, 0x2e, 0x45,
	0x56, 0x4d, 0x32, 0x45, 0x56, 0x4d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x36, 0x0a, 0x07, 0x74, 0x78, 0x5f, 0x6d, 0x65, 0x74,
	0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6c, 0x6f, 0x6f, 0x70, 0x2e, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x70, 0x62, 0x2e, 0x63, 0x63, 0x69, 0x70, 0x2e,
	0x54, 0x78, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x06, 0x74, 0x78, 0x4d, 0x65, 0x74, 0x61, 0x32, 0xb5,
	0x03, 0x0a, 0x0c, 0x4f, 0x6e, 0x52, 0x61, 0x6d, 0x70, 0x52, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12,
	0x99, 0x01, 0x0a, 0x1c, 0x47, 0x65, 0x74, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x42, 0x65, 0x74, 0x77, 0x65, 0x65, 0x6e, 0x53, 0x65, 0x71, 0x4e, 0x75, 0x6d, 0x73,
	0x12, 0x3a, 0x2e, 0x6c, 0x6f, 0x6f, 0x70, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2e, 0x70, 0x62, 0x2e, 0x63, 0x63, 0x69, 0x70, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x65, 0x6e, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x65, 0x74, 0x77, 0x65, 0x65, 0x6e, 0x53, 0x65,
	0x71, 0x4e, 0x75, 0x6d, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x3b, 0x2e, 0x6c,
	0x6f, 0x6f, 0x70, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x70, 0x62, 0x2e,
	0x63, 0x63, 0x69, 0x70, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x42, 0x65, 0x74, 0x77, 0x65, 0x65, 0x6e, 0x53, 0x65, 0x71, 0x4e, 0x75, 0x6d,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x57, 0x0a, 0x0d, 0x52,
	0x6f, 0x75, 0x74, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x2c, 0x2e, 0x6c, 0x6f, 0x6f, 0x70, 0x2e, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x70, 0x62, 0x2e, 0x63, 0x63, 0x69, 0x70, 0x2e, 0x52, 0x6f, 0x75,
	0x74, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x2c, 0x2e, 0x6c, 0x6f, 0x6f, 0x70, 0x2e, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x70, 0x62, 0x2e, 0x63, 0x63, 0x69, 0x70, 0x2e,
	0x4f, 0x6e, 0x72, 0x61, 0x6d, 0x70, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5d, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x44, 0x79,
	0x6e, 0x61, 0x6d, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x2f, 0x2e, 0x6c, 0x6f, 0x6f, 0x70, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2e, 0x70, 0x62, 0x2e, 0x63, 0x63, 0x69, 0x70, 0x2e, 0x47, 0x65, 0x74, 0x44,
	0x79, 0x6e, 0x61, 0x6d, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x4f, 0x5a, 0x4d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61,
	0x63, 0x74, 0x6b, 0x69, 0x74, 0x2f, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x6c, 0x69, 0x6e, 0x6b, 0x2d,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6c, 0x6f, 0x6f, 0x70, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x62, 0x2f, 0x63, 0x63, 0x69, 0x70,
	0x3b, 0x63, 0x63, 0x69, 0x70, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_onramp_proto_rawDescOnce sync.Once
	file_onramp_proto_rawDescData = file_onramp_proto_rawDesc
)

func file_onramp_proto_rawDescGZIP() []byte {
	file_onramp_proto_rawDescOnce.Do(func() {
		file_onramp_proto_rawDescData = protoimpl.X.CompressGZIP(file_onramp_proto_rawDescData)
	})
	return file_onramp_proto_rawDescData
}

var file_onramp_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_onramp_proto_goTypes = []interface{}{
	(*GetSendRequestBetweenSeqNumsRequest)(nil),  // 0: loop.internal.pb.ccip.GetSendRequestBetweenSeqNumsRequest
	(*GetSendRequestBetweenSeqNumsResponse)(nil), // 1: loop.internal.pb.ccip.GetSendRequestBetweenSeqNumsResponse
	(*RouterAddressResponse)(nil),                // 2: loop.internal.pb.ccip.RouterAddressResponse
	(*OnrampAddressResponse)(nil),                // 3: loop.internal.pb.ccip.OnrampAddressResponse
	(*GetDynamicConfigResponse)(nil),             // 4: loop.internal.pb.ccip.GetDynamicConfigResponse
	(*OnRampDynamicConfig)(nil),                  // 5: loop.internal.pb.ccip.OnRampDynamicConfig
	(*EVM2EVMMessageWithTxMeta)(nil),             // 6: loop.internal.pb.ccip.EVM2EVMMessageWithTxMeta
	(*EVM2EVMMessage)(nil),                       // 7: loop.internal.pb.ccip.EVM2EVMMessage
	(*TxMeta)(nil),                               // 8: loop.internal.pb.ccip.TxMeta
	(*emptypb.Empty)(nil),                        // 9: google.protobuf.Empty
}
var file_onramp_proto_depIdxs = []int32{
	6, // 0: loop.internal.pb.ccip.GetSendRequestBetweenSeqNumsResponse.send_requests:type_name -> loop.internal.pb.ccip.EVM2EVMMessageWithTxMeta
	5, // 1: loop.internal.pb.ccip.GetDynamicConfigResponse.dynamic_config:type_name -> loop.internal.pb.ccip.OnRampDynamicConfig
	7, // 2: loop.internal.pb.ccip.EVM2EVMMessageWithTxMeta.message:type_name -> loop.internal.pb.ccip.EVM2EVMMessage
	8, // 3: loop.internal.pb.ccip.EVM2EVMMessageWithTxMeta.tx_meta:type_name -> loop.internal.pb.ccip.TxMeta
	0, // 4: loop.internal.pb.ccip.OnRampReader.GetSendRequestBetweenSeqNums:input_type -> loop.internal.pb.ccip.GetSendRequestBetweenSeqNumsRequest
	9, // 5: loop.internal.pb.ccip.OnRampReader.RouterAddress:input_type -> google.protobuf.Empty
	9, // 6: loop.internal.pb.ccip.OnRampReader.Address:input_type -> google.protobuf.Empty
	9, // 7: loop.internal.pb.ccip.OnRampReader.GetDynamicConfig:input_type -> google.protobuf.Empty
	1, // 8: loop.internal.pb.ccip.OnRampReader.GetSendRequestBetweenSeqNums:output_type -> loop.internal.pb.ccip.GetSendRequestBetweenSeqNumsResponse
	2, // 9: loop.internal.pb.ccip.OnRampReader.RouterAddress:output_type -> loop.internal.pb.ccip.RouterAddressResponse
	3, // 10: loop.internal.pb.ccip.OnRampReader.Address:output_type -> loop.internal.pb.ccip.OnrampAddressResponse
	4, // 11: loop.internal.pb.ccip.OnRampReader.GetDynamicConfig:output_type -> loop.internal.pb.ccip.GetDynamicConfigResponse
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_onramp_proto_init() }
func file_onramp_proto_init() {
	if File_onramp_proto != nil {
		return
	}
	file_models_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_onramp_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSendRequestBetweenSeqNumsRequest); i {
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
		file_onramp_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSendRequestBetweenSeqNumsResponse); i {
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
		file_onramp_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouterAddressResponse); i {
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
		file_onramp_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnrampAddressResponse); i {
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
		file_onramp_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDynamicConfigResponse); i {
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
		file_onramp_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnRampDynamicConfig); i {
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
		file_onramp_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EVM2EVMMessageWithTxMeta); i {
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
			RawDescriptor: file_onramp_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_onramp_proto_goTypes,
		DependencyIndexes: file_onramp_proto_depIdxs,
		MessageInfos:      file_onramp_proto_msgTypes,
	}.Build()
	File_onramp_proto = out.File
	file_onramp_proto_rawDesc = nil
	file_onramp_proto_goTypes = nil
	file_onramp_proto_depIdxs = nil
}
