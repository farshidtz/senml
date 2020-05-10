// Code generated by protoc-gen-go. DO NOT EDIT.
// source: senml.proto

package senml_protobuf

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Record struct {
	// base meta fields
	BaseName    string  `protobuf:"bytes,1,opt,name=baseName,proto3" json:"baseName,omitempty"`
	BaseTime    float64 `protobuf:"fixed64,2,opt,name=baseTime,proto3" json:"baseTime,omitempty"`
	BaseUnit    string  `protobuf:"bytes,3,opt,name=baseUnit,proto3" json:"baseUnit,omitempty"`
	BaseVersion int32   `protobuf:"varint,4,opt,name=baseVersion,proto3" json:"baseVersion,omitempty"`
	// Types that are valid to be assigned to BaseValueOptional:
	//	*Record_BaseValue
	BaseValueOptional isRecord_BaseValueOptional `protobuf_oneof:"baseValueOptional"`
	// Types that are valid to be assigned to BaseSumOptional:
	//	*Record_BaseSum
	BaseSumOptional isRecord_BaseSumOptional `protobuf_oneof:"baseSumOptional"`
	// meta fields
	Name       string  `protobuf:"bytes,7,opt,name=name,proto3" json:"name,omitempty"`
	Unit       string  `protobuf:"bytes,8,opt,name=unit,proto3" json:"unit,omitempty"`
	Time       float64 `protobuf:"fixed64,9,opt,name=time,proto3" json:"time,omitempty"`
	UpdateTime float64 `protobuf:"fixed64,10,opt,name=updateTime,proto3" json:"updateTime,omitempty"`
	// value fields
	//
	// Types that are valid to be assigned to ValueOneof:
	//	*Record_Value
	//	*Record_StringValue
	//	*Record_DataValue
	//	*Record_BoolValue
	ValueOneof isRecord_ValueOneof `protobuf_oneof:"valueOneof"`
	// Types that are valid to be assigned to SumOptional:
	//	*Record_Sum
	SumOptional          isRecord_SumOptional `protobuf_oneof:"sumOptional"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Record) Reset()         { *m = Record{} }
func (m *Record) String() string { return proto.CompactTextString(m) }
func (*Record) ProtoMessage()    {}
func (*Record) Descriptor() ([]byte, []int) {
	return fileDescriptor_d2d574953ec1dc46, []int{0}
}

func (m *Record) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Record.Unmarshal(m, b)
}
func (m *Record) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Record.Marshal(b, m, deterministic)
}
func (m *Record) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Record.Merge(m, src)
}
func (m *Record) XXX_Size() int {
	return xxx_messageInfo_Record.Size(m)
}
func (m *Record) XXX_DiscardUnknown() {
	xxx_messageInfo_Record.DiscardUnknown(m)
}

var xxx_messageInfo_Record proto.InternalMessageInfo

func (m *Record) GetBaseName() string {
	if m != nil {
		return m.BaseName
	}
	return ""
}

func (m *Record) GetBaseTime() float64 {
	if m != nil {
		return m.BaseTime
	}
	return 0
}

func (m *Record) GetBaseUnit() string {
	if m != nil {
		return m.BaseUnit
	}
	return ""
}

func (m *Record) GetBaseVersion() int32 {
	if m != nil {
		return m.BaseVersion
	}
	return 0
}

type isRecord_BaseValueOptional interface {
	isRecord_BaseValueOptional()
}

type Record_BaseValue struct {
	BaseValue float64 `protobuf:"fixed64,5,opt,name=baseValue,proto3,oneof"`
}

func (*Record_BaseValue) isRecord_BaseValueOptional() {}

func (m *Record) GetBaseValueOptional() isRecord_BaseValueOptional {
	if m != nil {
		return m.BaseValueOptional
	}
	return nil
}

func (m *Record) GetBaseValue() float64 {
	if x, ok := m.GetBaseValueOptional().(*Record_BaseValue); ok {
		return x.BaseValue
	}
	return 0
}

type isRecord_BaseSumOptional interface {
	isRecord_BaseSumOptional()
}

type Record_BaseSum struct {
	BaseSum float64 `protobuf:"fixed64,6,opt,name=baseSum,proto3,oneof"`
}

func (*Record_BaseSum) isRecord_BaseSumOptional() {}

func (m *Record) GetBaseSumOptional() isRecord_BaseSumOptional {
	if m != nil {
		return m.BaseSumOptional
	}
	return nil
}

func (m *Record) GetBaseSum() float64 {
	if x, ok := m.GetBaseSumOptional().(*Record_BaseSum); ok {
		return x.BaseSum
	}
	return 0
}

func (m *Record) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Record) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

func (m *Record) GetTime() float64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *Record) GetUpdateTime() float64 {
	if m != nil {
		return m.UpdateTime
	}
	return 0
}

type isRecord_ValueOneof interface {
	isRecord_ValueOneof()
}

type Record_Value struct {
	Value float64 `protobuf:"fixed64,11,opt,name=value,proto3,oneof"`
}

type Record_StringValue struct {
	StringValue string `protobuf:"bytes,12,opt,name=stringValue,proto3,oneof"`
}

type Record_DataValue struct {
	DataValue string `protobuf:"bytes,13,opt,name=dataValue,proto3,oneof"`
}

type Record_BoolValue struct {
	BoolValue bool `protobuf:"varint,14,opt,name=boolValue,proto3,oneof"`
}

func (*Record_Value) isRecord_ValueOneof() {}

func (*Record_StringValue) isRecord_ValueOneof() {}

func (*Record_DataValue) isRecord_ValueOneof() {}

func (*Record_BoolValue) isRecord_ValueOneof() {}

func (m *Record) GetValueOneof() isRecord_ValueOneof {
	if m != nil {
		return m.ValueOneof
	}
	return nil
}

func (m *Record) GetValue() float64 {
	if x, ok := m.GetValueOneof().(*Record_Value); ok {
		return x.Value
	}
	return 0
}

func (m *Record) GetStringValue() string {
	if x, ok := m.GetValueOneof().(*Record_StringValue); ok {
		return x.StringValue
	}
	return ""
}

func (m *Record) GetDataValue() string {
	if x, ok := m.GetValueOneof().(*Record_DataValue); ok {
		return x.DataValue
	}
	return ""
}

func (m *Record) GetBoolValue() bool {
	if x, ok := m.GetValueOneof().(*Record_BoolValue); ok {
		return x.BoolValue
	}
	return false
}

type isRecord_SumOptional interface {
	isRecord_SumOptional()
}

type Record_Sum struct {
	Sum float64 `protobuf:"fixed64,15,opt,name=sum,proto3,oneof"`
}

func (*Record_Sum) isRecord_SumOptional() {}

func (m *Record) GetSumOptional() isRecord_SumOptional {
	if m != nil {
		return m.SumOptional
	}
	return nil
}

func (m *Record) GetSum() float64 {
	if x, ok := m.GetSumOptional().(*Record_Sum); ok {
		return x.Sum
	}
	return 0
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Record) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Record_BaseValue)(nil),
		(*Record_BaseSum)(nil),
		(*Record_Value)(nil),
		(*Record_StringValue)(nil),
		(*Record_DataValue)(nil),
		(*Record_BoolValue)(nil),
		(*Record_Sum)(nil),
	}
}

// this contains the senml pack (array of records)
type Message struct {
	Pack                 []*Record `protobuf:"bytes,1,rep,name=pack,proto3" json:"pack,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_d2d574953ec1dc46, []int{1}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetPack() []*Record {
	if m != nil {
		return m.Pack
	}
	return nil
}

