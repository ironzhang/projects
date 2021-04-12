package base58

import (
	"bytes"
	"testing"
)

func TestBase58Encode(t *testing.T) {
	v := 1234

	tests := []struct {
		src []byte
		dst []byte
	}{
		{
			src: []byte{},
			dst: []byte(""),
		},
		{
			src: []byte{0},
			dst: []byte("1"),
		},
		{
			src: []byte{byte(v >> 8), byte(v)},
			dst: []byte("NH"),
		},
		{
			src: []byte{byte(v), byte(v >> 8)},
			dst: []byte("Gyy"),
		},
	}

	enc := NewEncoding(encodeStd)
	for i, tt := range tests {
		if got, want := enc.Encode(tt.src), tt.dst; !bytes.Equal(got, want) {
			t.Errorf("%d: encode: got %q, want %q", i, got, want)
		} else {
			t.Logf("%d: encode: got %q", i, got)
		}
	}
}

func TestBase58Decode(t *testing.T) {
	v := 1234

	tests := []struct {
		src []byte
		dst []byte
	}{
		{
			src: []byte(""),
			dst: []byte{},
		},
		{
			src: []byte("1"),
			dst: []byte{0},
		},
		{
			src: []byte("NH"),
			dst: []byte{byte(v >> 8), byte(v)},
		},
		{
			src: []byte("Gyy"),
			dst: []byte{byte(v), byte(v >> 8)},
		},
	}
	enc := NewEncoding(encodeStd)
	for i, tt := range tests {
		if got, want := enc.Decode(tt.src), tt.dst; !bytes.Equal(got, want) {
			t.Errorf("%d: encode: got %q, want %q", i, got, want)
		} else {
			t.Logf("%d: encode: got %q", i, got)
		}
	}
}

func TestBase58(t *testing.T) {
	tests := [][]byte{
		{},
		{0},
		{0, 1},
		{1, 0},
		{1, 0},
		{0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 10, 0},
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}
	for i, tt := range tests {
		encb := StdEncoding.Encode(tt)
		decb := StdEncoding.Decode(encb)
		if got, want := decb, tt; !bytes.Equal(got, want) {
			t.Errorf("%d: encode(%v) decode: got %q, want %q", i, tt, got, want)
		} else {
			t.Logf("%d: encode(%v): %q, decode(%q): %v", i, tt, encb, encb, decb)
		}
	}
}
