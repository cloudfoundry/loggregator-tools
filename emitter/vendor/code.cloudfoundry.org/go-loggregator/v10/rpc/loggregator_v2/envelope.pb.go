// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: loggregator-api/v2/envelope.proto

package loggregator_v2

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

type Log_Type int32

const (
	Log_OUT Log_Type = 0
	Log_ERR Log_Type = 1
)

// Enum value maps for Log_Type.
var (
	Log_Type_name = map[int32]string{
		0: "OUT",
		1: "ERR",
	}
	Log_Type_value = map[string]int32{
		"OUT": 0,
		"ERR": 1,
	}
)

func (x Log_Type) Enum() *Log_Type {
	p := new(Log_Type)
	*p = x
	return p
}

func (x Log_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Log_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_loggregator_api_v2_envelope_proto_enumTypes[0].Descriptor()
}

func (Log_Type) Type() protoreflect.EnumType {
	return &file_loggregator_api_v2_envelope_proto_enumTypes[0]
}

func (x Log_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Log_Type.Descriptor instead.
func (Log_Type) EnumDescriptor() ([]byte, []int) {
	return file_loggregator_api_v2_envelope_proto_rawDescGZIP(), []int{3, 0}
}

type Envelope struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp      int64             `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	SourceId       string            `protobuf:"bytes,2,opt,name=source_id,proto3" json:"source_id,omitempty"`
	InstanceId     string            `protobuf:"bytes,8,opt,name=instance_id,proto3" json:"instance_id,omitempty"`
	DeprecatedTags map[string]*Value `protobuf:"bytes,3,rep,name=deprecated_tags,proto3" json:"deprecated_tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Tags           map[string]string `protobuf:"bytes,9,rep,name=tags,proto3" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Types that are assignable to Message:
	//
	//	*Envelope_Log
	//	*Envelope_Counter
	//	*Envelope_Gauge
	//	*Envelope_Timer
	//	*Envelope_Event
	Message isEnvelope_Message `protobuf_oneof:"message"`
}

func (x *Envelope) Reset() {
	*x = Envelope{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_api_v2_envelope_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Envelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Envelope) ProtoMessage() {}

func (x *Envelope) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_api_v2_envelope_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Envelope.ProtoReflect.Descriptor instead.
func (*Envelope) Descriptor() ([]byte, []int) {
	return file_loggregator_api_v2_envelope_proto_rawDescGZIP(), []int{0}
}

func (x *Envelope) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Envelope) GetSourceId() string {
	if x != nil {
		return x.SourceId
	}
	return ""
}

func (x *Envelope) GetInstanceId() string {
	if x != nil {
		return x.InstanceId
	}
	return ""
}

func (x *Envelope) GetDeprecatedTags() map[string]*Value {
	if x != nil {
		return x.DeprecatedTags
	}
	return nil
}

func (x *Envelope) GetTags() map[string]string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (m *Envelope) GetMessage() isEnvelope_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *Envelope) GetLog() *Log {
	if x, ok := x.GetMessage().(*Envelope_Log); ok {
		return x.Log
	}
	return nil
}

func (x *Envelope) GetCounter() *Counter {
	if x, ok := x.GetMessage().(*Envelope_Counter); ok {
		return x.Counter
	}
	return nil
}

func (x *Envelope) GetGauge() *Gauge {
	if x, ok := x.GetMessage().(*Envelope_Gauge); ok {
		return x.Gauge
	}
	return nil
}

func (x *Envelope) GetTimer() *Timer {
	if x, ok := x.GetMessage().(*Envelope_Timer); ok {
		return x.Timer
	}
	return nil
}

func (x *Envelope) GetEvent() *Event {
	if x, ok := x.GetMessage().(*Envelope_Event); ok {
		return x.Event
	}
	return nil
}

type isEnvelope_Message interface {
	isEnvelope_Message()
}

type Envelope_Log struct {
	Log *Log `protobuf:"bytes,4,opt,name=log,proto3,oneof"`
}

type Envelope_Counter struct {
	Counter *Counter `protobuf:"bytes,5,opt,name=counter,proto3,oneof"`
}

