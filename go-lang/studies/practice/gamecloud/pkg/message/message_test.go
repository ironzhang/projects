package message

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/ironzhang/gamecloud/pkg/message/testpb"
)

func TestNewMessage(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "testpb.Header",
			err:  nil,
		},
		{
			name: "testpb.GetAccessPointReq",
			err:  nil,
		},
		{
			name: "testpb.GetAccessPointRsp",
			err:  nil,
		},
		{
			name: "packet.Packet",
			err:  nil,
		},
		{
			name: "Packet",
			err:  fmt.Errorf("%q message type not found", "Packet"),
		},
		{
			name: "testpb.Transport",
			err:  fmt.Errorf("%q message type not found", "testpb.Transport"),
		},
	}

	for i, test := range tests {
		_, err := NewMessage(test.name)
		if !reflect.DeepEqual(err, test.err) {
			t.Errorf("testcase(%d): %v != %v", i, err, test.err)
		}
		//	if msg != nil {
		//		t.Logf("testcase(%d): Type(%[2]T), Value(%[2]v)\n", i, msg)
		//	}
	}
}

func TestReadWriteMessage(t *testing.T) {
	wMessages := []proto.Message{
		&testpb.Header{Seq: 1},
		&testpb.GetAccessPointReq{Header: &testpb.Header{Seq: 2}, AccessType: "MQTT"},
		&testpb.GetAccessPointRsp{Addrs: []string{"127.0.0.1", "127.0.0.2"}},
	}
	rMessages := []proto.Message{}

	var buf bytes.Buffer
	for _, msg := range wMessages {
		if err := WriteMessage(&buf, msg); err != nil {
			t.Fatalf("WriteMessage: %v", err)
		}
	}
	for {
		msg, err := ReadMessage(&buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatalf("ReadMessage: %v", err)
		}
		rMessages = append(rMessages, msg)
	}
	if !reflect.DeepEqual(rMessages, wMessages) {
		t.Errorf("%v != %v", rMessages, wMessages)
	}
}
