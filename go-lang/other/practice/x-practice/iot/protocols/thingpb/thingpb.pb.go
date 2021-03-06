// Code generated by protoc-gen-go. DO NOT EDIT.
// source: thingpb.proto

package thingpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Head struct {
	MsgName              string   `protobuf:"bytes,1,opt,name=MsgName" json:"MsgName,omitempty"`
	Sequence             int32    `protobuf:"varint,2,opt,name=Sequence" json:"Sequence,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Head) Reset()         { *m = Head{} }
func (m *Head) String() string { return proto.CompactTextString(m) }
func (*Head) ProtoMessage()    {}
func (*Head) Descriptor() ([]byte, []int) {
	return fileDescriptor_thingpb_5672aa721159fdbd, []int{0}
}
func (m *Head) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Head.Unmarshal(m, b)
}
func (m *Head) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Head.Marshal(b, m, deterministic)
}
func (dst *Head) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Head.Merge(dst, src)
}
func (m *Head) XXX_Size() int {
	return xxx_messageInfo_Head.Size(m)
}
func (m *Head) XXX_DiscardUnknown() {
	xxx_messageInfo_Head.DiscardUnknown(m)
}

var xxx_messageInfo_Head proto.InternalMessageInfo

func (m *Head) GetMsgName() string {
	if m != nil {
		return m.MsgName
	}
	return ""
}

func (m *Head) GetSequence() int32 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

type ErrorNotice struct {
	Code                 int32    `protobuf:"varint,1,opt,name=Code" json:"Code,omitempty"`
	Details              string   `protobuf:"bytes,2,opt,name=Details" json:"Details,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ErrorNotice) Reset()         { *m = ErrorNotice{} }
func (m *ErrorNotice) String() string { return proto.CompactTextString(m) }
func (*ErrorNotice) ProtoMessage()    {}
func (*ErrorNotice) Descriptor() ([]byte, []int) {
	return fileDescriptor_thingpb_5672aa721159fdbd, []int{1}
}
func (m *ErrorNotice) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ErrorNotice.Unmarshal(m, b)
}
func (m *ErrorNotice) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ErrorNotice.Marshal(b, m, deterministic)
}
func (dst *ErrorNotice) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrorNotice.Merge(dst, src)
}
func (m *ErrorNotice) XXX_Size() int {
	return xxx_messageInfo_ErrorNotice.Size(m)
}
func (m *ErrorNotice) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrorNotice.DiscardUnknown(m)
}

var xxx_messageInfo_ErrorNotice proto.InternalMessageInfo

func (m *ErrorNotice) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *ErrorNotice) GetDetails() string {
	if m != nil {
		return m.Details
	}
	return ""
}

func init() {
	proto.RegisterType((*Head)(nil), "thingpb.Head")
	proto.RegisterType((*ErrorNotice)(nil), "thingpb.ErrorNotice")
}

func init() { proto.RegisterFile("thingpb.proto", fileDescriptor_thingpb_5672aa721159fdbd) }

var fileDescriptor_thingpb_5672aa721159fdbd = []byte{
	// 139 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0xc9, 0xc8, 0xcc,
	0x4b, 0x2f, 0x48, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x72, 0x95, 0x6c, 0xb8,
	0x58, 0x3c, 0x52, 0x13, 0x53, 0x84, 0x24, 0xb8, 0xd8, 0x7d, 0x8b, 0xd3, 0xfd, 0x12, 0x73, 0x53,
	0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x60, 0x5c, 0x21, 0x29, 0x2e, 0x8e, 0xe0, 0xd4, 0xc2,
	0xd2, 0xd4, 0xbc, 0xe4, 0x54, 0x09, 0x26, 0x05, 0x46, 0x0d, 0xd6, 0x20, 0x38, 0x5f, 0xc9, 0x9a,
	0x8b, 0xdb, 0xb5, 0xa8, 0x28, 0xbf, 0xc8, 0x2f, 0xbf, 0x24, 0x33, 0x39, 0x55, 0x48, 0x88, 0x8b,
	0xc5, 0x39, 0x3f, 0x05, 0x62, 0x02, 0x6b, 0x10, 0x98, 0x0d, 0x32, 0xd8, 0x25, 0xb5, 0x24, 0x31,
	0x33, 0xa7, 0x18, 0xac, 0x9b, 0x33, 0x08, 0xc6, 0x4d, 0x62, 0x03, 0x3b, 0xc5, 0x18, 0x10, 0x00,
	0x00, 0xff, 0xff, 0xce, 0x17, 0x0e, 0x5f, 0x9b, 0x00, 0x00, 0x00,
}
