package codec

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

type Codec interface {
	ContentType() string
	Encode(io.Writer, interface{}) error
	Decode(io.Reader, interface{}) error
	EncodeError(io.Writer, Error) error
	DecodeError(io.Reader, *Error) error
}

var DefaultCodec JSONCodec

var _ Codec = JSONCodec{}

type JSONCodec struct{}

func (c JSONCodec) ContentType() string {
	return "application/json"
}

func (c JSONCodec) Encode(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func (c JSONCodec) Decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func (c JSONCodec) EncodeError(w io.Writer, e Error) error {
	return json.NewEncoder(w).Encode(e)
}

func (c JSONCodec) DecodeError(r io.Reader, e *Error) error {
	return json.NewDecoder(r).Decode(e)
}

var _ Codec = XMLCodec{}

type XMLCodec struct{}

func (c XMLCodec) ContentType() string {
	return "application/xml"
}

func (c XMLCodec) Encode(w io.Writer, v interface{}) error {
	return xml.NewEncoder(w).Encode(v)
}

func (c XMLCodec) Decode(r io.Reader, v interface{}) error {
	return xml.NewDecoder(r).Decode(v)
}

func (c XMLCodec) EncodeError(w io.Writer, e Error) error {
	return xml.NewEncoder(w).Encode(e)
}

func (c XMLCodec) DecodeError(r io.Reader, e *Error) error {
	return xml.NewDecoder(r).Decode(e)
}
