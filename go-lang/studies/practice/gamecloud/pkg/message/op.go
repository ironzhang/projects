package message

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/golang/protobuf/proto"
)

var overflow = errors.New("binary: varint overflows a 64-bit integer")

func readUvarintBuf(r io.Reader) (buf []byte, err error) {
	b := make([]byte, 1)
	for {
		if len(buf) >= binary.MaxVarintLen64 {
			return buf, overflow
		}
		_, err = r.Read(b)
		if err != nil {
			return buf, err
		}
		buf = append(buf, b[0])
		if b[0] < 0x80 {
			return buf, nil
		}
	}
}

func readLength(r io.Reader) (uint64, error) {
	b, err := readUvarintBuf(r)
	if err != nil {
		return 0, err
	}
	x, _ := binary.Uvarint(b)
	return x, nil
}

func writeLength(w io.Writer, x uint64) error {
	b := make([]byte, binary.MaxVarintLen64)
	i := binary.PutUvarint(b, x)
	_, err := w.Write(b[:i])
	return err
}

func readPacket(r io.Reader, m proto.Message) error {
	len, err := readLength(r)
	if err != nil {
		return err
	}
	buf := make([]byte, len)
	if _, err = io.ReadFull(r, buf); err != nil {
		return err
	}
	return proto.Unmarshal(buf, m)
}

func writePacket(w io.Writer, m proto.Message) error {
	buf, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	if err = writeLength(w, uint64(len(buf))); err != nil {
		return err
	}
	if _, err = w.Write(buf); err != nil {
		return err
	}
	return nil
}
