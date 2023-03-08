package message

import (
	"fmt"
	"hash/crc64"
	"io"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/ironzhang/gamecloud/pkg/message/packet"
)

// NewMessage 根据消息名称构建消息
func NewMessage(name string) (proto.Message, error) {
	t := proto.MessageType(name)
	if t == nil {
		return nil, fmt.Errorf("%q message type not found", name)
	}
	i := reflect.New(t.Elem()).Interface()
	return i.(proto.Message), nil
}

var tab = crc64.MakeTable(crc64.ISO)

// ReadMessage 读取消息
func ReadMessage(r io.Reader) (msg proto.Message, err error) {
	var p packet.Packet
	if err = readPacket(r, &p); err != nil {
		return nil, err
	}
	if sum := crc64.Checksum(p.Payload, tab); sum != p.Checksum {
		return nil, fmt.Errorf("checksum is wrong, %d != %d", sum, p.Checksum)
	}
	msg, err = NewMessage(p.Name)
	if err != nil {
		return nil, err
	}
	if err = proto.Unmarshal(p.Payload, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

// WriteMessage 写入消息
func WriteMessage(w io.Writer, msg proto.Message) (err error) {
	var p packet.Packet
	p.Name = proto.MessageName(msg)
	if p.Payload, err = proto.Marshal(msg); err != nil {
		return err
	}
	p.Checksum = crc64.Checksum(p.Payload, tab)
	return writePacket(w, &p)
}
