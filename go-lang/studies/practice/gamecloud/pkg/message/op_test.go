package message

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/ironzhang/gamecloud/pkg/message/testpb"
)

func PutUvarints(nums []uint64) []byte {
	pos := 0
	buf := make([]byte, binary.MaxVarintLen64*len(nums))
	for _, x := range nums {
		pos += binary.PutUvarint(buf[pos:], x)
	}
	return buf[:pos]
}

func GetUvarints(buf []byte) []uint64 {
	pos := 0
	nums := []uint64{}
	for pos < len(buf) {
		x, i := binary.Uvarint(buf[pos:])
		nums = append(nums, x)
		pos += i
	}
	return nums
}

func ReadLengths(r io.Reader) ([]uint64, error) {
	nums := []uint64{}
	for {
		x, err := readLength(r)
		if err != nil {
			if err == io.EOF {
				return nums, nil
			}
			return nums, err
		}
		nums = append(nums, x)
	}
}

func WriteLengths(w io.Writer, nums []uint64) error {
	for _, x := range nums {
		if err := writeLength(w, x); err != nil {
			return err
		}
	}
	return nil
}

func TestReadLength(t *testing.T) {
	want := []uint64{
		0, 1, 2, 3, 4, 100, 200, 300, 400, 1000, 2000, 3000, 4000, 10000, 20000, 100000, 200000,
		0xFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xEFFFFFFFFFFFFFFF,
	}

	b := PutUvarints(want)
	r := bytes.NewReader(b)
	got, err := ReadLengths(r)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestWriteLength(t *testing.T) {
	want := []uint64{
		0, 1, 2, 3, 4, 100, 200, 300, 400, 1000, 2000, 3000, 4000, 10000, 20000, 100000, 200000,
		0xFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xEFFFFFFFFFFFFFFF,
	}

	var buf bytes.Buffer
	err := WriteLengths(&buf, want)
	if err != nil {
		t.Fatal(err)
	}
	got := GetUvarints(buf.Bytes())
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestReadWriteLength(t *testing.T) {
	want := []uint64{
		0, 1, 2, 3, 4, 100, 200, 300, 400, 1000, 2000, 3000, 4000, 10000, 20000, 100000, 200000,
		0xFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xEFFFFFFFFFFFFFFF,
	}

	var buf bytes.Buffer
	err := WriteLengths(&buf, want)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ReadLengths(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestReadWritePacket(t *testing.T) {
	wMessages := []proto.Message{
		&testpb.Header{},
		&testpb.GetAccessPointReq{Header: &testpb.Header{Seq: 1}, AccessType: "MQTT"},
		&testpb.GetAccessPointRsp{Addrs: []string{"127.0.0.1", "127.0.0.2"}},
	}
	rMessages := []proto.Message{
		&testpb.Header{},
		&testpb.GetAccessPointReq{},
		&testpb.GetAccessPointRsp{},
	}

	var buf bytes.Buffer
	for _, msg := range wMessages {
		if err := writePacket(&buf, msg); err != nil {
			t.Fatalf("WriteProtoMessage: %v", err)
		}
	}
	for _, msg := range rMessages {
		if err := readPacket(&buf, msg); err != nil {
			if err == io.EOF {
				break
			}
			t.Fatalf("ReadProtoMessage: %v", err)
		}
	}
	if !reflect.DeepEqual(rMessages, wMessages) {
		t.Errorf("%v != %v", rMessages, wMessages)
	}
}
