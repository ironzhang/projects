package thingmsg

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"

	pb "github.com/ironzhang/practice/x-practice/iot/protocols/thingpb"
)

func TestNewMessageBody(t *testing.T) {
	tests := []struct {
		name string
		msg  proto.Message
		err  bool
	}{
		{"Head", nil, true},
		{"ErrorNotice", nil, true},
		{"thingpb.Head", &pb.Head{}, false},
		{"thingpb.ErrorNotice", &pb.ErrorNotice{}, false},
	}

	for i, tt := range tests {
		msg, err := newMessageBody(tt.name)
		if tt.err {
			if err == nil {
				t.Fatalf("%d: new message body: err is nil", i, err)
			} else {
				t.Logf("%d: new message body: err[%v] is not nil", i, err)
			}
		} else {
			if err != nil {
				t.Fatalf("%d: new message body: %v", i, err)
			}
			if got, want := reflect.TypeOf(msg), reflect.TypeOf(tt.msg); got != want {
				t.Fatalf("%d: message type: %v != %v", i, got, want)
			}
		}
	}
}

func TestMessage(t *testing.T) {
	tests := []struct {
		msg Message
		err bool
	}{
		{
			msg: Message{
				Head: pb.Head{MsgName: "thingpb.ErrorNotice"},
				Body: &pb.ErrorNotice{},
			},
			err: false,
		},
		{
			msg: Message{
				Head: pb.Head{MsgName: "thingpb.ErrorNotice"},
				Body: &pb.ErrorNotice{Code: 1, Details: "details"},
			},
			err: false,
		},
		{
			msg: Message{
				Head: pb.Head{MsgName: "thingpb.ErrorNotice1"},
				Body: &pb.ErrorNotice{Code: 1},
			},
			err: true,
		},
	}
	for i, tt := range tests {
		data, err := tt.msg.Encode()
		if err != nil {
			t.Fatalf("%d: encode: %v", i, err)
		}

		var m Message
		err = m.Decode(data)
		if tt.err {
			if err == nil {
				t.Fatalf("%d: decode: err is nil", i)
			} else {
				t.Logf("%d: decode: err[%v] is not nil", i, err)
			}
		} else {
			if err != nil {
				t.Fatalf("%d: decode: %v", i, err)
			}
			if got, want := m.String(), tt.msg.String(); got != want {
				t.Fatalf("%d: message: %v != %v", i, got, want)
			} else {
				t.Logf("%d: message string: %v", i, got)
			}
		}
	}
}

func TestPacket(t *testing.T) {
	messages := []Message{
		{
			Head: pb.Head{MsgName: "thingpb.ErrorNotice"},
			Body: &pb.ErrorNotice{},
		},
		//		{
		//			Head: pb.Head{MsgName: "thingpb.ErrorNotice"},
		//			Body: &pb.ErrorNotice{Code: 1},
		//		},
	}

	var buf bytes.Buffer
	for i, msg := range messages {
		if err := Write(&buf, &msg); err != nil {
			t.Fatalf("%d: write: %v", i, err)
		}
	}
	t.Logf("bytes: %d: %v", buf.Len(), buf.Bytes())
	for i, msg := range messages {
		m, err := Read(&buf)
		if err != nil {
			t.Fatalf("%d: read: %v", i, err)
		}
		if got, want := m.String(), msg.String(); got != want {
			t.Fatalf("%d: message: %v != %v", i, got, want)
		} else {
			t.Logf("%d: message string: %v", i, got)
		}
	}
}