type Envelope_Gauge struct {
	Gauge *Gauge `protobuf:"bytes,6,opt,name=gauge,proto3,oneof"`
}

type Envelope_Timer struct {
	Timer *Timer `protobuf:"bytes,7,opt,name=timer,proto3,oneof"`
}

type Envelope_Event struct {
	Event *Event `protobuf:"bytes,10,opt,name=event,proto3,oneof"`
}

func (*Envelope_Log) isEnvelope_Message() {}

func (*Envelope_Counter) isEnvelope_Message() {}

func (*Envelope_Gauge) isEnvelope_Message() {}

func (*Envelope_Timer) isEnvelope_Message() {}

func (*Envelope_Event) isEnvelope_Message() {}

type EnvelopeBatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Batch []*Envelope `protobuf:"bytes,1,rep,name=batch,proto3" json:"batch,omitempty"`
}

func (x *EnvelopeBatch) Reset() {
	*x = EnvelopeBatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_api_v2_envelope_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvelopeBatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvelopeBatch) ProtoMessage() {}

func (x *EnvelopeBatch) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_api_v2_envelope_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvelopeBatch.ProtoReflect.Descriptor instead.
func (*EnvelopeBatch) Descriptor() ([]byte, []int) {
	return file_loggregator_api_v2_envelope_proto_rawDescGZIP(), []int{1}
}

func (x *EnvelopeBatch) GetBatch() []*Envelope {
	if x != nil {
		return x.Batch
	}
	return nil
}

type Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//
	//	*Value_Text
	//	*Value_Integer
	//	*Value_Decimal
	Data isValue_Data `protobuf_oneof:"data"`
}

func (x *Value) Reset() {
	*x = Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_api_v2_envelope_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_api_v2_envelope_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Value.ProtoReflect.Descriptor instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_loggregator_api_v2_envelope_proto_rawDescGZIP(), []int{2}
}

func (m *Value) GetData() isValue_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *Value) GetText() string {
	if x, ok := x.GetData().(*Value_Text); ok {
		return x.Text
	}
	return ""
}

func (x *Value) GetInteger() int64 {
	if x, ok := x.GetData().(*Value_Integer); ok {
		return x.Integer
	}
	return 0
}

func (x *Value) GetDecimal() float64 {
	if x, ok := x.GetData().(*Value_Decimal); ok {
		return x.Decimal
	}
	return 0
}

type isValue_Data interface {
	isValue_Data()
}

type Value_Text struct {
	Text string `protobuf:"bytes,1,opt,name=text,proto3,oneof"`
}

type Value_Integer struct {
	Integer int64 `protobuf:"varint,2,opt,name=integer,proto3,oneof"`
}

type Value_Decimal struct {
	Decimal float64 `protobuf:"fixed64,3,opt,name=decimal,proto3,oneof"`
}

func (*Value_Text) isValue_Data() {}

func (*Value_Integer) isValue_Data() {}

func (*Value_Decimal) isValue_Data() {}

type Log struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload []byte   `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	Type    Log_Type `protobuf:"varint,2,opt,name=type,proto3,enum=loggregator.v2.Log_Type" json:"type,omitempty"`
}

func (x *Log) Reset() {
	*x = Log{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_api_v2_envelope_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Log) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Log) ProtoMessage() {}

func (x *Log) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_api_v2_envelope_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Log.ProtoReflect.Descriptor instead.
func (*Log) Descriptor() ([]byte, []int) {
	return file_loggregator_api_v2_envelope_proto_rawDescGZIP(), []int{3}
}

func (x *Log) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *Log) GetType() Log_Type {
	if x != nil {
		return x.Type
	}
	return Log_OUT
}

type Counter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Delta uint64 `protobuf:"varint,2,opt,name=delta,proto3" json:"delta,omitempty"`
	Total uint64 `protobuf:"varint,3,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *Counter) Reset() {
	*x = Counter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_api_v2_envelope_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Counter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Counter) ProtoMessage() {}

func (x *Counter) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_api_v2_envelope_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Counter.ProtoReflect.Descriptor instead.
func (*Counter) Descriptor() ([]byte, []int) {
	return file_loggregator_api_v2_envelope_proto_rawDescGZIP(), []int{4}
}

