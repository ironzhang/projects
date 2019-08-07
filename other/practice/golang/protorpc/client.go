package protorpc

import (
	"fmt"
	"io"
	"net"
	"net/rpc"

	"github.com/golang/protobuf/proto"

	"github.com/ironzhang/golang/protorpc/protorpcpb"
)

type clientCodec struct {
	rwc io.ReadWriteCloser
}

func (c *clientCodec) WriteRequest(r *rpc.Request, body interface{}) error {
	var param proto.Message
	if body != nil {
		var ok bool
		if param, ok = body.(proto.Message); !ok {
			return fmt.Errorf("%T does not implement proto.Message", body)
		}
	}

	var req protorpcpb.Request
	req.Method = proto.String(r.ServiceMethod)
	req.Seq = proto.Uint64(r.Seq)
	return writeMessage(c.rwc, &req, param)
}

func (c *clientCodec) ReadResponseHeader(r *rpc.Response) error {
	var rsp protorpcpb.Response
	if err := readFrame(c.rwc, &rsp); err != nil {
		return err
	}
	r.ServiceMethod = rsp.GetMethod()
	r.Seq = rsp.GetSeq()
	r.Error = rsp.GetErr()
	return nil
}

func (c *clientCodec) ReadResponseBody(body interface{}) error {
	var param proto.Message
	if body != nil {
		var ok bool
		if param, ok = body.(proto.Message); !ok {
			return fmt.Errorf("%T does not implement proto.Message", body)
		}
	}
	return readFrame(c.rwc, param)
}

func (c *clientCodec) Close() error {
	return c.rwc.Close()
}

func NewClientCodec(conn io.ReadWriteCloser) rpc.ClientCodec {
	return &clientCodec{rwc: conn}
}

func NewClient(conn io.ReadWriteCloser) *rpc.Client {
	return rpc.NewClientWithCodec(NewClientCodec(conn))
}

func Dial(network, address string) (*rpc.Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return NewClient(conn), nil
}
