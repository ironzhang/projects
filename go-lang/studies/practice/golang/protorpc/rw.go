package protorpc

import (
	"encoding/binary"
	"io"

	"github.com/golang/protobuf/proto"
)

func writeFrame(w io.Writer, frame proto.Message) error {
	if frame != nil {
		data, err := proto.Marshal(frame)
		if err != nil {
			return err
		}
		size := int32(len(data))
		if err = binary.Write(w, binary.BigEndian, size); err != nil {
			return err
		}
		if _, err = w.Write(data); err != nil {
			return err
		}
	}
	return nil
}

func writeMessage(w io.Writer, head proto.Message, body proto.Message) error {
	if err := writeFrame(w, head); err != nil {
		return err
	}
	if err := writeFrame(w, body); err != nil {
		return err
	}
	return nil
}

func readFrame(r io.Reader, frame proto.Message) error {
	if frame != nil {
		var size int32
		if err := binary.Read(r, binary.BigEndian, &size); err != nil {
			return err
		}
		data := make([]byte, size)
		if _, err := io.ReadFull(r, data); err != nil {
			return err
		}
		if err := proto.Unmarshal(data, frame); err != nil {
			return err
		}
	}
	return nil
}
