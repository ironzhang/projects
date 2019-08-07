package thingmsg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"math"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/ironzhang/pearls/endian"

	pb "github.com/ironzhang/practice/x-practice/iot/protocols/thingpb"
)

var (
	ErrInvalidChecksum          = errors.New("invalid checksum")
	ErrInvalidPacketLength      = errors.New("invalid packet length")
	ErrInvalidMessageHeadLength = errors.New("invalid message head length")
)

/*
Packet format
+----------+
| Length   | Varint32, 1 - 5 bytes
+----------+
| Payload  | [Length-4]byte, Length-4 bytes
+----------+
| Checksum | uint32, BigEndian, 4 bytes
+----------+

Message format
+---------+
| HeadLen | Varint32, 1 - 5 bytes
+---------+
| Head    | pb.Head
+---------+
| Body    | proto.Message
+---------+
*/

type Message struct {
	Head pb.Head
	Body proto.Message
}

func (m *Message) String() string {
	return fmt.Sprintf("Head: { %s}, Body: { %s}", m.Head.String(), m.Body.String())
}

func (m *Message) TextString() string {
	return fmt.Sprint(proto.MarshalTextString(&m.Head), proto.MarshalTextString(m.Body))
}

func (m *Message) Encode() ([]byte, error) {
	head, err := proto.Marshal(&m.Head)
	if err != nil {
		return nil, err
	}
	hlen := len(head)
	if hlen < 0 || hlen > math.MaxInt32 {
		return nil, ErrInvalidMessageHeadLength
	}
	body, err := proto.Marshal(m.Body)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err = endian.EncodeVarint(&buf, int64(hlen)); err != nil {
		return nil, err
	}
	if _, err = buf.Write(head); err != nil {
		return nil, err
	}
	if _, err = buf.Write(body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Message) Decode(data []byte) (err error) {
	hlen, n := binary.Varint(data)
	if n <= 0 || hlen < 0 || hlen > math.MaxInt32 || int(hlen) > len(data)-n {
		return ErrInvalidMessageHeadLength
	}

	hend := n + int(hlen)
	if err = proto.Unmarshal(data[n:hend], &m.Head); err != nil {
		return err
	}
	if m.Body, err = newMessageBody(m.Head.GetMsgName()); err != nil {
		return err
	}
	if err = proto.Unmarshal(data[hend:], m.Body); err != nil {
		return err
	}
	return nil
}

func newMessageBody(name string) (proto.Message, error) {
	typ := proto.MessageType(name)
	if typ == nil || typ.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("unknown message name: %s", name)
	}
	return reflect.New(typ.Elem()).Interface().(proto.Message), nil
}

func Write(w io.Writer, m *Message) error {
	payload, err := m.Encode()
	if err != nil {
		return err
	}
	length := len(payload) + 4
	if length <= 0 || length > math.MaxInt32 {
		return ErrInvalidPacketLength
	}
	checksum := crc32.ChecksumIEEE(payload)

	if err = endian.EncodeVarint(w, int64(length)); err != nil {
		return err
	}
	if _, err = w.Write(payload); err != nil {
		return err
	}
	if err = endian.BigEndian.EncodeUint32(w, checksum); err != nil {
		return err
	}
	return nil
}

func Read(r io.Reader) (*Message, error) {
	length, err := endian.DecodeVarint(r)
	if err != nil {
		return nil, err
	}
	if length <= 0 || length > math.MaxInt32 {
		return nil, ErrInvalidPacketLength
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	payload := buf[:length-4]
	checksum := binary.BigEndian.Uint32(buf[length-4:])
	if checksum != crc32.ChecksumIEEE(payload) {
		return nil, ErrInvalidChecksum
	}

	var m Message
	if err := m.Decode(payload); err != nil {
		return nil, err
	}
	return &m, nil
}