func (x *Counter) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Counter) GetDelta() uint64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *Counter) GetTotal() uint64 {
	if x != nil {
		return x.Total
	}
	return 0
}

type Gauge struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics map[string]*GaugeValue `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Gauge) Reset() {
	*x = Gauge{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_api_v2_envelope_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Gauge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Gauge) ProtoMessage() {}

func (x *Gauge) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_api_v2_envelope_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Gauge.ProtoReflect.Descriptor instead.
func (*Gauge) Descriptor() ([]byte, []int) {
	return file_loggregator_api_v2_envelope_proto_rawDescGZIP(), []int{5}
}

func (x *Gauge) GetMetrics() map[string]*GaugeValue {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type GaugeValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Unit  string  `protobuf:"bytes,1,opt,name=unit,proto3" json:"unit,omitempty"`
	Value float64 `protobuf:"fixed64,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *GaugeValue) Reset() {
	*x = GaugeValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_api_v2_envelope_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GaugeValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GaugeValue) ProtoMessage() {}

func (x *GaugeValue) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_api_v2_envelope_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GaugeValue.ProtoReflect.Descriptor instead.
func (*GaugeValue) Descriptor() ([]byte, []int) {
	return file_loggregator_api_v2_envelope_proto_rawDescGZIP(), []int{6}
}

func (x *GaugeValue) GetUnit() string {
	if x != nil {
		return x.Unit
	}
	return ""
}

func (x *GaugeValue) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type Timer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Start int64  `protobuf:"varint,2,opt,name=start,proto3" json:"start,omitempty"`
	Stop  int64  `protobuf:"varint,3,opt,name=stop,proto3" json:"stop,omitempty"`
}

func (x *Timer) Reset() {
	*x = Timer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_api_v2_envelope_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Timer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Timer) ProtoMessage() {}

func (x *Timer) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_api_v2_envelope_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Timer.ProtoReflect.Descriptor instead.
func (*Timer) Descriptor() ([]byte, []int) {
	return file_loggregator_api_v2_envelope_proto_rawDescGZIP(), []int{7}
}

func (x *Timer) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Timer) GetStart() int64 {
	if x != nil {
		return x.Start
	}
	return 0
}

