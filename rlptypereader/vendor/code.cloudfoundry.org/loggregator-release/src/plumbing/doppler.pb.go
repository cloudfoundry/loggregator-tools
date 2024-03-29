// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.4
// source: loggregator-release/src/plumbing/doppler.proto

package plumbing

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

type EnvelopeData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *EnvelopeData) Reset() {
	*x = EnvelopeData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_release_src_plumbing_doppler_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvelopeData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvelopeData) ProtoMessage() {}

func (x *EnvelopeData) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_release_src_plumbing_doppler_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvelopeData.ProtoReflect.Descriptor instead.
func (*EnvelopeData) Descriptor() ([]byte, []int) {
	return file_loggregator_release_src_plumbing_doppler_proto_rawDescGZIP(), []int{0}
}

func (x *EnvelopeData) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type PushResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PushResponse) Reset() {
	*x = PushResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_release_src_plumbing_doppler_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushResponse) ProtoMessage() {}

func (x *PushResponse) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_release_src_plumbing_doppler_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushResponse.ProtoReflect.Descriptor instead.
func (*PushResponse) Descriptor() ([]byte, []int) {
	return file_loggregator_release_src_plumbing_doppler_proto_rawDescGZIP(), []int{1}
}

type SubscriptionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShardID string  `protobuf:"bytes,1,opt,name=shardID,proto3" json:"shardID,omitempty"`
	Filter  *Filter `protobuf:"bytes,2,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (x *SubscriptionRequest) Reset() {
	*x = SubscriptionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_release_src_plumbing_doppler_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscriptionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscriptionRequest) ProtoMessage() {}

func (x *SubscriptionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_release_src_plumbing_doppler_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscriptionRequest.ProtoReflect.Descriptor instead.
func (*SubscriptionRequest) Descriptor() ([]byte, []int) {
	return file_loggregator_release_src_plumbing_doppler_proto_rawDescGZIP(), []int{2}
}

func (x *SubscriptionRequest) GetShardID() string {
	if x != nil {
		return x.ShardID
	}
	return ""
}

func (x *SubscriptionRequest) GetFilter() *Filter {
	if x != nil {
		return x.Filter
	}
	return nil
}

type Filter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AppID string `protobuf:"bytes,1,opt,name=appID,proto3" json:"appID,omitempty"`
	// Types that are assignable to Message:
	//	*Filter_Log
	//	*Filter_Metric
	Message isFilter_Message `protobuf_oneof:"Message"`
}

func (x *Filter) Reset() {
	*x = Filter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_release_src_plumbing_doppler_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Filter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Filter) ProtoMessage() {}

func (x *Filter) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_release_src_plumbing_doppler_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Filter.ProtoReflect.Descriptor instead.
func (*Filter) Descriptor() ([]byte, []int) {
	return file_loggregator_release_src_plumbing_doppler_proto_rawDescGZIP(), []int{3}
}

func (x *Filter) GetAppID() string {
	if x != nil {
		return x.AppID
	}
	return ""
}

func (m *Filter) GetMessage() isFilter_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *Filter) GetLog() *LogFilter {
	if x, ok := x.GetMessage().(*Filter_Log); ok {
		return x.Log
	}
	return nil
}

func (x *Filter) GetMetric() *MetricFilter {
	if x, ok := x.GetMessage().(*Filter_Metric); ok {
		return x.Metric
	}
	return nil
}

type isFilter_Message interface {
	isFilter_Message()
}

type Filter_Log struct {
	Log *LogFilter `protobuf:"bytes,2,opt,name=log,proto3,oneof"`
}

type Filter_Metric struct {
	Metric *MetricFilter `protobuf:"bytes,3,opt,name=metric,proto3,oneof"`
}

func (*Filter_Log) isFilter_Message() {}

func (*Filter_Metric) isFilter_Message() {}

type LogFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LogFilter) Reset() {
	*x = LogFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_release_src_plumbing_doppler_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogFilter) ProtoMessage() {}

func (x *LogFilter) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_release_src_plumbing_doppler_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogFilter.ProtoReflect.Descriptor instead.
func (*LogFilter) Descriptor() ([]byte, []int) {
	return file_loggregator_release_src_plumbing_doppler_proto_rawDescGZIP(), []int{4}
}

type MetricFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MetricFilter) Reset() {
	*x = MetricFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_release_src_plumbing_doppler_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricFilter) ProtoMessage() {}

func (x *MetricFilter) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_release_src_plumbing_doppler_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricFilter.ProtoReflect.Descriptor instead.
func (*MetricFilter) Descriptor() ([]byte, []int) {
	return file_loggregator_release_src_plumbing_doppler_proto_rawDescGZIP(), []int{5}
}

// Note: Ideally this would be EnvelopeData but for the time being we do not
// want to pay the cost of planning an upgrade path for this to be renamed.
type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_release_src_plumbing_doppler_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_release_src_plumbing_doppler_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_loggregator_release_src_plumbing_doppler_proto_rawDescGZIP(), []int{6}
}

func (x *Response) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type BatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload [][]byte `protobuf:"bytes,1,rep,name=payload,proto3" json:"payload,omitempty"`
}

func (x *BatchResponse) Reset() {
	*x = BatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_release_src_plumbing_doppler_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchResponse) ProtoMessage() {}

func (x *BatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_release_src_plumbing_doppler_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchResponse.ProtoReflect.Descriptor instead.
func (*BatchResponse) Descriptor() ([]byte, []int) {
	return file_loggregator_release_src_plumbing_doppler_proto_rawDescGZIP(), []int{7}
}

func (x *BatchResponse) GetPayload() [][]byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

var File_loggregator_release_src_plumbing_doppler_proto protoreflect.FileDescriptor

var file_loggregator_release_src_plumbing_doppler_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2d, 0x72, 0x65,
	0x6c, 0x65, 0x61, 0x73, 0x65, 0x2f, 0x73, 0x72, 0x63, 0x2f, 0x70, 0x6c, 0x75, 0x6d, 0x62, 0x69,
	0x6e, 0x67, 0x2f, 0x64, 0x6f, 0x70, 0x70, 0x6c, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x70, 0x6c, 0x75, 0x6d, 0x62, 0x69, 0x6e, 0x67, 0x22, 0x28, 0x0a, 0x0c, 0x45, 0x6e,
	0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x22, 0x0e, 0x0a, 0x0c, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x59, 0x0a, 0x13, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x68, 0x61, 0x72, 0x64, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x68,
	0x61, 0x72, 0x64, 0x49, 0x44, 0x12, 0x28, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x6c, 0x75, 0x6d, 0x62, 0x69, 0x6e, 0x67,
	0x2e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x22,
	0x84, 0x01, 0x0a, 0x06, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x70,
	0x70, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x70, 0x70, 0x49, 0x44,
	0x12, 0x27, 0x0a, 0x03, 0x6c, 0x6f, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x70, 0x6c, 0x75, 0x6d, 0x62, 0x69, 0x6e, 0x67, 0x2e, 0x4c, 0x6f, 0x67, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x48, 0x00, 0x52, 0x03, 0x6c, 0x6f, 0x67, 0x12, 0x30, 0x0a, 0x06, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x6c, 0x75, 0x6d,
	0x62, 0x69, 0x6e, 0x67, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x46, 0x69, 0x6c, 0x74, 0x65,
	0x72, 0x48, 0x00, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x42, 0x09, 0x0a, 0x07, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x0b, 0x0a, 0x09, 0x4c, 0x6f, 0x67, 0x46, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x22, 0x0e, 0x0a, 0x0c, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x46, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x22, 0x24, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x29, 0x0a, 0x0d, 0x42, 0x61, 0x74,
	0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x32, 0x9b, 0x01, 0x0a, 0x07, 0x44, 0x6f, 0x70, 0x70, 0x6c, 0x65, 0x72,
	0x12, 0x42, 0x0a, 0x09, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12, 0x1d, 0x2e,
	0x70, 0x6c, 0x75, 0x6d, 0x62, 0x69, 0x6e, 0x67, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x70,
	0x6c, 0x75, 0x6d, 0x62, 0x69, 0x6e, 0x67, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x30, 0x01, 0x12, 0x4c, 0x0a, 0x0e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x53, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12, 0x1d, 0x2e, 0x70, 0x6c, 0x75, 0x6d, 0x62, 0x69, 0x6e,
	0x67, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x70, 0x6c, 0x75, 0x6d, 0x62, 0x69, 0x6e, 0x67,
	0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x30, 0x01, 0x32, 0x4f, 0x0a, 0x0f, 0x44, 0x6f, 0x70, 0x70, 0x6c, 0x65, 0x72, 0x49, 0x6e, 0x67,
	0x65, 0x73, 0x74, 0x6f, 0x72, 0x12, 0x3c, 0x0a, 0x06, 0x50, 0x75, 0x73, 0x68, 0x65, 0x72, 0x12,
	0x16, 0x2e, 0x70, 0x6c, 0x75, 0x6d, 0x62, 0x69, 0x6e, 0x67, 0x2e, 0x45, 0x6e, 0x76, 0x65, 0x6c,
	0x6f, 0x70, 0x65, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x16, 0x2e, 0x70, 0x6c, 0x75, 0x6d, 0x62, 0x69,
	0x6e, 0x67, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x28, 0x01, 0x42, 0x2c, 0x5a, 0x2a, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75,
	0x64, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x72, 0x79, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x6c, 0x6f, 0x67,
	0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x70, 0x6c, 0x75, 0x6d, 0x62, 0x69, 0x6e,
	0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_loggregator_release_src_plumbing_doppler_proto_rawDescOnce sync.Once
	file_loggregator_release_src_plumbing_doppler_proto_rawDescData = file_loggregator_release_src_plumbing_doppler_proto_rawDesc
)

func file_loggregator_release_src_plumbing_doppler_proto_rawDescGZIP() []byte {
	file_loggregator_release_src_plumbing_doppler_proto_rawDescOnce.Do(func() {
		file_loggregator_release_src_plumbing_doppler_proto_rawDescData = protoimpl.X.CompressGZIP(file_loggregator_release_src_plumbing_doppler_proto_rawDescData)
	})
	return file_loggregator_release_src_plumbing_doppler_proto_rawDescData
}

var file_loggregator_release_src_plumbing_doppler_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_loggregator_release_src_plumbing_doppler_proto_goTypes = []interface{}{
	(*EnvelopeData)(nil),        // 0: plumbing.EnvelopeData
	(*PushResponse)(nil),        // 1: plumbing.PushResponse
	(*SubscriptionRequest)(nil), // 2: plumbing.SubscriptionRequest
	(*Filter)(nil),              // 3: plumbing.Filter
	(*LogFilter)(nil),           // 4: plumbing.LogFilter
	(*MetricFilter)(nil),        // 5: plumbing.MetricFilter
	(*Response)(nil),            // 6: plumbing.Response
	(*BatchResponse)(nil),       // 7: plumbing.BatchResponse
}
var file_loggregator_release_src_plumbing_doppler_proto_depIdxs = []int32{
	3, // 0: plumbing.SubscriptionRequest.filter:type_name -> plumbing.Filter
	4, // 1: plumbing.Filter.log:type_name -> plumbing.LogFilter
	5, // 2: plumbing.Filter.metric:type_name -> plumbing.MetricFilter
	2, // 3: plumbing.Doppler.Subscribe:input_type -> plumbing.SubscriptionRequest
	2, // 4: plumbing.Doppler.BatchSubscribe:input_type -> plumbing.SubscriptionRequest
	0, // 5: plumbing.DopplerIngestor.Pusher:input_type -> plumbing.EnvelopeData
	6, // 6: plumbing.Doppler.Subscribe:output_type -> plumbing.Response
	7, // 7: plumbing.Doppler.BatchSubscribe:output_type -> plumbing.BatchResponse
	1, // 8: plumbing.DopplerIngestor.Pusher:output_type -> plumbing.PushResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_loggregator_release_src_plumbing_doppler_proto_init() }
func file_loggregator_release_src_plumbing_doppler_proto_init() {
	if File_loggregator_release_src_plumbing_doppler_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_loggregator_release_src_plumbing_doppler_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvelopeData); i {
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
		file_loggregator_release_src_plumbing_doppler_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushResponse); i {
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
		file_loggregator_release_src_plumbing_doppler_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscriptionRequest); i {
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
		file_loggregator_release_src_plumbing_doppler_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Filter); i {
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
		file_loggregator_release_src_plumbing_doppler_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogFilter); i {
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
		file_loggregator_release_src_plumbing_doppler_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MetricFilter); i {
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
		file_loggregator_release_src_plumbing_doppler_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_loggregator_release_src_plumbing_doppler_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchResponse); i {
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
	file_loggregator_release_src_plumbing_doppler_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*Filter_Log)(nil),
		(*Filter_Metric)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_loggregator_release_src_plumbing_doppler_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_loggregator_release_src_plumbing_doppler_proto_goTypes,
		DependencyIndexes: file_loggregator_release_src_plumbing_doppler_proto_depIdxs,
		MessageInfos:      file_loggregator_release_src_plumbing_doppler_proto_msgTypes,
	}.Build()
	File_loggregator_release_src_plumbing_doppler_proto = out.File
	file_loggregator_release_src_plumbing_doppler_proto_rawDesc = nil
	file_loggregator_release_src_plumbing_doppler_proto_goTypes = nil
	file_loggregator_release_src_plumbing_doppler_proto_depIdxs = nil
}
