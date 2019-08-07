package uuid

import "encoding/hex"

type UUID []byte

// New returns a UUID
func New() UUID {
	return NewRandom()
}

func (u UUID) String() string {
	return string(encodeHex(u))
}

func encodeHex(uuid UUID) []byte {
	buf := make([]byte, 36)
	hex.Encode(buf[:], uuid[:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], uuid[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], uuid[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], uuid[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], uuid[10:])
	return buf
}