func (x *Timer) GetStop() int64 {
	if x != nil {
		return x.Stop
	}
	return 0
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Body  string `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loggregator_api_v2_envelope_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_loggregator_api_v2_envelope_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_loggregator_api_v2_envelope_proto_rawDescGZIP(), []int{8}
}

func (x *Event) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Event) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

var File_loggregator_api_v2_envelope_proto protoreflect.FileDescriptor

var file_loggregator_api_v2_envelope_proto_rawDesc = []byte{
	0x0a, 0x21, 0x6c, 0x6f, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2d, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x32, 0x2f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x6c, 0x6f, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72,
	0x2e, 0x76, 0x32, 0x22, 0x81, 0x05, 0x0a, 0x08, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1c,
	0x0a, 0x09, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b,
	0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x12, 0x56,
	0x0a, 0x0f, 0x64, 0x65, 0x70, 0x72, 0x65, 0x63, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x74, 0x61, 0x67,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x72, 0x65,
	0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70,
	0x65, 0x2e, 0x44, 0x65, 0x70, 0x72, 0x65, 0x63, 0x61, 0x74, 0x65, 0x64, 0x54, 0x61, 0x67, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0f, 0x64, 0x65, 0x70, 0x72, 0x65, 0x63, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x74, 0x61, 0x67, 0x73, 0x12, 0x36, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x09,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74,
	0x6f, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x2e, 0x54,
	0x61, 0x67, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x27,
	0x0a, 0x03, 0x6c, 0x6f, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6c, 0x6f,
	0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x4c, 0x6f, 0x67,
	0x48, 0x00, 0x52, 0x03, 0x6c, 0x6f, 0x67, 0x12, 0x33, 0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x72,
	0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65,
	0x72, 0x48, 0x00, 0x52, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x2d, 0x0a, 0x05,
	0x67, 0x61, 0x75, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6c, 0x6f,
	0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x47, 0x61, 0x75,
	0x67, 0x65, 0x48, 0x00, 0x52, 0x05, 0x67, 0x61, 0x75, 0x67, 0x65, 0x12, 0x2d, 0x0a, 0x05, 0x74,
	0x69, 0x6d, 0x65, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6c, 0x6f, 0x67,
	0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x72, 0x48, 0x00, 0x52, 0x05, 0x74, 0x69, 0x6d, 0x65, 0x72, 0x12, 0x2d, 0x0a, 0x05, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6c, 0x6f, 0x67, 0x67,
	0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x48, 0x00, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x58, 0x0a, 0x13, 0x44, 0x65, 0x70,
	0x72, 0x65, 0x63, 0x61, 0x74, 0x65, 0x64, 0x54, 0x61, 0x67, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x2b, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x15, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e,
	0x76, 0x32, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x1a, 0x37, 0x0a, 0x09, 0x54, 0x61, 0x67, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x09, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x3f, 0x0a, 0x0d, 0x45, 0x6e, 0x76, 0x65, 0x6c,
	0x6f, 0x70, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x2e, 0x0a, 0x05, 0x62, 0x61, 0x74, 0x63,
	0x68, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x72, 0x65,
	0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70,
	0x65, 0x52, 0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x22, 0x5d, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x14, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x1a, 0x0a, 0x07, 0x69, 0x6e, 0x74, 0x65, 0x67,
	0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x07, 0x69, 0x6e, 0x74, 0x65,
	0x67, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x07, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x01, 0x48, 0x00, 0x52, 0x07, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x42,
	0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x67, 0x0a, 0x03, 0x4c, 0x6f, 0x67, 0x12, 0x18,
	0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x2c, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x72, 0x65, 0x67,
	0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x4c, 0x6f, 0x67, 0x2e, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x18, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x07,
	0x0a, 0x03, 0x4f, 0x55, 0x54, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x45, 0x52, 0x52, 0x10, 0x01,
	0x22, 0x49, 0x0a, 0x07, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05,
	0x64, 0x65, 0x6c, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x22, 0x9d, 0x01, 0x0a, 0x05,
	0x47, 0x61, 0x75, 0x67, 0x65, 0x12, 0x3c, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x72, 0x65, 0x67,
	0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x47, 0x61, 0x75, 0x67, 0x65, 0x2e, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x1a, 0x56, 0x0a, 0x0c, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x30, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74,
	0x6f, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x47, 0x61, 0x75, 0x67, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x36, 0x0a, 0x0a, 0x47,
	0x61, 0x75, 0x67, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x6e, 0x69,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x22, 0x45, 0x0a, 0x05, 0x54, 0x69, 0x6d, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x74, 0x6f, 0x70, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x74, 0x6f, 0x70, 0x22, 0x31, 0x0a, 0x05, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x42, 0x72, 0x0a,
	0x1f, 0x6f, 0x72, 0x67, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x72,
	0x79, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x32,
	0x42, 0x13, 0x4c, 0x6f, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x45, 0x6e, 0x76,
	0x65, 0x6c, 0x6f, 0x70, 0x65, 0x5a, 0x3a, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75,
	0x64, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x72, 0x79, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x67, 0x6f, 0x2d,
	0x6c, 0x6f, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x76, 0x39, 0x2f, 0x72,
	0x70, 0x63, 0x2f, 0x6c, 0x6f, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x76,
	0x32, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_loggregator_api_v2_envelope_proto_rawDescOnce sync.Once
	file_loggregator_api_v2_envelope_proto_rawDescData = file_loggregator_api_v2_envelope_proto_rawDesc
)

func file_loggregator_api_v2_envelope_proto_rawDescGZIP() []byte {
	file_loggregator_api_v2_envelope_proto_rawDescOnce.Do(func() {
		file_loggregator_api_v2_envelope_proto_rawDescData = protoimpl.X.CompressGZIP(file_loggregator_api_v2_envelope_proto_rawDescData)
	})
	return file_loggregator_api_v2_envelope_proto_rawDescData
}

var file_loggregator_api_v2_envelope_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_loggregator_api_v2_envelope_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_loggregator_api_v2_envelope_proto_goTypes = []interface{}{
	(Log_Type)(0),         // 0: loggregator.v2.Log.Type
	(*Envelope)(nil),      // 1: loggregator.v2.Envelope
	(*EnvelopeBatch)(nil), // 2: loggregator.v2.EnvelopeBatch
	(*Value)(nil),         // 3: loggregator.v2.Value
	(*Log)(nil),           // 4: loggregator.v2.Log
	(*Counter)(nil),       // 5: loggregator.v2.Counter
	(*Gauge)(nil),         // 6: loggregator.v2.Gauge
	(*GaugeValue)(nil),    // 7: loggregator.v2.GaugeValue
	(*Timer)(nil),         // 8: loggregator.v2.Timer
	(*Event)(nil),         // 9: loggregator.v2.Event
	nil,                   // 10: loggregator.v2.Envelope.DeprecatedTagsEntry
	nil,                   // 11: loggregator.v2.Envelope.TagsEntry
	nil,                   // 12: loggregator.v2.Gauge.MetricsEntry
}
var file_loggregator_api_v2_envelope_proto_depIdxs = []int32{
	10, // 0: loggregator.v2.Envelope.deprecated_tags:type_name -> loggregator.v2.Envelope.DeprecatedTagsEntry
	11, // 1: loggregator.v2.Envelope.tags:type_name -> loggregator.v2.Envelope.TagsEntry
	4,  // 2: loggregator.v2.Envelope.log:type_name -> loggregator.v2.Log
	5,  // 3: loggregator.v2.Envelope.counter:type_name -> loggregator.v2.Counter
	6,  // 4: loggregator.v2.Envelope.gauge:type_name -> loggregator.v2.Gauge
	8,  // 5: loggregator.v2.Envelope.timer:type_name -> loggregator.v2.Timer
	9,  // 6: loggregator.v2.Envelope.event:type_name -> loggregator.v2.Event
	1,  // 7: loggregator.v2.EnvelopeBatch.batch:type_name -> loggregator.v2.Envelope
	0,  // 8: loggregator.v2.Log.type:type_name -> loggregator.v2.Log.Type
	12, // 9: loggregator.v2.Gauge.metrics:type_name -> loggregator.v2.Gauge.MetricsEntry
	3,  // 10: loggregator.v2.Envelope.DeprecatedTagsEntry.value:type_name -> loggregator.v2.Value
	7,  // 11: loggregator.v2.Gauge.MetricsEntry.value:type_name -> loggregator.v2.GaugeValue
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_loggregator_api_v2_envelope_proto_init() }
func file_loggregator_api_v2_envelope_proto_init() {
	if File_loggregator_api_v2_envelope_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_loggregator_api_v2_envelope_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Envelope); i {
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
		file_loggregator_api_v2_envelope_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvelopeBatch); i {
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
		file_loggregator_api_v2_envelope_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Value); i {
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
		file_loggregator_api_v2_envelope_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Log); i {
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
		file_loggregator_api_v2_envelope_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Counter); i {
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
		file_loggregator_api_v2_envelope_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Gauge); i {
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
		file_loggregator_api_v2_envelope_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GaugeValue); i {
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
		file_loggregator_api_v2_envelope_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Timer); i {
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
		file_loggregator_api_v2_envelope_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
	file_loggregator_api_v2_envelope_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Envelope_Log)(nil),
		(*Envelope_Counter)(nil),
		(*Envelope_Gauge)(nil),
		(*Envelope_Timer)(nil),
		(*Envelope_Event)(nil),
	}
	file_loggregator_api_v2_envelope_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Value_Text)(nil),
		(*Value_Integer)(nil),
		(*Value_Decimal)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_loggregator_api_v2_envelope_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_loggregator_api_v2_envelope_proto_goTypes,
		DependencyIndexes: file_loggregator_api_v2_envelope_proto_depIdxs,
		EnumInfos:         file_loggregator_api_v2_envelope_proto_enumTypes,
		MessageInfos:      file_loggregator_api_v2_envelope_proto_msgTypes,
	}.Build()
	File_loggregator_api_v2_envelope_proto = out.File
	file_loggregator_api_v2_envelope_proto_rawDesc = nil
	file_loggregator_api_v2_envelope_proto_goTypes = nil
	file_loggregator_api_v2_envelope_proto_depIdxs = nil
}