func init() {
	proto.RegisterType((*Record)(nil), "senml_protobuf.Record")
	proto.RegisterType((*Message)(nil), "senml_protobuf.Message")
}

func init() {
	proto.RegisterFile("senml.proto", fileDescriptor_d2d574953ec1dc46)
}

var fileDescriptor_d2d574953ec1dc46 = []byte{
	// 338 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x92, 0xcb, 0x6e, 0xea, 0x30,
	0x10, 0x86, 0x8f, 0x09, 0xd7, 0x09, 0x17, 0xe1, 0x23, 0x21, 0x8b, 0x05, 0x8a, 0x58, 0x45, 0x67,
	0xc1, 0xe2, 0x54, 0x7d, 0x01, 0x56, 0x6c, 0x5a, 0x24, 0xf7, 0xb2, 0xad, 0x0c, 0x18, 0x14, 0x35,
	0xb1, 0xa3, 0xd8, 0xe9, 0x4b, 0xf7, 0x25, 0xaa, 0x19, 0x43, 0x42, 0x77, 0xf3, 0x7f, 0xff, 0x24,
	0xf3, 0x8f, 0x35, 0x10, 0x3b, 0x6d, 0x8a, 0x7c, 0x53, 0x56, 0xd6, 0x5b, 0x3e, 0x25, 0xf1, 0x41,
	0xe2, 0x50, 0x9f, 0xd7, 0xdf, 0x11, 0xf4, 0xa5, 0x3e, 0xda, 0xea, 0xc4, 0x97, 0x30, 0x3c, 0x28,
	0xa7, 0x9f, 0x55, 0xa1, 0x05, 0x4b, 0x58, 0x3a, 0x92, 0x8d, 0xbe, 0x79, 0xaf, 0x59, 0xa1, 0x45,
	0x27, 0x61, 0x29, 0x93, 0x8d, 0xbe, 0x79, 0x6f, 0x26, 0xf3, 0x22, 0x6a, 0xbf, 0x43, 0xcd, 0x13,
	0x88, 0xb1, 0x7e, 0xd7, 0x95, 0xcb, 0xac, 0x11, 0xdd, 0x84, 0xa5, 0x3d, 0x79, 0x8f, 0xf8, 0x0a,
	0x46, 0x24, 0x55, 0x5e, 0x6b, 0xd1, 0xc3, 0x5f, 0xef, 0xfe, 0xc8, 0x16, 0xf1, 0x25, 0x0c, 0x50,
	0xbc, 0xd4, 0x85, 0xe8, 0x93, 0xcb, 0xe4, 0x0d, 0x70, 0x0e, 0x5d, 0x83, 0x69, 0x07, 0x34, 0x95,
	0x6a, 0x64, 0x35, 0x26, 0x19, 0x06, 0x86, 0x35, 0x32, 0x8f, 0xc9, 0x47, 0x94, 0x9c, 0x6a, 0xbe,
	0x02, 0xa8, 0xcb, 0x93, 0xf2, 0x61, 0x27, 0x20, 0xe7, 0x8e, 0xf0, 0x05, 0xf4, 0xbe, 0x28, 0x53,
	0x4c, 0x53, 0x3b, 0x32, 0x48, 0xbe, 0x86, 0xd8, 0xf9, 0x2a, 0x33, 0x97, 0x90, 0x78, 0x8c, 0x63,
	0x76, 0x1d, 0x79, 0x0f, 0x71, 0xa7, 0x93, 0xf2, 0x2a, 0x74, 0x4c, 0xae, 0x1d, 0x2d, 0xa2, 0x9d,
	0xad, 0xcd, 0x83, 0x3f, 0x4d, 0x58, 0x3a, 0x44, 0xbf, 0x41, 0x9c, 0x43, 0xe4, 0xea, 0x42, 0xcc,
	0x68, 0x72, 0x24, 0x51, 0x6c, 0xff, 0xc2, 0xbc, 0x79, 0x94, 0x7d, 0xe9, 0x33, 0x6b, 0x54, 0xbe,
	0x9d, 0xc3, 0xec, 0xfa, 0x16, 0x0d, 0x1a, 0x03, 0x50, 0xd0, 0xbd, 0xd1, 0xf6, 0xbc, 0x9d, 0x40,
	0xec, 0x5a, 0x73, 0xfd, 0x08, 0x83, 0x27, 0xed, 0x9c, 0xba, 0x68, 0xfe, 0x0f, 0xba, 0xa5, 0x3a,
	0x7e, 0x0a, 0x96, 0x44, 0x69, 0xfc, 0x7f, 0xb1, 0xf9, 0x7d, 0x17, 0x9b, 0x70, 0x13, 0x92, 0x7a,
	0x0e, 0x7d, 0xc2, 0x0f, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x0b, 0xda, 0x03, 0x92, 0x4a, 0x02,
	0x00, 0x00,
}