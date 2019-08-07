// Code generated by protoc-gen-go.
// source: pb.proto
// DO NOT EDIT!

/*
Package testpb is a generated protocol buffer package.

It is generated from these files:
	pb.proto

It has these top-level messages:
	Request
	Result
*/
package testpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Request struct {
	Query            *string `protobuf:"bytes,1,req,name=query" json:"query,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Request) GetQuery() string {
	if m != nil && m.Query != nil {
		return *m.Query
	}
	return ""
}

type Result struct {
	Title            *string `protobuf:"bytes,1,req,name=title" json:"title,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Result) GetTitle() string {
	if m != nil && m.Title != nil {
		return *m.Title
	}
	return ""
}

func init() {
	proto.RegisterType((*Request)(nil), "testpb.Request")
	proto.RegisterType((*Result)(nil), "testpb.Result")
}

var fileDescriptor0 = []byte{
	// 83 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x28, 0x48, 0xd2, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2b, 0x49, 0x2d, 0x2e, 0x29, 0x48, 0x52, 0x92, 0xe0, 0x62,
	0x0f, 0x4a, 0x2d, 0x2c, 0x05, 0x72, 0x84, 0x78, 0xb9, 0x58, 0x81, 0x8c, 0xa2, 0x4a, 0x09, 0x46,
	0x05, 0x26, 0x0d, 0x4e, 0x25, 0x71, 0x2e, 0xb6, 0xa0, 0xd4, 0xe2, 0xd2, 0x1c, 0xb0, 0x44, 0x49,
	0x66, 0x49, 0x4e, 0x2a, 0x44, 0x02, 0x10, 0x00, 0x00, 0xff, 0xff, 0x17, 0x03, 0x9e, 0xd8, 0x45,
	0x00, 0x00, 0x00,
}
