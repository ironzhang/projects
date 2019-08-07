package protorpc

import (
	"bytes"
	"fmt"
	"net/rpc"
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/ironzhang/golang/protorpc/testpb"
)

type Buffer struct {
	bytes.Buffer
}

func (b *Buffer) Close() error {
	return nil
}

func writeRequest(c *clientCodec) error {
	head := rpc.Request{}
	head.ServiceMethod = "TestFunc"
	head.Seq = 1

	body := testpb.Request{}
	body.Query = proto.String("ping www.tencent.com")
	return c.WriteRequest(&head, &body)
}

func readRequest(c *serverCodec) error {
	head := rpc.Request{}
	if err := c.ReadRequestHeader(&head); err != nil {
		return err
	}

	body := testpb.Request{}
	if err := c.ReadRequestBody(&body); err != nil {
		return err
	}

	fmt.Printf("head: %v, body: %v\n", head, body.GetQuery())
	return nil
}

func writeResponse(c *serverCodec) error {
	head := rpc.Response{}
	head.ServiceMethod = "TestFunc"
	head.Seq = 2
	head.Error = "call failed"

	body := testpb.Result{}
	body.Title = proto.String("www.zhihu.com")
	return c.WriteResponse(&head, &body)
}

func readResponse(c *clientCodec) error {
	head := rpc.Response{}
	if err := c.ReadResponseHeader(&head); err != nil {
		return err
	}

	body := testpb.Result{}
	if err := c.ReadResponseBody(&body); err != nil {
		return err
	}

	fmt.Printf("head: %v, body: %v\n", head, body.GetTitle())
	return nil
}

func TestCodec(t *testing.T) {
	buffer := &Buffer{}
	ccodec := clientCodec{rwc: buffer}
	scodec := serverCodec{rwc: buffer}

	if err := writeRequest(&ccodec); err != nil {
		t.Errorf("clientcodec write request error, %v", err)
		return
	}

	if err := readRequest(&scodec); err != nil {
		t.Errorf("servercodec read request error, %v", err)
		return
	}

	if err := writeResponse(&scodec); err != nil {
		t.Errorf("servercodec write response error, %v", err)
		return
	}

	if err := readResponse(&ccodec); err != nil {
		t.Errorf("clientcodec read response error, %v", err)
		return
	}
}
