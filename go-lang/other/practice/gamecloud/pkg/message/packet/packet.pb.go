// Code generated by protoc-gen-go.
// source: packet.proto
// DO NOT EDIT!

/*
Package packet is a generated protocol buffer package.

It is generated from these files:
	packet.proto

It has these top-level messages:
	Packet
*/
package packet

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

type Packet struct {
	Name     string `protobuf:"bytes,1,opt,name=Name" json:"Name,omitempty"`
	Payload  []byte `protobuf:"bytes,2,opt,name=Payload,proto3" json:"Payload,omitempty"`
	Checksum uint64 `protobuf:"varint,3,opt,name=Checksum" json:"Checksum,omitempty"`
}

func (m *Packet) Reset()                    { *m = Packet{} }
func (m *Packet) String() string            { return proto.CompactTextString(m) }
func (*Packet) ProtoMessage()               {}
func (*Packet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Packet) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Packet) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Packet) GetChecksum() uint64 {
	if m != nil {
		return m.Checksum
	}
	return 0
}

func init() {
	proto.RegisterType((*Packet)(nil), "packet.Packet")
}

func init() { proto.RegisterFile("packet.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 109 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x48, 0x4c, 0xce,
	0x4e, 0x2d, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x83, 0xf0, 0x94, 0x82, 0xb8, 0xd8,
	0x02, 0xc0, 0x2c, 0x21, 0x21, 0x2e, 0x16, 0xbf, 0xc4, 0xdc, 0x54, 0x09, 0x46, 0x05, 0x46, 0x0d,
	0xce, 0x20, 0x30, 0x5b, 0x48, 0x82, 0x8b, 0x3d, 0x20, 0xb1, 0x32, 0x27, 0x3f, 0x31, 0x45, 0x82,
	0x09, 0x28, 0xcc, 0x13, 0x04, 0xe3, 0x0a, 0x49, 0x71, 0x71, 0x38, 0x67, 0xa4, 0x26, 0x67, 0x17,
	0x97, 0xe6, 0x4a, 0x30, 0x03, 0xa5, 0x58, 0x82, 0xe0, 0xfc, 0x24, 0x36, 0xb0, 0x15, 0xc6, 0x80,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x0d, 0x23, 0x33, 0xe7, 0x72, 0x00, 0x00, 0x00,
}