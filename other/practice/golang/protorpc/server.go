package protorpc

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"

	"github.com/golang/protobuf/proto"

	"github.com/ironzhang/golang/protorpc/protorpcpb"
)

type serverCodec struct {
	rwc io.ReadWriteCloser
}

func (c *serverCodec) ReadRequestHeader(r *rpc.Request) error {
	var req protorpcpb.Request
	if err := readFrame(c.rwc, &req); err != nil {
		return err
	}
	r.ServiceMethod = req.GetMethod()
	r.Seq = req.GetSeq()
	return nil
}

func (c *serverCodec) ReadRequestBody(body interface{}) error {
	var param proto.Message
	if body != nil {
		var ok bool
		if param, ok = body.(proto.Message); !ok {
			return fmt.Errorf("%T does not implement proto.Message", body)
		}
	}
	return readFrame(c.rwc, param)
}

func (c *serverCodec) WriteResponse(r *rpc.Response, body interface{}) error {
	var param proto.Message
	if body != nil {
		var ok bool
		if param, ok = body.(proto.Message); !ok {
			return fmt.Errorf("%T does not implement proto.Message", body)
		}
	}

	var rsp protorpcpb.Response
	rsp.Method = proto.String(r.ServiceMethod)
	rsp.Seq = proto.Uint64(r.Seq)
	rsp.Err = proto.String(r.Error)
	return writeMessage(c.rwc, &rsp, param)
}

func (c *serverCodec) Close() error {
	return c.rwc.Close()
}

func NewServerCodec(conn io.ReadWriteCloser) rpc.ServerCodec {
	return &serverCodec{rwc: conn}
}

func ServeConn(conn io.ReadWriteCloser) {
	rpc.ServeCodec(NewServerCodec(conn))
}

func Serve(s *rpc.Server, l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("rpc.Serve: accept:", err.Error())
		}
		go s.ServeCodec(NewServerCodec(conn))
	}
}
